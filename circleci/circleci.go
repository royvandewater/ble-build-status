package circleci

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Build struct {
	Outcome string `json:"outcome"`
}

// GetLatestBuild will retrieve the latest build from CircleCI
func GetLatestBuild(username, project string) (*Build, error) {
	statusURL := fmt.Sprintf("https://circleci.com/api/v1.1/project/github/%s/%s", username, project)

	res, err := doJSONRequest(statusURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to get %v: %v", statusURL, err.Error())
	}

	if res.StatusCode != 200 {
		log.Fatalf("Expected 200, received %v status code back from GET \"%v\"\n", res.StatusCode, statusURL)
	}

	return parseLatestBuild(res.Body)
}

func doJSONRequest(requestURL string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	return http.DefaultClient.Do(req)
}

func parseLatestBuild(body io.ReadCloser) (*Build, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var builds []*Build

	err = json.Unmarshal(data, &builds)
	if err != nil {
		return nil, err
	}

	if len(builds) == 0 {
		return nil, fmt.Errorf("Found an empty list of builds")
	}

	return builds[0], nil
}
