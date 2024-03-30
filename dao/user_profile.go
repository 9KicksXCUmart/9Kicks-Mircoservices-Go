package dao

import (
	"9Kicks/config"
	"9Kicks/model/auth_model"
	"9Kicks/util"
	"context"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	dynamoDBClient = config.GetDynamoDBClient()
	tableName      = os.Getenv("DB_TABLE_NAME")
	emailIndexName = "User-email-index"
)

func GetUserProfileByEmail(email string) ([]auth_model.UserProfile, error) {
	var userProfiles []auth_model.UserProfile

	queryParams := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String(emailIndexName),
		KeyConditionExpression: aws.String("#pk = :email"),
		ExpressionAttributeNames: map[string]string{
			"#pk": "email",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	}

	result, err := dynamoDBClient.Query(context.TODO(), queryParams)
	if err != nil || len(result.Items) == 0 {
		return userProfiles, err
	}
	for _, item := range result.Items {
		var userProfile auth_model.UserProfile
		err := attributevalue.UnmarshalMap(item, &userProfile)
		if err != nil {
			return userProfiles, err
		}
		userProfiles = append(userProfiles, userProfile)
	}

	return userProfiles, nil
}

func AddNewUserProfile(userProfile auth_model.UserProfile) (success bool) {
	// Convert the struct to dynamodb.AttributeValue
	profileItem, err := util.StructToAttributeValue(userProfile)
	if err != nil {
		return false
	}

	// Put the user profile item into DynamoDB
	putParams := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      profileItem,
	}

	_, err = dynamoDBClient.PutItem(context.TODO(), putParams)
	if err != nil {
		return false
	}

	return true
}

func UpdateVerificationToken(userId, verificationToken string, tokenExpirationTime int64) error {
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

	_, err := dynamoDBClient.UpdateItem(context.TODO(), updateInput)
	if err != nil {
		return err
	}

	return nil
}

func VerifyUserEmail(userId string) error {
	updateParams := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: userId},
			"SK": &types.AttributeValueMemberS{Value: "USER_PROFILE"},
		},
		UpdateExpression: aws.String("SET #isVerified = :isVerified"),
		ExpressionAttributeNames: map[string]string{
			"#isVerified": "isVerified",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":isVerified": &types.AttributeValueMemberBOOL{Value: true},
		},
	}

	_, err := dynamoDBClient.UpdateItem(context.TODO(), updateParams)
	if err != nil {
		return err
	}

	return nil
}
