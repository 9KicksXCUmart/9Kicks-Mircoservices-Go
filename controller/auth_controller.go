package controller

import (
	"9Kicks/config"
	"9Kicks/model/auth_model"
	"9Kicks/service/auth_service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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
	tableName             = os.Getenv("DB_TABLE_NAME")
	gsiName               = "User-email-index"
	indexPartitionKeyName = "email"
)

func Signup(c *gin.Context) {
	var user auth_model.UserSignUpForm
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, _ := auth_service.CheckEmailExists(user.Email)
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	verificationToken, success := auth_service.CreateUser(user.Email, user.FirstName, user.LastName, user.Password)
	if !success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Send verification email
	err := auth_service.SendEmailTo(user.Email, verificationToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "A verification email has been sent to your email address. Please verify your email to complete registration"})
}

func Login(c *gin.Context) {
	var user auth_model.UserLoginForm

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
	claims := &auth_model.Claims{
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

	token, err := jwt.ParseWithClaims(tokenString, &auth_model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(*auth_model.Claims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"email": claims.Email})
}

func ResendVerificationEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token or email"})
		return
	}
	// Generate verification token for email verification
	verificationToken := uuid.New().String()
	// Set verification token expiry to be 5 minutes

	// TODO: Update token expiry time
	tokenExpirationTime := time.Now().Add(time.Minute * 5).Unix()

	queryParams := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String("User-email-index"),
		KeyConditionExpression: aws.String("#pk = :email"),
		ExpressionAttributeNames: map[string]string{
			"#pk": "email",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	}

	result, err := dynamoDBClient.Query(context.TODO(), queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email"})
		return
	}

	userId := result.Items[0]["PK"].(*types.AttributeValueMemberS).Value

	// Define the expression attribute names and values
	expr := aws.String("SET #verificationToken = :verificationToken, #tokenExpiry = :tokenExpiry")
	exprNames := map[string]string{
		"#verificationToken": "verificationToken",
		"#tokenExpiry":       "tokenExpiry",
	}
	exprValues := map[string]types.AttributeValue{
		":verificationToken": &types.AttributeValueMemberS{Value: verificationToken},
		":tokenExpiry":       &types.AttributeValueMemberN{Value: strconv.FormatInt(tokenExpirationTime, 10)},
	}

	// Construct the UpdateItemInput
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: userId},
			"SK": &types.AttributeValueMemberS{Value: "USER_PROFILE"},
		},
		UpdateExpression:          expr,
		ExpressionAttributeNames:  exprNames,
		ExpressionAttributeValues: exprValues,
	}

	_, err = dynamoDBClient.UpdateItem(context.TODO(), updateInput)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update verification token"})
		return
	}

	err = auth_service.SendEmailTo(email, verificationToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "A verification email has been sent to your email address. Please verify your email to complete registration"})
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
	isVerified := item["isVerified"].(*types.AttributeValueMemberBOOL).Value
	if isVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email already verified"})
		return
	}

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
