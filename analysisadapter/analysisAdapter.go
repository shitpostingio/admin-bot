package analysisadapter

import (
	"encoding/json"
	"fmt"
	"github.com/shitpostingio/admin-bot/api"
	"github.com/shitpostingio/analysis-commons/structs"
	"io/ioutil"
	"net/http"
	"time"
)

func GetAnalysis(uniqueFileID, fileID string) (analysis structs.Analysis, err error) {

	file, err := api.GetTelegramFile(uniqueFileID, fileID)
	if err != nil {
		err = fmt.Errorf("GetAnalysis: unable to retrieve telegram file path: %s", err)
		return analysis, err
	}

	if file.FileSize > cfg.FileSizeThreshold {
		err = fmt.Errorf("GetAnalysis: file size too big: %d", file.FileSize)
		return analysis, err
	}

	endpoint := getAnalysisEndpoint(file.FileID, file.FileUniqueID)
	client := &http.Client{Timeout: time.Second * 30}
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		err = fmt.Errorf("GetAnalysis: can't set up request for media with fileID %s: %s", fileID, err)
		return analysis, err
	}

	request.Header.Add(cfg.AuthorizationHeaderName, cfg.AuthorizationHeaderValue)
	request.Header.Add(cfg.CallerAPIKeyHeaderName, botToken)
	request.Header.Add(cfg.FilePathHeaderName, file.FilePath)
	webResponse, err := client.Do(request)
	if err != nil {
		err = fmt.Errorf("GetAnalysis: unable to perform request: %s", err)
		return analysis, err
	}
	defer closeSafely(webResponse.Body)

	bodyResult, err := ioutil.ReadAll(webResponse.Body)
	if err != nil {
		err = fmt.Errorf("GetAnalysis: error while reading response: %s", err)
		return analysis, err
	}

	fmt.Println(string(bodyResult))

	err = json.Unmarshal(bodyResult, &analysis)
	if err != nil {
		err = fmt.Errorf("GetAnalysis: error while unmarshaling response: %s", err)
		return analysis, err
	}

	if analysis.FingerprintErrorString != "" {
		err = fmt.Errorf("GetAnalysis: %w: %s", FingerprintError, analysis.FingerprintErrorString)
		return analysis, err
	}

	if analysis.NSFWErrorString != "" {
		err = fmt.Errorf("GetAnalysis: %w: %s", NSFWError, analysis.NSFWErrorString)
	}

	return analysis, err

}
