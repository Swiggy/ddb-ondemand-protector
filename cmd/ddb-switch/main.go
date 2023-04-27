package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"

	"github.com/Swiggy/ddb-ondemand-protector/config"
	ddb_switch "github.com/Swiggy/ddb-ondemand-protector/internal/ddb-switch"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	var cfg config.Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		logger.Fatal("Error in loading configuration", zap.Error(err))
	}
	handler, err := ddb_switch.NewDdbHandler(logger, svc, cfg)
	if err != nil {
		logger.Fatal("Error initializing DynamoDB handler", zap.Error(err))
	}
	lambda.Start(handler.Handle)
}
