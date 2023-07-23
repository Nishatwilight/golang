package downloadfile

import (
	"fmt"
	"log"
	"os"
	extractfolder "s3_file/extract_folder"
	"s3_file/helper"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func DownloadFromS3(logger *zap.Logger, redisClient redis.Client) (bool,error) {

	fmt.Println(os.Getenv("SERVER_URL") + os.Getenv("CALLBACK_URL_ENDPOINT"))

	fileName := os.Getenv("COMPRESSED_CASE_FILE_NAME")
	BucketName := os.Getenv("S3_BUCKET_NAME")

	switch {
	case fileName == "":
		log.Fatal("FileName should not be empty")
	case BucketName == "":
		log.Fatal("BucketName should not be empty")
	}

	logger.Info(fmt.Sprintf("Received fileName %s, BucketName %s", fileName, BucketName))

	roleSess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		helper.ErrorFormat(logger, redisClient, err.Error(), "500")
		return false,err
	}

	// Create a STS client
	svc := sts.New(roleSess)
	roleToAssumeArn := "arn:aws:iam::489532684731:role/s2h-cmf-s3-read-write-access"
	sessionName := "session_test"
	externalId := "s2h-cmf-5ffdf1cd-c54e-41b3-a60b-7e74be162beb"

	result, err := svc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         &roleToAssumeArn,
		RoleSessionName: &sessionName,
		ExternalId:      &externalId,
	})

	if err != nil {
		helper.ErrorFormat(logger, redisClient, err.Error(), "500")
		return false,err
	}

	logger.Info("AssumeRole created successfully")

	file, err := os.Create(fileName)
	if err != nil {
		helper.ErrorFormat(logger, redisClient, err.Error(), "500")
		return false,err
	}

	s3sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			*result.Credentials.AccessKeyId,
			*result.Credentials.SecretAccessKey,
			*result.Credentials.SessionToken,
		),
	})

	if err != nil {
		helper.ErrorFormat(logger, redisClient, err.Error(), "500")
		return false,err
	}

	logger.Info("New session created successfully")

	fmt.Println("Downloading Log of " + fileName)

	fmt.Println()

	t1 := time.Now()

	fmt.Println("Starting Time =", time.Now().String())

	fmt.Println()

	downloader := s3manager.NewDownloader(s3sess, func(d *s3manager.Downloader) {
		d.PartSize = 5 * 1024 * 1024
		d.Concurrency = 100
	})

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(BucketName),
			Key:    aws.String(fileName),
		})

	if err != nil {
		helper.ErrorFormat(logger, redisClient, err.Error(), "500")
		return false,err
	}

	logger.Info(fmt.Sprintf("Case file downloaded successfully %s ", file.Name()))

	fmt.Println("Download completed", file.Name(), numBytes, "bytes")

	fi, _ := os.Stat(fileName)

	// get the size
	size := fi.Size()

	fmt.Println()

	fmt.Println("file size", size, numBytes, "bytes")

	fmt.Println()

	t2 := time.Now()

	fmt.Println("Ending Time =", time.Now().String())

	fmt.Println()

	t3 := t2.Sub(t1)
	
	fmt.Println("Total Time =", t3.String())

	checkextract,err := extractfolder.Extract(fileName, logger, redisClient)

	return checkextract,err

}

    