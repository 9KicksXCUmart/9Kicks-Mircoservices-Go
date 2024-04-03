package dao

import (
	"9Kicks/model/review"
	"9Kicks/util"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func AddNewReview(productReview review.ProductReview) bool {

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

func GetReviewList(productId string) []review.ProductReview {
	var productReviews []review.ProductReview

	queryParams := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("#pk = :productId"),
		ExpressionAttributeNames: map[string]string{
			"#pk": "PK",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":productId": &types.AttributeValueMemberS{Value: productId},
		},
	}

	result, err := dynamoDBClient.Query(context.TODO(), queryParams)
	if err != nil || len(result.Items) == 0 {
		return productReviews
	}
	for _, item := range result.Items {
		var productReview review.ProductReview
		err := util.AttributeValueToStruct(item, &productReview)
		if err != nil {
			return productReviews
		}
		productReviews = append(productReviews, productReview)
	}

	return productReviews
}
