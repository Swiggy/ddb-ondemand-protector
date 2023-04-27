package ddbswitch

/*
BillingMode denotes the billing mode of dynamodb which can be provisioned or on-demand (pay per request)
*/
type BillingMode int

const (
	/*
		PROVISIONED capacity mode of dynamodb
	*/
	PROVISIONED BillingMode = iota
	/*
		PAY_PER_REQUEST denotes on-demand capacity mode of on-demand
	*/
	PAY_PER_REQUEST
)
