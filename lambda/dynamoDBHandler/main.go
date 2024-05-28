package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"encoding/json"
	"net/http"
)

var (
	tableName = "MyDynamoDB"
	dynamoDb  = dynamodb.New(session.Must(session.NewSession()))
)

type Item struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func getItem(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.QueryStringParameters["id"]

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	result, err := dynamoDb.GetItem(input)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	item := Item{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	body, _ := json.Marshal(item)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func putItem(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	item := Item{}
	err := json.Unmarshal([]byte(request.Body), &item)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(item.ID),
			},
			"content": {
				S: aws.String(item.Content),
			},
		},
	}
	_, err = dynamoDb.PutItem(input)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
	}, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "GET":
		return getItem(request)
	case "POST":
		return putItem(request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
		}, nil
	}
}

func main() {
	lambda.Start(handler)
}