package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func doRequest(method, url string, body, rsp interface{}) (int, error) {
	requestBody := bytes.NewBuffer(nil)
	if body != nil {
		bodyJson, _ := json.Marshal(body)
		requestBody = bytes.NewBuffer(bodyJson)
	}

	request, _ := http.NewRequest(method, url, requestBody)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("url[%s] failed:%s\n", url, err.Error())
		return 0, err
	}
	defer response.Body.Close()

	rspBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("url[%s] read response body error:%s\n", url, err.Error())
		return response.StatusCode, err
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusAccepted {
		fmt.Printf("url[%s] not 200/202, response\n", url)
		fmt.Println(string(rspBody))
	}

	if rsp != nil {
		err = json.Unmarshal(rspBody, rsp)
		if err != nil {
			fmt.Println("unmarshal error:", err, string(rspBody))
			return response.StatusCode, err
		}

	}

	return response.StatusCode, nil
}

// DoGet 'rspReference' used to return the unmarshal response body
func DoGet(url string, rspReference interface{}) (int, error) {
	return doRequest(http.MethodGet, url, nil, rspReference)
}

// DoPost 'rspReference' used to return the unmarshal response body
func DoPost(url string, body, rspReference interface{}) (int, error) {
	return doRequest(http.MethodPost, url, body, rspReference)
}
