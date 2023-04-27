package ddbswitch

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetBillingModePayPerRequest(t *testing.T) {
	rcu := int64(0)
	wcu := int64(0)
	provisionedThroughputDescription := &dynamodb.ProvisionedThroughputDescription{ReadCapacityUnits: &rcu, WriteCapacityUnits: &wcu}
	output := &dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{
			ProvisionedThroughput: provisionedThroughputDescription,
			TableName:             aws.String("tableName"),
		},
	}
	billingMode := GetBillingMode(output)
	assert.Equal(t, billingMode, PAY_PER_REQUEST)
}

func TestGetBillingModeProvisioned(t *testing.T) {
	rcu := int64(10)
	wcu := int64(10)
	provisionedThroughputDescription := &dynamodb.ProvisionedThroughputDescription{ReadCapacityUnits: &rcu, WriteCapacityUnits: &wcu}
	output := &dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{
			ProvisionedThroughput: provisionedThroughputDescription,
			TableName:             aws.String("tableName"),
		},
	}
	billingMode := GetBillingMode(output)
	assert.Equal(t, billingMode, PROVISIONED)
}

func TestGetTableName(t *testing.T) {
	tableName := "testTable"
	cwEvent, _ := GetNewCloudWatchEvent(tableName)
	tableName, err := GetTableName(cwEvent)
	assert.Equal(t, tableName, "testTable")
	assert.NoError(t, err)
}

func TestGetTableNameError(t *testing.T) {
	_, err := GetTableName([]byte(nil))
	assert.Error(t, err)
}
