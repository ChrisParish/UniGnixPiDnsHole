package pihole

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	urlProcessor "net/url"
	"strings"
	"time"
)

type AuthSession struct {
	Valid    bool   `json:"valid"`
	TOTP     bool   `json:"totp"`
	SID      string `json:"sid"`
	CSRF     string `json:"csrf"`
	Validity int    `json:"validity"`
	Message  string `json:"message"`
}

type AuthResponse struct {
	Session AuthSession `json:"session"`
	Took    float64     `json:"took"`
	Expires time.Time
}

type AuthRequest struct {
	Password string `json:"password"`
}

var authResponse AuthResponse

func auth(url string, password string) error {
	if authResponse.Session.Valid && time.Now().Before(authResponse.Expires) {
		fmt.Println("Using cached session")
		return nil
	}
	fmt.Println("Authenticating to PiHole")

	// create the request
	bodyObject := AuthRequest{Password: password}
	jsonBody, err := json.Marshal(bodyObject)
	if err != nil {
		return err
	}
	jsonBodyBytes := []byte(jsonBody)
	bodyReader := strings.NewReader(string(jsonBodyBytes))

	req, err := http.NewRequest("POST", url+"/api/auth", bodyReader)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	// make the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("failed to authenticate to PiHole: %s", res.Status)
	}

	// parse the response
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var content AuthResponse
	err = json.Unmarshal(resBody, &content)
	if err != nil {
		return err
	}

	if !content.Session.Valid || content.Session.Message != "password correct" {
		return fmt.Errorf("failed to authenticate to PiHole: %s", content.Session.Message)
	}

	content.Expires = time.Now().Add(time.Duration(content.Session.Validity) * time.Second)
	authResponse = content

	return nil
}

func ClearAuth() {
	authResponse = AuthResponse{}
}

func addHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-FTL-SID", authResponse.Session.SID)
	req.Header.Add("X-FTL-CSRF", authResponse.Session.CSRF)
}

func GetLocalDns(url string, password string) ([]string, []string, error) {
	err := auth(url, password)
	if err != nil {
		return nil, nil, err
	}

	err = auth(url, password)
	if err != nil {
		return nil, nil, err
	}

	// create the request
	req, err := http.NewRequest("GET", url+"/api/config", nil)
	if err != nil {
		return nil, nil, err
	}
	addHeaders(req)

	// make the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if res.StatusCode != 200 {
		return nil, nil, fmt.Errorf("failed to get PiHole config: %s", res.Status)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	var content ConfigResponse
	err = json.Unmarshal(resBody, &content)
	if err != nil {
		return nil, nil, err
	}

	return content.Config.DNS.Hosts, content.Config.DNS.CnameRecords, nil

}

func AddLocalDns(url string, password string, host string, ip string) error {
	return localDns(url, password, host, ip, "PUT")
}

func RemoveLocalDns(url string, password string, host string, ip string) error {
	return localDns(url, password, host, ip, "DELETE")
}

func AddCname(url string, password string, host string, redirect string) error {
	return cname(url, password, host, redirect, "PUT")
}

func RemoveCname(url string, password string, host string, redirect string) error {
	return cname(url, password, host, redirect, "DELETE")
}

func cname(url string, password string, host string, redirect string, verb string) error {
	err := auth(url, password)
	if err != nil {
		return err
	}
	// create the request
	baseUrl, err := urlProcessor.Parse(url)
	if err != nil {
		return err
	}
	baseUrl.Path += "/api/config/dns/cnameRecords/"
	baseUrl.Path += fmt.Sprintf("%s,%s", host, redirect)

	req, err := http.NewRequest(verb, baseUrl.String(), nil)
	if err != nil {
		return err
	}
	addHeaders(req)

	// Make the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if verb == "PUT" && res.StatusCode != 201 && res.StatusCode != 400 { // 400 is returned if the record already exists
		return fmt.Errorf("failed to %s local CNAME record: %s", verb, res.Status)
	}
	if verb == "DELETE" && res.StatusCode != 204 {
		return fmt.Errorf("failed to %s local CNAME record: %s", verb, res.Status)
	}

	return nil
}

func localDns(url string, password string, host string, ip string, verb string) error {
	err := auth(url, password)
	if err != nil {
		return err
	}
	// create the request
	baseUrl, err := urlProcessor.Parse(url)
	if err != nil {
		return err
	}
	baseUrl.Path += "/api/config/dns/hosts/"
	baseUrl.Path += fmt.Sprintf("%s %s", ip, host)

	req, err := http.NewRequest(verb, baseUrl.String(), nil)
	if err != nil {
		return err
	}
	addHeaders(req)

	// Make the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if verb == "PUT" && res.StatusCode != 201 && res.StatusCode != 400 { // 400 is returned if the record already exists
		return fmt.Errorf("failed to %s local DNS record: %s", verb, res.Status)
	}
	if verb == "DELETE" && res.StatusCode != 204 {
		return fmt.Errorf("failed to %s local DNS record: %s", verb, res.Status)
	}

	return nil
}
