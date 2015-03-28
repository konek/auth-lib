package auth

import (
	"bytes"
	"errors"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"go.konek.io/auth-server/controllers"
)

func (a Auth) Login(username string, password string, domain string) (ok bool, sid string, err error) {
	var loginRes controllers.LoginResponse

	url := a.URL + "/session"
	request := controllers.LoginRequest{
		Domain: domain,
		Username: username,
		Password: password,
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
	err = json.Unmarshal(b, &loginRes)
	if err != nil {
		return false, "", err
	}
	return true, loginRes.Session.Token, nil
}
