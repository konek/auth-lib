
package auth

import (
	"bytes"
	"errors"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"go.konek.io/auth-server/models"
	"go.konek.io/auth-server/controllers"
)

type Error struct{
	Code int `json:"code"`
	Msg string `json:"msg"`
}

func (e Error) Error() string {
	return e.Msg
}

func (e Error) StatusCode() int {
	return e.Code
}

func (a Auth) CreateUser(username string, password string, domains []string, variables map[string]interface{}) (ok bool, uid string, err error) {
	var createRes controllers.CreateResponse
	url := a.URL + "/user"
	request := models.User{
		Username: username,
		Password: password,
		Domains: domains,
		Enable: true,
		Variables: variables,
	}
	b, err := json.Marshal(request)
	if err != nil {
		return false, "", err
	}
	buf := bytes.NewBuffer(b)
	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		return false, "", err
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}
	if resp.StatusCode != 200 {
		var e Error
		err = json.Unmarshal(b, &e)
		if err != nil {
			return false, "", errors.New("An unexpected error occured")
		}
		return false, "", e
	}
	err = json.Unmarshal(b, &createRes)
	if err != nil {
		return false, "", err
	}
	return true, createRes.UserID, nil
}

