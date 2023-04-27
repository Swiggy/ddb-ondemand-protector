package ddbswitch

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"go.uber.org/zap"

	"github.com/Swiggy/ddb-ondemand-protector/config"
	ddb_switch "github.com/Swiggy/ddb-ondemand-protector/pkg/ddb-switch"
)

/*
Handler is a struct to handle the lambda events, which has a logger, the ddb service and config
*/
type Handler struct {
	logger *zap.Logger
	svc    dynamodbiface.DynamoDBAPI
	config config.Config
}

/*
NewDdbHandler is a function to instantiate a new handler to handle lambda events
*/
func NewDdbHandler(logger *zap.Logger, svc dynamodbiface.DynamoDBAPI, config config.Config) (*Handler, error) {
	return &Handler{logger, svc, config}, nil
}

/*
updateTableMode updates the table capacity billing mode.
For the switch from on-demand to provisioned, the provisioned RCU WCU are loaded from the configuration file.
*/
func (handler *Handler) updateTableMode(result *dynamodb.DescribeTableOutput) (*dynamodb.UpdateTableInput, error) {
	var updateInput *dynamodb.UpdateTableInput

	if ddb_switch.GetBillingMode(result) == ddb_switch.PAY_PER_REQUEST {
		rcuLimit, err := strconv.Atoi(handler.config.ProvisionedCapacityUnits.Rcu)
		if err != nil {
			return nil, err
		}
		wcuLimit, err := strconv.Atoi(handler.config.ProvisionedCapacityUnits.Wcu)
		if err != nil {
			return nil, err
		}
		updateInput = &dynamodb.UpdateTableInput{
			BillingMode: aws.String("PROVISIONED"),
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(int64(rcuLimit)),
				WriteCapacityUnits: aws.Int64(int64(wcuLimit)),
			},
			TableName: result.Table.TableName,
		}
	} else {
		updateInput = &dynamodb.UpdateTableInput{
			BillingMode: aws.String("PAY_PER_REQUEST"),
			TableName:   result.Table.TableName,
		}
	}
	_, err := handler.svc.UpdateTable(updateInput)
	if err != nil {
		return nil, err
	}

	return updateInput, nil
}

/*
Handle will be the handler receiving all the lambda events for the conversion of provisioned capacity
*/
func (handler *Handler) Handle(e json.RawMessage) error {
	/*
		The handler receives the cloudwatch event to switch the dynamodb capacity mode from provisioned to on-demand and
		vice-versa.
	*/
	handler.logger.Info("Processing DynamoDB scaling event", zap.String("details", string(e)))

	table, err := ddb_switch.GetTableName(e)
	if err != nil {
		return err
	}

	result, err := handler.svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	})
	if err != nil {
		return err
	}

	handler.logger.Info("Updating the billing mode of table", zap.String("table", *result.Table.TableName))
	updateInput, err := handler.updateTableMode(result)
	if err != nil {
		return err
	}
	handler.logger.Info("Billing Mode of the table has been updated successfully",
		zap.String("table", *result.Table.TableName),
		zap.String("billing", *updateInput.BillingMode))

	return nil
}
