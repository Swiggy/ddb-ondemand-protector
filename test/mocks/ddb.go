package mocks

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/mock"
)

/*
Mocks is a struct to define mock and a ddb interface to mock the ddb API
*/
type Mocks struct {
	mock.Mock
	dynamodbiface.DynamoDBAPI
}

/*
DescribeTable is a mock for dynamodb describe table function
*/
func (m *Mocks) DescribeTable(input *dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	args := m.Called(input)
	tableOutput := args.Get(0).(*dynamodb.DescribeTableOutput)
	return tableOutput, args.Error(1)
}

/*
UpdateTable is a mock for dynamodb update table function
*/
func (m *Mocks) UpdateTable(input *dynamodb.UpdateTableInput) (*dynamodb.UpdateTableOutput, error) {
	args := m.Called(input)
	return nil, args.Error(0)
}
