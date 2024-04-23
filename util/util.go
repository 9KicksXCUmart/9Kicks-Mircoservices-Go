/*
Package util provides utility functions for the application.
*/
package util

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func StructToAttributeValue(s interface{}) (map[string]types.AttributeValue, error) {
	av, err := attributevalue.MarshalMap(s)
	if err != nil {
		return nil, err
	}

	return av, nil
}

func GetCurrentTime() string {
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	layout := "2006-01-02T15:04:05"

	return time.Now().In(loc).Format(layout)
}
