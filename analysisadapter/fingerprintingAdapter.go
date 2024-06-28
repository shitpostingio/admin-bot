package analysisadapter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/shitpostingio/admin-bot/api"
	analysis "github.com/shitpostingio/analysis-commons/structs"
)

// GetFingerprint gets the fingerprint values of a media given its file id
func GetFingerprint(uniqueFileID, fileID string) (fingerprint analysis.FingerprintResponse, err error) {

	file, err := api.GetTelegramFile(uniqueFileID, fileID)
	if err != nil {
		err = fmt.Errorf("GetFingerprint: unable to retrieve telegram file path: %s", err)
		return fingerprint, err
	}

	if file.FileSize > cfg.FileSizeThreshold {
		err = fmt.Errorf("GetFingerprint: file size too big: %d", file.FileSize)
		return fingerprint, err
	}

	endpoint := getFingerprintEndpoint(file.FileID, file.FileUniqueID)
	client := &http.Client{Timeout: time.Second * 30}
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		err = fmt.Errorf("GetFingerprint: can't set up request for media with fileID %s: %s", fileID, err)
		return fingerprint, err
	}

	request.Header.Add(cfg.AuthorizationHeaderName, cfg.AuthorizationHeaderValue)
	request.Header.Add(cfg.CallerAPIKeyHeaderName, botToken)
	request.Header.Add(cfg.FilePathHeaderName, file.FilePath)
	webResponse, err := client.Do(request)
	if err != nil {
		err = fmt.Errorf("GetFingerprint: unable to perform request: %s", err)
		return fingerprint, err
	}
	defer closeSafely(webResponse.Body)

	bodyResult, err := ioutil.ReadAll(webResponse.Body)
	if err != nil {
		err = fmt.Errorf("GetFingerprint: error while reading response: %s", err)
		return fingerprint, err
	}

	err = json.Unmarshal(bodyResult, &fingerprint)
	if err != nil {
		err = fmt.Errorf("GetFingerprint: error while unmarshaling response: %s", err)
		return fingerprint, err
	}

	return fingerprint, err

}
