package ddbswitch

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

/*
GetBillingMode Identifies the billing mode of the table using provisioned throughput
*/
func GetBillingMode(describeTableOutput *dynamodb.DescribeTableOutput) BillingMode {
	if *(describeTableOutput.Table.ProvisionedThroughput.ReadCapacityUnits) == 0 && *(describeTableOutput.Table.ProvisionedThroughput.WriteCapacityUnits) == 0 {
		return PAY_PER_REQUEST
	}
	return PROVISIONED

}

/*
GetTableName - This function parses the cloudwatch event to get the table name
*/
func GetTableName(e json.RawMessage) (string, error) {
	event := &DdbScalingEvent{}
	err := json.Unmarshal(e, event)
	if err != nil {
		return "", err
	}
	return event.TableName, nil
}

/*
GetNewCloudWatchEvent returns a cloudwatch event to be created basis a tablename
*/
func GetNewCloudWatchEvent(tableName string) (json.RawMessage, error) {
	event := &DdbScalingEvent{
		TableName: tableName,
	}
	ddbEvent, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	return ddbEvent, nil
}
