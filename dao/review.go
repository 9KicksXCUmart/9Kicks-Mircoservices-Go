/*
dao package responsible for handling all the data access operations related to reviews.
*/
package dao

import (
	"9Kicks/model/review"
	"9Kicks/util"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const IndexName = "email-SK-index"

func AddNewReview(productReview review.ProductReview) (success bool) {
	reviewItem, err := util.StructToAttributeValue(productReview)
	if err != nil {
		return false
	}
	_, err = dynamoDBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      reviewItem,
	})
	if err != nil {
		return false
	}
	return true
}

func GetReviewList(productId string) ([]review.ProductReview, error) {
	var productReviews []review.ProductReview
	productId = "PRODUCT#" + productId

	queryParams := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("#pk = :productId AND begins_with(#sk, :skPrefix)"),
		ExpressionAttributeNames: map[string]string{
			"#pk": "PK",
			"#sk": "SK",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":productId": &types.AttributeValueMemberS{Value: productId},
			":skPrefix":  &types.AttributeValueMemberS{Value: "REVIEW#"},
		},
	}

	result, err := dynamoDBClient.Query(context.TODO(), queryParams)
	if err != nil || len(result.Items) == 0 {
		return productReviews, err
	}
	for _, item := range result.Items {
		var productReview review.ProductReview
		err := attributevalue.UnmarshalMap(item, &productReview)
		if err != nil {
			return productReviews, err
		}
		productReviews = append(productReviews, productReview)
	}

	return productReviews, nil
}

func GetUserReviews(email string) ([]review.ProductReview, error) {
	var userReviews []review.ProductReview

	queryParams := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String(IndexName),
		KeyConditionExpression: aws.String("#email = :email AND begins_with(#sk, :skPrefix)"),
		ExpressionAttributeNames: map[string]string{
			"#email": "email",
			"#sk":    "SK",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email":    &types.AttributeValueMemberS{Value: email},
			":skPrefix": &types.AttributeValueMemberS{Value: "REVIEW#"},
		},
	}

	result, err := dynamoDBClient.Query(context.TODO(), queryParams)
	if err != nil || len(result.Items) == 0 {
		return userReviews, err
	}
	for _, item := range result.Items {
		var userReview review.ProductReview
		err := attributevalue.UnmarshalMap(item, &userReview)
		if err != nil {
			return userReviews, err
		}
		userReviews = append(userReviews, userReview)
	}

	return userReviews, nil
}
