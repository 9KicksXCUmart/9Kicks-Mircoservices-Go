package controller

import (
	"9Kicks/config"
	"9Kicks/model/auth"
	"9Kicks/util"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtKey                = []byte(config.GetJWTSecrets().JWTUserSecret)
	dynamoDBClient        = config.GetDynamoDBClient()
	tableName             = "9Kicks-test"
	gsiName               = "User-email-index"
	indexPartitionKeyName = "email"
)

func Signup(c *gin.Context) {
	var user auth.UserSignUpForm
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the email already exists
	queryParams := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("User-email-index"),
		KeyConditionExpression: aws.String("#pk = :email"),
		ExpressionAttributeNames: map[string]string{
			"#pk": "email",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: user.Email},
		},
	}

	result, err := dynamoDBClient.Query(context.TODO(), queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email"})
		return
	}

	if len(result.Items) > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Hash the input password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Generate verification token for email verification
	verificationToken := uuid.New().String()
	tokenExpirationTime := time.Now().Add(time.Minute * 5).Unix()

	// Construct the user profile item to be stored in DynamoDB
	userProfile := auth.UserProfile{
		PK:                "USER#" + uuid.New().String(),
		SK:                "USER_PROFILE",
		Email:             user.Email,
		Password:          string(hashedPassword),
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		VerificationToken: verificationToken,
		TokenExpiry:       tokenExpirationTime,
	}

	// Convert the struct to dynamodb.AttributeValue
	profileItem, err := util.StructToAttributeValue(userProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert struct to attribute value"})
		return
	}

	// Put the user profile item into DynamoDB
	putParams := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      profileItem,
	}

	_, err = dynamoDBClient.PutItem(context.TODO(), putParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
	var user auth.UserLoginForm

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Query the user profile by email
	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String(gsiName),
		KeyConditionExpression: aws.String("#pk = :email"),
		ExpressionAttributeNames: map[string]string{
			"#pk": indexPartitionKeyName,
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: user.Email},
		},
	}

	result, err := dynamoDBClient.Query(context.TODO(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if len(result.Items) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	item := result.Items[0]

	// Compare the input password with the stored password
	storedPassword := item["password"].(*types.AttributeValueMemberS).Value

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check if the email is verified
	isVerified := item["isVerified"].(*types.AttributeValueMemberBOOL).Value
	if !isVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not verified"})
		return
	}

	// Set jwt token expiration to be 1 hour
	expirationTime := time.Now().Add(time.Hour)
	claims := &auth.Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set the jwt token in a cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode, // Set SameSite policy to Strict
	})

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func ValidateToken(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	tokenString := parts[1]

	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(*auth.Claims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"email": claims.Email})
}

func VerifyEmail(c *gin.Context) {
	// Get the verification token and email from the request parameters
	token := c.Query("token")
	email := c.Query("email")
	if token == "" || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token or email"})
		return
	}

	// Check if the token and email match the stored values in your database
	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String(gsiName),
		KeyConditionExpression: aws.String("#pk = :email"),
		ExpressionAttributeNames: map[string]string{
			"#pk": indexPartitionKeyName,
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	}

	result, err := dynamoDBClient.Query(context.TODO(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email"})
		return
	}

	if len(result.Items) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email does not exist"})
		return
	}

	item := result.Items[0]

	storedToken := item["verificationToken"].(*types.AttributeValueMemberS).Value
	tokenExpirationTime, _ := strconv.ParseInt(item["tokenExpiry"].(*types.AttributeValueMemberN).Value, 10, 64)
	fmt.Println(tokenExpirationTime)

	if storedToken != token {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid verification token"})
		return
	}

	// Check if the token has expired
	if time.Now().Unix() > tokenExpirationTime {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Verification token has expired"})
		return
	}

	// Update the user profile item in DynamoDB to mark the user as verified
	updateParams := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: item["PK"].(*types.AttributeValueMemberS).Value},
			"SK": &types.AttributeValueMemberS{Value: item["SK"].(*types.AttributeValueMemberS).Value},
		},
		UpdateExpression: aws.String("SET #isVerified = :isVerified"),
		ExpressionAttributeNames: map[string]string{
			"#isVerified": "isVerified",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":isVerified": &types.AttributeValueMemberBOOL{Value: true},
		},
	}

	_, err = dynamoDBClient.UpdateItem(context.TODO(), updateParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
