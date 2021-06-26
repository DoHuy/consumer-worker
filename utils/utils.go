package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func JsonToJson(obj1 interface{}, obj2 interface{}) error{
	raw, err := json.Marshal(obj1)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, obj2); err != nil {
		return err
	}
	return nil
}

func MakePOSTRequestAPI(url, body string) (string, error){
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	bodyResp, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(bodyResp))
	}
	return string(bodyResp), err
}