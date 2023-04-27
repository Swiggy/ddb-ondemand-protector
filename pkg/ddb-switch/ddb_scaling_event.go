package ddbswitch

/*
DdbScalingEvent is a struct which is used to create a cloudwatch event.
*/
type DdbScalingEvent struct {
	TableName string `json:"table_name"`
}
