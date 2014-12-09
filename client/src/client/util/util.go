package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"client/settings"
)

func Get(address string) (map[string]interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", settings.ServerAddress, address))
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
	resp, err := http.PostForm(fmt.Sprintf("%s/%s", settings.ServerAddress, address), values)
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
