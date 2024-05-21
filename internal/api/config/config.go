package config

import "os"

const BucketName = "default"

var ApiKey = os.Getenv("API_KEY")
