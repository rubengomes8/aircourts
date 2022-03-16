package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func HTTPGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err.Error())
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading GET response: %v", err.Error())
	}

	return body, nil
}
