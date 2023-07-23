package setcasedata

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func SetCaseData(filePath string)error {
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Println(file.Name())
			os.Setenv("CASE_FILE_NAME", file.Name())
			break
		}
	}

	return nil
}
