
package auth

import (
  "bytes"
  "errors"
  "net/http"
  "io/ioutil"
  "encoding/json"

  "gopkg.in/konek/auth-server.v1/controllers"
)

// Auth main structure
type Auth struct {
  URL string
}

// Auth sends a authentication request to the server.
//
// If authentication is successful, it returns ok = true, uid = <UserID>, err = nil.
//
// In any other cases, it returns ok = false, uid = "" and the appropriate error
func (a Auth) Auth(Domain, Username, Password string) (ok bool, uid string, err error) {
  var authRes controllers.AuthResponse

  url := a.URL + "/auth"
  request := controllers.AuthRequest{
    Domain: Domain,
    Username: Username,
    Password: Password,
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
    return false, "", errors.New("server returned an error : " + resp.Status)
  }
  err = json.Unmarshal(b, &authRes)
  if err != nil {
    return false, "", err
  }
  return true, authRes.UserID, nil
}

