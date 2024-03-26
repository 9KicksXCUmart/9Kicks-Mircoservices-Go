package controller

import (
	"9Kicks/config"
	"9Kicks/model/auth"
	"9Kicks/util"
	"context"
	"fmt"
	"net/http"
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
	jwtKey         = []byte(config.GetJWTSecrets().JWTUserSecret)
	dynamoDBClient = config.GetDynamoDBClient()
	tableName      = "9Kicks-test"
	gsiName        = "User-email-index"
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	userProfile := auth.UserProfile{
		PK:        "USER#" + uuid.New().String(),
		SK:        "USER_PROFILE",
		Email:     user.Email,
		Password:  string(hashedPassword),
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	profileItem, err := util.StructToAttributeValue(userProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert struct to attribute value"})
		return
	}

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
	fmt.Println("DLLM")
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	IndexPartitionKeyName := "email"

	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String(gsiName),
		KeyConditionExpression: aws.String("#pk = :email"),
		ExpressionAttributeNames: map[string]string{
			"#pk": IndexPartitionKeyName,
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

	storedPassword := item["password"].(*types.AttributeValueMemberS).Value

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Set token expiration to be 1 hour
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
