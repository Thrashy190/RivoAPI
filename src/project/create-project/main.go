package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/rivo-api/shared/models/project"
)

var (
	client    *dynamodb.Client
	tableName string
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client = dynamodb.NewFromConfig(cfg)

	var ok bool
	tableName, ok = os.LookupEnv("TABLE_NAME")
	if !ok || tableName == "" {
		log.Fatal("TABLE_NAME env variable required")
	}
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var body project.CreateProjectRequest

	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error":"invalid body"}`,
		}, nil
	}

	if body.Name == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error":"name is required"}`,
		}, nil
	}

	p := project.Project{
		ID:          uuid.NewString(),
		Name:        body.Name,
		Description: body.Description,
		CreatedAt:   time.Now().UTC().UnixMilli(),
	}

	item, err := attributevalue.MarshalMap(p)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	respBody, _ := json.Marshal(p)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(respBody),
	}, nil
}

func main() {
	lambda.Start(handler)
}
