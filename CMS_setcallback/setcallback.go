package cmssetcallback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func SetCallBackURL() error {

	inputBody := map[string]string{
		"callbackUrl": fmt.Sprintf(os.Getenv("SERVER_URL") + os.Getenv("CALLBACK_URL_ENDPOINT")),
		"caseId":      os.Getenv("CMS_CASE_ID"),
		"authToken":   os.Getenv("CMS_AUTH_TOKEN"),
		"cmfApiKey":   os.Getenv("CMS_AUTH_TOKEN"),
		"sessionId":   os.Getenv("SESSION_ID"),
	}

	json_data, err := json.Marshal(inputBody)

	if err != nil {
		log.Println(err)
		return err
	}

	url := fmt.Sprintf(os.Getenv("CMS_URL") + "callback-url")

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		log.Println(err)
		return err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("appian-api-key", os.Getenv("CMS_AUTH_TOKEN"))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Something Went Wrong to set the call back URL")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return nil
}
