package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type responseBody any

func DoRequest[T responseBody](method string, path string, queryParams map[string]string) (T, error) {
	var rspBody T

	queryParamsStr := buildQueryParams(queryParams)
	url := fmt.Sprintf("%s%s%s", os.Getenv("API_URL"), path, queryParamsStr)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return rspBody, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-cg-demo-api-key", os.Getenv("API_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return rspBody, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return rspBody, err
	}

	rspBody, err = parseResponse[T](body)
	return rspBody, err
}

func buildQueryParams(queryParams map[string]string) string {
	queryParamsStr := ""
	if len(queryParams) > 0 {
		queryParamsStr += "?"
		for key, value := range queryParams {
			queryParamsStr += fmt.Sprintf("%s=%s&", key, value)
		}
	}
	return queryParamsStr
}

func parseResponse[T responseBody](rspBodyByte []byte) (T, error) {
	var rspBody T
	err := json.Unmarshal(rspBodyByte, &rspBody)
	if err != nil {
		return rspBody, err
	}
	return rspBody, nil
}
