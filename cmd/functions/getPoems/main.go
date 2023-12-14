package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

var TableName string = os.Getenv("POEMS_TABLE_NAME")

var db dynamodb.Client

func init() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	db = *dynamodb.NewFromConfig(sdkConfig)
}

type Poem struct {
	Id   string `json:"id" dynamodbav:"id"`
	Poem string `json:"poem" dynamodbav:"poem"`
}

func listItems(ctx context.Context) ([]Poem, error) {
	poems := make([]Poem, 0)
	var token map[string]types.AttributeValue

	for {
		input := &dynamodb.ScanInput{
			TableName:         aws.String(TableName),
			ExclusiveStartKey: token,
		}

		result, err := db.Scan(ctx, input)
		if err != nil {
			return nil, err
		}

		var fetchedPoems []Poem
		err = attributevalue.UnmarshalListOfMaps(result.Items, &fetchedPoems)
		if err != nil {
			return nil, err
		}

		poems = append(poems, fetchedPoems...)
		token = result.LastEvaluatedKey
		if token == nil {
			break
		}
	}

	return poems, nil
}

func getItem(ctx context.Context, id string) (*Poem, error) {
	key, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"id": key,
		},
	}

	log.Printf("Calling Dynamodb with input: %v", input)
	result, err := db.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}
	log.Printf("Executed GetItem DynamoDb successfully. Result: %#v", result)

	if result.Item == nil {
		return nil, nil
	}

	poem := new(Poem)
	err = attributevalue.UnmarshalMap(result.Item, poem)
	if err != nil {
		return nil, err
	}

	return poem, nil
}

func processGet(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := req.PathParameters["id"]
	if !ok {
		return processGetPoems(ctx)
	} else {
		return processGetPoem(ctx, id)
	}
}

func processGetPoem(ctx context.Context, id string) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received GET poem request with id = %s", id)

	poem, err := getItem(ctx, id)
	if err != nil {
		return serverError(err)
	}

	if poem == nil {
		return clientError(http.StatusNotFound)
	}

	json, err := json.Marshal(poem)
	if err != nil {
		return serverError(err)
	}
	log.Printf("Successfully fetched poem item %s", json)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(json),
	}, nil
}

func processGetPoems(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	log.Print("Received GET poems request")

	poems, err := listItems(ctx)
	if err != nil {
		return serverError(err)
	}

	json, err := json.Marshal(poems)
	if err != nil {
		return serverError(err)
	}
	log.Printf("Successfully fetched poems: %s", json)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(json),
	}, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {

	return events.APIGatewayProxyResponse{
		Body:       http.StatusText(status),
		StatusCode: status,
	}, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	log.Println(err.Error())

	return events.APIGatewayProxyResponse{
		Body:       http.StatusText(http.StatusInternalServerError),
		StatusCode: http.StatusInternalServerError,
	}, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("Received request %#v", request)

	switch request.HTTPMethod {
	case "GET":
		return processGet(ctx, request)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func main() {
	lambda.Start(handler)
}
