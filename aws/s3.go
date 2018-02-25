package aws

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/viper"
)

const cS3bucket = "songbeamer"

var pSession *session.Session
var sRegion string
var pSvc *s3.S3

// GetS3Files Download files
func GetS3Files() error {
	getS3Services()
	input := &s3.ListObjectsInput{
		Bucket:  aws.String(getBucket()),
		MaxKeys: aws.Int64(4),
		Prefix:  aws.String("songs/"),
	}

	log.Print("get objects")
	result, err := pSvc.ListObjects(input)
	log.Print("have objects")

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err
	}

	fmt.Println(result)
	for _, key := range result.Contents {
		fmt.Println(*key.Key)
		ob := &s3.GetObjectInput{
			Bucket: aws.String(getBucket()),
			Key:    aws.String(*key.Key),
		}
		object, err := pSvc.GetObject(ob)
		if err == nil {
			fmt.Print(reflect.TypeOf(object.Metadata))
		}
		fmt.Printf("%+v\n", object.Metadata)
	}
	return nil
}

func getS3Services() *s3.S3 {
	if pSvc == nil {
		pSvc = s3.New(getSession())
	}
	return nil
}

func getRegionForBucket(creds *credentials.Credentials) string {
	if sRegion == "" {
		cfg := aws.NewConfig().WithCredentials(creds)
		sess := session.Must(session.NewSession(cfg))
		ctx := context.Background()
		bucket := getBucket()
		// can be any region to connect initially to
		region, err := s3manager.GetBucketRegion(ctx, sess, bucket, "us-west-1")
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
				panic(fmt.Errorf("unable to find bucket %s's region not found", bucket))
			}
		}
		log.Printf("bucket %v is in using region %v", bucket, region)
		sRegion = region
	}
	return sRegion
}

func getSession() *session.Session {
	if pSession == nil {
		creds := getCredentials()
		region := getRegionForBucket(creds)
		cfg := aws.NewConfig().WithRegion(region).WithCredentials(creds)
		pSession = session.Must(session.NewSession(cfg))
	}
	return pSession
}

func getBucket() string {
	viper.SetDefault("uploader.s3bucket", cS3bucket)
	bucket := viper.GetString("uploader.s3bucket")
	if bucket == "" {
		log.Fatal("invalid bucket")
	}
	return bucket
}
