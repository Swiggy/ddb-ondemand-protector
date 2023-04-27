package ddbswitch

import (
	"fmt"
	"github.com/Swiggy/ddb-ondemand-protector/config"
	ddb_switch "github.com/Swiggy/ddb-ondemand-protector/pkg/ddb-switch"
	"github.com/Swiggy/ddb-ondemand-protector/test/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

/*
TestSwitchToProvisioned function tests the basic functionality which updates the table billing mode from on-demand to provisioned.
*/
func TestSwitchToProvisioned(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ddbAPI := new(mocks.Mocks)
	cfg := config.Config{}
	cfg.ProvisionedCapacityUnits.Rcu = "10"
	cfg.ProvisionedCapacityUnits.Wcu = "10"
	handler, _ := NewDdbHandler(logger, ddbAPI, cfg)
	tableName := "testTable"
	cwEvent, _ := ddb_switch.GetNewCloudWatchEvent(tableName)
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}
	rcu := int64(0)
	wcu := int64(0)
	provisionedThroughputDescription := &dynamodb.ProvisionedThroughputDescription{ReadCapacityUnits: &rcu, WriteCapacityUnits: &wcu}
	output := &dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{
			ProvisionedThroughput: provisionedThroughputDescription,
			TableName:             aws.String(string(*input.TableName)),
		},
	}
	ddbAPI.On("DescribeTable", input).Return(output, nil)
	updateInput := &dynamodb.UpdateTableInput{
		BillingMode: aws.String("PROVISIONED"),
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(int64(10)),
			WriteCapacityUnits: aws.Int64(int64(10)),
		},
		TableName: aws.String(tableName),
	}
	ddbAPI.On("UpdateTable", updateInput).Return(nil, nil)
	err := handler.Handle(cwEvent)

	assert.NoError(t, err)
}

/*
TestSwitchToPayPerRequest function tests the basic functionality which updates the table billing mode from provisioned to on-demand.
*/
func TestSwitchToPayPerRequest(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ddbAPI := new(mocks.Mocks)
	cfg := config.Config{}
	cfg.ProvisionedCapacityUnits.Rcu = "10"
	cfg.ProvisionedCapacityUnits.Wcu = "10"
	handler, _ := NewDdbHandler(logger, ddbAPI, cfg)
	tableName := "testTable"
	cwEvent, _ := ddb_switch.GetNewCloudWatchEvent(tableName)
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}
	rcu := int64(10)
	wcu := int64(10)
	provisionedThroughputDescription := &dynamodb.ProvisionedThroughputDescription{ReadCapacityUnits: &rcu, WriteCapacityUnits: &wcu}
	output := &dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{
			ProvisionedThroughput: provisionedThroughputDescription,
			TableName:             aws.String(string(*input.TableName)),
		},
	}
	ddbAPI.On("DescribeTable", input).Return(output, nil)
	updateInput := &dynamodb.UpdateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		TableName:   aws.String(tableName),
	}
	ddbAPI.On("UpdateTable", updateInput).Return(nil, nil)
	err := handler.Handle(cwEvent)

	assert.NoError(t, err)
}

/*
TestErrorInDescribeTable function tests the error handling for the describe table function
In case any error is thrown by the api, the lambda should fail with the same error
*/
func TestErrorInDescribeTable(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ddbAPI := new(mocks.Mocks)
	cfg := config.Config{}
	cfg.ProvisionedCapacityUnits.Rcu = "10"
	cfg.ProvisionedCapacityUnits.Wcu = "10"
	handler, _ := NewDdbHandler(logger, ddbAPI, cfg)
	tableName := "testTable"
	cwEvent, err := ddb_switch.GetNewCloudWatchEvent(tableName)
	assert.Nil(t, err)
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}
	rcu := int64(10)
	wcu := int64(10)
	provisionedThroughputDescription := &dynamodb.ProvisionedThroughputDescription{ReadCapacityUnits: &rcu, WriteCapacityUnits: &wcu}
	output := &dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{
			ProvisionedThroughput: provisionedThroughputDescription,
			TableName:             aws.String(string(*input.TableName)),
		},
	}
	ddbAPI.On("DescribeTable", input).Return(output, fmt.Errorf("error describing table: oops"))
	err = handler.Handle(cwEvent)

	assert.Error(t, err)
	assert.Equal(t, "error describing table: oops", err.Error())
}

/*
TestErrorInUpdateTable function tests the error handling for the update table function
In case any error is thrown by the api, the lambda should fail with the same error
*/
func TestErrorInUpdateTable(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ddbAPI := new(mocks.Mocks)
	cfg := config.Config{}
	cfg.ProvisionedCapacityUnits.Rcu = "10"
	cfg.ProvisionedCapacityUnits.Wcu = "10"
	handler, _ := NewDdbHandler(logger, ddbAPI, cfg)
	tableName := "testTable"
	cwEvent, _ := ddb_switch.GetNewCloudWatchEvent(tableName)
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}
	rcu := int64(10)
	wcu := int64(10)
	provisionedThroughputDescription := &dynamodb.ProvisionedThroughputDescription{ReadCapacityUnits: &rcu, WriteCapacityUnits: &wcu}
	output := &dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{
			ProvisionedThroughput: provisionedThroughputDescription,
			TableName:             aws.String(string(*input.TableName)),
		},
	}
	ddbAPI.On("DescribeTable", input).Return(output, nil)
	updateInput := &dynamodb.UpdateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		TableName:   aws.String(tableName),
	}
	ddbAPI.On("UpdateTable", updateInput).Return(fmt.Errorf("error updating table: oops"))
	err := handler.Handle(cwEvent)

	assert.Error(t, err)
	assert.Equal(t, "error updating table: oops", err.Error())
}
