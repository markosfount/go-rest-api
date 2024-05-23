package config

import "os"

const (
	BucketName  = "default"
	BrokerLink  = "localhost:29092"
	Topic       = "movies"
	SyncPublish = false
)

var ApiKey = os.Getenv("API_KEY")
