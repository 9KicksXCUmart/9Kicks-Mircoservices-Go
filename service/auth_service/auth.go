package auth_service

import (
	"9Kicks/dao"
	"9Kicks/model/auth_model"
	"log"
	"time"

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

func generateVerificationToken() (string, int64) {
	// Set verification token expiry to be 5 minutes
	tokenExpirationTime := time.Now().Add(time.Minute * 5).Unix()
	verificationToken := uuid.New().String()
	return verificationToken, tokenExpirationTime
}
