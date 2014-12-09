package util

import (
	"bytes"
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

func Post(address string, values interface{}) (map[string]interface{}, error) {
	jsonStr, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", settings.ServerAddress, address), bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	respJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// TODO: Remove this
	ioutil.WriteFile("error", respJson, 0666)
	respObj := map[string]interface{}{}
	err = json.Unmarshal(respJson, &respObj)
	if err != nil {
		return nil, err
	}
	return respObj, nil
}

func PostForm(address string, values url.Values) (map[string]interface{}, error) {
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
