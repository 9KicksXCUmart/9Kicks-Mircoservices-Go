package auth_service

import (
	"9Kicks/dao"
	"9Kicks/model/auth_model"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CheckEmailExists(email string) (bool, error) {
	userProfiles, err := dao.GetUserProfileByEmail(email)
	if err != nil {
		return false, err
	}

	return len(userProfiles) > 0, nil
}

func GetUserProfileByEmail(email string) (auth_model.UserProfile, error) {
	userProfiles, err := dao.GetUserProfileByEmail(email)
	if err != nil {
		return userProfiles[0], err
	}

	return userProfiles[0], nil
}

func CreateUser(email, firstName, lastName, password string) (token string, success bool) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return "", false
	}

	verificationToken, tokenExpirationTime := generateVerificationToken()

	// Construct the user profile item to be stored in DynamoDB
	userProfile := auth_model.UserProfile{
		PK:                "USER#" + uuid.New().String(),
		SK:                "USER_PROFILE",
		Email:             email,
		Password:          string(hashedPassword),
		FirstName:         firstName,
		LastName:          lastName,
		VerificationToken: verificationToken,
		TokenExpiry:       tokenExpirationTime,
	}

	return verificationToken, dao.AddNewUserProfile(userProfile)
}

func IsValidPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateJWT(secretKey, email, userID string) (string, time.Time, error) {
	// Set jwt token expiration to be 1 hour
	expirationTime := time.Now().Add(time.Hour)
	claims := &auth_model.Claims{
		Email:  email,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func UpdateVerificationToken(userId string) (string, int64, error) {
	verificationToken, tokenExpirationTime := generateVerificationToken()
	err := dao.UpdateVerificationToken(userId, verificationToken, tokenExpirationTime)
	if err != nil {
		return "", 0, err
	}

	return verificationToken, tokenExpirationTime, nil
}

func generateVerificationToken() (string, int64) {
	// Set verification token expiry to be 5 minutes
	tokenExpirationTime := time.Now().Add(time.Minute * 5).Unix()
	verificationToken := uuid.New().String()
	return verificationToken, tokenExpirationTime
}

func VerifyUserEmail(userId string) error {
	return dao.VerifyUserEmail(userId)
}
