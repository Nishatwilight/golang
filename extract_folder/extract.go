package extractfolder

import (
	"context"
	"fmt"

	"io"
	"os"
	"path/filepath"
	"s3_file/CMS_setcallback"
	"s3_file/helper"
	"s3_file/set_casedata"

	"github.com/codeclysm/extract"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func Extract(filename string, logger *zap.Logger, redisClient redis.Client) (bool,error) {

	fmt.Println(filename)

	extension := filepath.Ext(filename)

	f, err := os.Open(filename)
	if err != nil {
		helper.ErrorFormat(logger, redisClient, err.Error(),"500")
		return false, nil
	}
	defer f.Close()

	var archiveFile io.Reader
	archiveFile = f

	err = os.Mkdir("CASEDIRECTORY", 0755)
	if err != nil {
		helper.ErrorFormat(logger, redisClient, err.Error(),"500")
	}

	switch extension {
	case ".zip":
		err = extract.Zip(context.Background(), archiveFile, "CASEDIRECTORY/", nil)
		if err != nil {
			return true,err
		}
	case ".gz":
		err = extract.Gz(context.Background(), archiveFile, "CASEDIRECTORY/", nil)
		if err != nil {
			return true,err
		}
	case ".tar":
		err = extract.Tar(context.Background(), archiveFile, "CASEDIRECTORY/", nil)
		if err != nil {
			return true,err
		}
	}

	logger.Info(fmt.Sprintf("File extracted succesfully %s", filename))

	err = setcasedata.SetCaseData("CASEDIRECTORY/")
	if err != nil {
		helper.ErrorFormat(logger, redisClient, err.Error(),"500")
		return false,err
	}

	logger.Info(fmt.Sprintf("Stored case data to the CASEDIRECTORY"))

	cmssetcallback.SetCallBackURL()
	
	logger.Info("Successfully posted the callbackurl to the CMS emulator")

	helper.ReadFileandStoreRedis(redisClient,"200")

	fmt.Println("successfully completed at the end")

	return false, nil

}
