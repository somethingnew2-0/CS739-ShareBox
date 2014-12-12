package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"client/settings"
)

func Get(o *settings.Options, address string) (map[string]interface{}, error) {
	log.Println("GET: ", address)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", settings.ServerAddress, address), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("USERID", o.UserId)
	req.Header.Set("AUTH", o.AuthToken)

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
	respObj := map[string]interface{}{}
	err = json.Unmarshal(respJson, &respObj)
	if err != nil {
		return nil, err
	}
	return respObj, nil
}

func Post(o *settings.Options, address string, values interface{}) (map[string]interface{}, error) {
	jsonStr, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}
	log.Println("POST: ", address, string(jsonStr))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", settings.ServerAddress, address), bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("USERID", o.UserId)
	req.Header.Set("AUTH", o.AuthToken)
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
	respObj := map[string]interface{}{}
	err = json.Unmarshal(respJson, &respObj)
	if err != nil {
		// TODO: Remove this
		ioutil.WriteFile("error.html", respJson, 0666)
		return nil, err
	}
	return respObj, nil
}
