package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	downloadfile "s3_file/download_file"
	"s3_file/error_log"
	"s3_file/helper"
	"s3_file/redisconfig"

	// "github.com/joho/godotenv"
)

func main() {

	logger := errorlog.ErrorLog()

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Println("Error loading .env file")
	// 	return
	// }

	redisClient, err := redisconfig.RedisConnect()
	if err != nil {
		log.Fatal(err)
	}

	downloadLog := []string{}

	for i := 0; i < 3; i++ {
		checkExtract, err := downloadfile.DownloadFromS3(logger, redisClient)
		if !checkExtract {
			message := fmt.Sprintf("Downloaded successfully in the first time %s and the caseId is %s ", os.Getenv("COMPRESSED_CASE_FILE_NAME"), os.Getenv("CMS_CASE_ID"))
			downloadLog = append(downloadLog, message)
			goto ENTRY
		} else {
			helper.FileRemove(os.Getenv("COMPRESSED_CASE_FILE_NAME"), "CASEDIRECTORY")
			downloadLog = append(downloadLog, err.Error()+"-"+string(i))
			helper.ErrorFormat(logger, redisClient, err.Error(), "500")
		}
	}

ENTRY:

	logforRedis, _ := json.Marshal(downloadLog)
	err = redisClient.HMSet(os.Getenv("CMS_CASE_ID")+":"+os.Getenv("SESSION_ID")+"-"+"download-log", map[string]interface{}{"message": string(logforRedis)}).Err()
	if err != nil {
		log.Println(err)
	}
}
