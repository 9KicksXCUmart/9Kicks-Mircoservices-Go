package util

import (
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
