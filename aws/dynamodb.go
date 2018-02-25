package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bitte-ein-bit/songbeamer-helper/util"
)

var dynamodbsvc *dynamodb.DynamoDB

const region = "eu-central-1"
const tablename = "songbeamer-helper"

type NumericItem struct {
	Value int    `json:"numericValue"`
	Key   string `json:"key"`
}

func GetDynamoDBNumericItem(key string) NumericItem {
	connectDynamoDB()
	result, err := dynamodbsvc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String(key),
			},
		},
	})
	util.CheckForError(err)

	item := NumericItem{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	util.CheckForError(err)

	return item
}

func UpdateDynamoDBItem(key string, value int) {
	connectDynamoDB()
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				N: aws.String(fmt.Sprintf("%v", value)),
			},
		},
		TableName: aws.String(tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String(key),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set numericValue = :r"),
	}

	_, err := dynamodbsvc.UpdateItem(input)
	util.CheckForError(err)
}

func connectDynamoDB() {
	if dynamodbsvc == nil {
		creds := getCredentials()
		cfg := aws.NewConfig().WithRegion(region).WithCredentials(creds)
		sess := session.Must(session.NewSession(cfg))
		// Create DynamoDB client
		dynamodbsvc = dynamodb.New(sess)
	}
}
