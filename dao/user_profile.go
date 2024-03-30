package dao

import (
	"9Kicks/config"
	"9Kicks/model/auth_model"
	"9Kicks/util"
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	dynamoDBClient = config.GetDynamoDBClient()
	tableName      = os.Getenv("DB_TABLE_NAME")
)

func GetUserProfileByEmail(email string) ([]auth_model.UserProfile, error) {
	var userProfiles []auth_model.UserProfile

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
