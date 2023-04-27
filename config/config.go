package config

/*
Config denotes the provisioned capacity units configuration imported from env variables.
*/
type Config struct {
	ProvisionedCapacityUnits struct {
		Rcu string `envconfig:"RCU"`
		Wcu string `envconfig:"WCU"`
	}
}