package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/spf13/viper"
)

const cawsAccessKeyID = ""
const cawsSecretAccessKey = ""
const cawsSessionToken = ""

func getCredentials() *credentials.Credentials {
	viper.SetDefault("aws.awsAccessKeyID", cawsAccessKeyID)
	awsAccessKeyID := viper.GetString("aws.awsAccessKeyID")

	viper.SetDefault("aws.awsSecretAccessKey", cawsSecretAccessKey)
	awsSecretAccessKey := viper.GetString("aws.awsSecretAccessKey")

	viper.SetDefault("aws.awsSessionToken", cawsSessionToken)
	awsSessionToken := viper.GetString("aws.awsSessionToken")

	var creds *credentials.Credentials
	if awsAccessKeyID != "" && awsSecretAccessKey != "" {
		log.Print("Using config / static credentials")
		creds = credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, awsSessionToken)
	} else {
		log.Print("Using env / shared credentials")
		creds = credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{},
			})
	}
	_, err := creds.Get()
	if err != nil {
		log.Fatalf("bad credentials: %s", err)
	}
	if creds.IsExpired() {
		log.Fatal("credentials have expired")
	}
	return creds
}
