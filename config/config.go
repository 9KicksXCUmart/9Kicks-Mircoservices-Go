package config

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/joho/godotenv"
	"github.com/opensearch-project/opensearch-go"
)

var (
	cfg              aws.Config
	awsErr           error
	openSearchClient *opensearch.Client
)

type Secrets struct {
	JWTUserSecret  string `json:"JWT_USER_SECRET"`
	JWTAdminSecret string `json:"JWT_ADMIN_SECRET"`
}

type TokenParams struct {
	ZOHOTokenParams string `json:"ZOHO_TOKEN_PARAMS"`
	ZOHOAccountNo   string `json:"ZOHO_ACCOUNT_NO"`
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

func GetS3Client() *s3.Client {
	return s3.NewFromConfig(cfg)
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

func GetTokenParams() TokenParams {
	secretName := "9Kicks-email"

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
	var tokenParams TokenParams
	err = json.Unmarshal([]byte(*result.SecretString), &tokenParams)
	if err != nil {
		log.Fatal(err.Error())
	}
	return tokenParams
}

func GetOpenSearchClient() *opensearch.Client {
	openSearchClient, _ = opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{os.Getenv("OPENSEARCH_ENDPOINT")},
		Username:  os.Getenv("OPENSEARCH_USERNAME"),
		Password:  os.Getenv("OPENSEARCH_PASSWORD"),
	})

	return openSearchClient
}
