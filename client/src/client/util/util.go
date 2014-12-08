package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Get(address string) (map[string]interface{}, error) {
	resp, err := http.Get(address)
	if err != nil {
		return nil, err
	}
	respJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respObj := map[string]interface{}{}
	err = json.Unmarshal(respJson, &respObj)
	if err != nil {
		return nil, err
	}
	return respObj, nil
}

func Post(address string, values url.Values) (map[string]interface{}, error) {
	resp, err := http.PostForm(address, values)
	if err != nil {
		return nil, err
	}
	respJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respObj := map[string]interface{}{}
	err = json.Unmarshal(respJson, &respObj)
	if err != nil {
		return nil, err
	}
	return respObj, nil
}
