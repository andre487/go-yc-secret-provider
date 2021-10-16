package ycSecretProvider

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type LockBoxResult struct {
	Entries []struct {
		Key         string
		TextValue   string
		BinaryValue string
	}
	VersionId string
}

var lockBoxHandler = "https://payload.lockbox.api.cloud.yandex.net/lockbox/v1/secrets"

func GetLockBoxTextValue(secId string, name string) (string, error) {
	result, err := requestLockBox(secId)
	if err != nil {
		return "", errors.New("LockBox text value error: " + err.Error())
	}

	value := ""
	for _, val := range result.Entries {
		if val.Key == name {
			value = val.TextValue
			break
		}
	}

	if len(value) == 0 {
		return "", errors.New("LockBox text value error: empty secret value")
	}

	return value, nil
}

func GetLockBoxBinaryValue(secId string, name string) ([]byte, error) {
	result, err := requestLockBox(secId)
	if err != nil {
		return []byte{}, errors.New("LockBox binary value error: " + err.Error())
	}

	value := ""
	for _, val := range result.Entries {
		if val.Key == name {
			value = val.BinaryValue
			break
		}
	}

	if len(value) == 0 {
		return []byte{}, errors.New("LockBox binary value error: empty secret value")
	}

	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return []byte{}, errors.New("LockBox binary value error: " + err.Error())
	}

	return decoded, nil
}

func requestLockBox(secId string) (LockBoxResult, error) {
	iamToken, err := GetIamToken()
	if err != nil {
		return LockBoxResult{}, errors.New("LockBox request error: " + err.Error())
	}

	url := fmt.Sprintf("%s/%s/payload", lockBoxHandler, secId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LockBoxResult{}, errors.New("LockBox request error: " + err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+iamToken)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return LockBoxResult{}, errors.New("LockBox request error: " + err.Error())
	}

	resultBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return LockBoxResult{}, errors.New("LockBox request error: " + err.Error())
	}

	var result LockBoxResult
	json.Unmarshal(resultBytes, &result)
	return result, nil
}
