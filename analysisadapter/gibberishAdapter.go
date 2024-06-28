package analysisadapter

import (
	"encoding/json"
	"fmt"
	analysis "github.com/shitpostingio/analysis-commons/structs"
	"io/ioutil"
	"net/http"
	"time"
)

// GetGibberishValues checks if the input string is gibberish
func GetGibberishValues(toCheck string) (gibberish analysis.GibberishResponse, err error) {

	client := &http.Client{Timeout: time.Second * 30}
	endpoint := getGibberishEndpoint()
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		err = fmt.Errorf("GetGibberishValues: can't set up request for string %s: %s", toCheck, err)
		return gibberish, err
	}

	request.Header.Add(cfg.AuthorizationHeaderName, cfg.AuthorizationHeaderValue)
	request.Header.Add(cfg.GibberishInputHeaderName, toCheck)
	webResponse, err := client.Do(request)
	if err != nil {
		err = fmt.Errorf("GetGibberishValues: unable to perform request: %s", err)
		return gibberish, err
	}
	defer closeSafely(webResponse.Body)

	bodyResult, err := ioutil.ReadAll(webResponse.Body)
	if err != nil {
		err = fmt.Errorf("GetGibberishValues: error while reading response: %s", err.Error())
		return gibberish, err
	}

	err = json.Unmarshal(bodyResult, &gibberish)
	if err != nil {
		err = fmt.Errorf("GetGibberishValues: error while unmarshaling response: %s", err)
	}

	return gibberish, err

}
