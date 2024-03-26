package config

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/joho/godotenv"
)

var (
	cfg    aws.Config
	awsErr error
)

type Secrets struct {
	JWTUserSecret  string `json:"JWT_USER_SECRET"`
	JWTAdminSecret string `json:"JWT_ADMIN_SECRET"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION_NAME")

	staticCredProvider := credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")
	cfg, awsErr = config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(staticCredProvider), config.WithRegion(awsRegion))

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
}

func GetDynamoDBClient() *dynamodb.Client {
	return dynamodb.NewFromConfig(cfg)
}

func GetJWTSecrets() Secrets {
	secretName := "9Kicks_secrets"

	// Create Secrets Manager client
	secretsManagerClient := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := secretsManagerClient.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Decrypts secret using the associated KMS key.
	var secrets Secrets
	err = json.Unmarshal([]byte(*result.SecretString), &secrets)
	if err != nil {
		log.Fatal(err.Error())
	}
	return secrets
}

//func init() {
//
//	fmt.Println("config init")
//awsAccessKeyID := "AKIA333U776S2ZMRH3LT"
//awsSecretAccessKey := "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"
//staticCredProvider := credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")
//cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(staticCredProvider), config.WithRegion("ap-southeast-1"))
//if err != nil {
//	log.Fatalf("unable to load SDK config, %v", err)
//}
//
//// Using the Config value, create the DynamoDB client
//svc := dynamodb.NewFromConfig(cfg)
//
//// Build the request with its input parameters
//resp, err := svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
//	Limit: aws.Int32(5),
//})
//if err != nil {
//	log.Fatalf("failed to list tables, %v", err)
//}
//
//fmt.Println("Tables:")
//for _, tableName := range resp.TableNames {
//	fmt.Println(tableName)
//}
//}
