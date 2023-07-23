package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	// "math/rand"
	// "time"

	"os"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func ReadFileandStoreRedis(redisClient redis.Client, statuscode string) {

	fileContent, err := ioutil.ReadFile("error.log")
	if err != nil {
		panic(err)
	}

	errorLog := ""

	if statuscode == "500" {
		fields := make(map[string]interface{})
		json.Unmarshal(fileContent, &fields)
		jsonby, _ := json.Marshal(fields)
		errorLog = string(jsonby)
	}
	if statuscode == "200" {
		errorLog = string(fileContent)
	}

	err = redisClient.HMSet(os.Getenv("CMS_CASE_ID"), map[string]interface{}{"message": string(errorLog), "statuscode": statuscode}).Err()
	if err != nil {
		fmt.Println(err)
	}

	// if statuscode == "500" {
	// 	fields := make(map[string]interface{})
	// 	json.Unmarshal(fileContent, &fields)
	// 	jsonby, _ := json.Marshal(fields)
	// 	errorLog = string(jsonby)
	// }
	// if statuscode == "200" {
	// 	errorLog = string(fileContent)
	// }

	// s := []string{"message", "statuscode"}

	// for _, j := range s {

	// 	keyValue := ""

	// 	if j == "message" {
	// 		keyValue = string(errorLog)
	// 		goto ENTRY

	// 	}

	// 	if j == "statuscode" {
	// 		keyValue = statuscode
	// 		goto ENTRY
	// 	}

	// ENTRY:

	// 	err = redisClient.HSet(os.Getenv("CMS_CASE_ID"), j, keyValue).Err()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }
}

func ErrorFormat(logger *zap.Logger, redisClient redis.Client, messagestring, statuscode string) {

	file, err := os.OpenFile("error.log", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	logger.Error(messagestring)
	ReadFileandStoreRedis(redisClient, statuscode)
}

func FileRemove(fileName, ExtractedFolderName string) {
	arr := []string{fileName, ExtractedFolderName}
	if len(arr) > 0 {
		for _, name := range arr {
			os.Remove(name)
		}
	}
}

// var letters = []rune("abcdefghijklmnopqrstuvwxyz")

// func RandSeq(n int) string {

// 	rand.Seed(time.Now().UnixNano())
// 	b := make([]rune, n)
// 	for i := range b {
// 		b[i] = letters[rand.Intn(len(letters))]
// 	}
// 	rand.Seed(time.Now().UnixNano())
// 	return string(b)
// }



    
  