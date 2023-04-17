package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const tableName = "site-analytics"
const keyName = "visitorCount"
const partitionKey = "metrics"
const site = "resume.braheezy.net"

var dbClient *dynamodb.Client

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleRequest(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var body string
	var err error
	switch req.HTTPMethod {
	case "GET":
		body, err = getCount()
	case "PUT":
		body, err = updateCount()
	default:
		// This is probably a terrible default case
		_, err = checkTable()
	}

	if err != nil {
		return &events.APIGatewayProxyResponse{Body: "Oh no! Encountered an error", StatusCode: 400}, err
	}

	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET,PUT",
	}

	return &events.APIGatewayProxyResponse{
			Body:       body,
			StatusCode: 200,
			Headers:    headers},
		nil
}

func getCount() (string, error) {
	key := map[string]types.AttributeValue{
		partitionKey: &types.AttributeValueMemberS{Value: site},
	}

	// Define the DynamoDB expression to retrieve the visitorCount attribute
	expr, err := expression.NewBuilder().
		WithProjection(expression.NamesList(expression.Name(keyName))).
		Build()
	check(err)

	// Define the GetItem input with the expression
	input := &dynamodb.GetItemInput{
		TableName:                aws.String(tableName),
		Key:                      key,
		ProjectionExpression:     expr.Projection(),
		ExpressionAttributeNames: expr.Names(),
	}

	// Retrieve the item from the DynamoDB table
	output, err := dbClient.GetItem(context.Background(), input)
	check(err)

	// Extract the value of the visitorCount attribute
	count := ""
	if countAttr, ok := output.Item[keyName]; ok {
		count = countAttr.(*types.AttributeValueMemberN).Value
	}

	return count, nil
}

func checkTable() (bool, error) {
	exists := true
	_, err := dbClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			log.Printf("Table %v does not exist.\n", tableName)
			err = nil
		} else {
			log.Printf("Couldn't determine existence of table %v. Here's why: %v\n", tableName, err)
		}
		exists = false
	}
	return exists, err
}

func updateCount() (string, error) {
	// Specify the key of the item to update
	key := map[string]types.AttributeValue{
		"metrics": &types.AttributeValueMemberS{Value: site},
	}

	updateExpr := fmt.Sprintf("SET %[1]v = %[1]v + :incr", keyName)
	updateInput := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeValues: map[string]types.AttributeValue{":incr": &types.AttributeValueMemberN{Value: "1"}},
		ReturnValues:              types.ReturnValueUpdatedNew,
	}

	// Execute the update operation and print the updated item
	output, err := dbClient.UpdateItem(context.Background(), updateInput)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`{"count":"%v"}`, output.Attributes["visitorCount"].(*types.AttributeValueMemberN).Value), nil
}
func main() {
	// Connect to DB
	cfg, err := config.LoadDefaultConfig(context.TODO())
	check(err)
	dbClient = dynamodb.NewFromConfig(cfg)

	// Handle lambda request from API gateway
	lambda.Start(handleRequest)
}
