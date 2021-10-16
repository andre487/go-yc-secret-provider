package ycSecretProvider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type IamTokenData struct {
	AccessToken string `json:"access_token"`
}

var iamToken string

func GetIamToken() (string, error) {
	if len(iamToken) > 0 {
		return iamToken, nil
	}

	var err error
	if os.Getenv("YC_SECRET_MODE") == "prod" {
		iamToken, err = getProdIamToken()
	} else {
		iamToken, err = getDevIamToken()
	}

	return iamToken, err
}

func getDevIamToken() (string, error) {
	cmd := exec.Command("yc", "iam", "create-token")

	out := bytes.Buffer{}
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", newIamError(err)
	}

	return strings.TrimSpace(out.String()), nil
}

func getProdIamToken() (string, error) {
	metaServiceHost := os.Getenv("YC_METADATA_SERVICE")
	if len(metaServiceHost) == 0 {
		metaServiceHost = "169.254.169.254"
	}
	url := fmt.Sprintf("http://%s/computeMetadata/v1/instance/service-accounts/default/token", metaServiceHost)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", newIamError(err)
	}
	req.Header.Set("Metadata-Flavor", "Google")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", newIamError(err)
	}

	resultBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", newIamError(err)
	}

	var tokenData IamTokenData
	json.Unmarshal(resultBytes, &tokenData)

	if len(tokenData.AccessToken) == 0 {
		return "", errors.New("IAM token error: no IAM token")
	}
	return tokenData.AccessToken, nil
}

func newIamError(err error) error {
	return errors.New("IAM token error: " + err.Error())
}
