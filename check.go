package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.konek.io/auth-server/controllers"
	"go.konek.io/auth-server/models"
)

func (a Auth) Check(domain, sid string) (bool, models.User, error) {
	var user controllers.InfosUserResponse
	var session controllers.CheckResponse

	// Checking session
	url := a.URL + "/session"

	request := controllers.CheckRequest{
		Token:  sid,
		Domain: domain,
	}
	b, err := json.Marshal(request)
	if err != nil {
		return false, user.Infos, err
	}
	buf := bytes.NewBuffer(b)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, buf)
	req.Header.Add("Content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return false, user.Infos, err
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, user.Infos, err
	}
	if resp.StatusCode != 200 {
		return false, user.Infos, errors.New("server returned an error : " + resp.Status)
	}
	err = json.Unmarshal(b, &session)
	if err != nil {
		return false, user.Infos, err
	}

	// Retrieving user
	url = a.URL + "/user/" + session.Session.UserID
	resp, err = http.Get(url)
	if err != nil {
		return false, user.Infos, err
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, user.Infos, err
	}
	if resp.StatusCode != 200 {
		fmt.Println("debug2:", url)
		return false, user.Infos, errors.New("server returned an error : " + resp.Status)
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return false, user.Infos, err
	}
	return true, user.Infos, nil
}
