package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const tableName = "ResumeAnalytics"

var dbClient *dynamodb.DynamoDB

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleRequest(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var body string
	switch req.HTTPMethod {
	case "GET":
		return getCount(req)
	case "PUT":
		// return SetCount(req)
		body = "Received PUT request"
	case "DELETE":
		// return ResetCount(req)
		body = "Received DELETE request"
	default:
		// return UnhandledMethod()
		body = "Received UnhandledMethod() request"
	}

	return &events.APIGatewayProxyResponse{Body: body, StatusCode: 200}, nil
}

func getCount(events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	input := &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("VisitorCount"),
		TableName:              aws.String(tableName),
	}
	result, err := dbClient.Query(input)
	check(err)
	return &events.APIGatewayProxyResponse{Body: result.String(), StatusCode: 200}, nil
}

func main() {
	// Connect to DB
	sess, err := session.NewSession()
	check(err)
	dbClient = dynamodb.New(sess)

	// Handle lambda request from API gateway
	lambda.Start(handleRequest)
}
