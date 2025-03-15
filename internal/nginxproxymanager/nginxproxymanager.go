package nginxproxymanager

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AuthResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

type AuthRequest struct {
	Identity string `json:"identity"`
	Secret   string `json:"secret"`
}

type Meta struct {
	LetsencryptAgree bool    `json:"letsencrypt_agree"`
	DNSChallenge     bool    `json:"dns_challenge"`
	NginxOnline      bool    `json:"nginx_online"`
	NginxErr         *string `json:"nginx_err"`
}

type Config struct {
	ID                    int       `json:"id"`
	CreatedOn             time.Time `json:"created_on"`
	ModifiedOn            time.Time `json:"modified_on"`
	OwnerUserID           int       `json:"owner_user_id"`
	DomainNames           []string  `json:"domain_names"`
	ForwardHost           string    `json:"forward_host"`
	ForwardPort           int       `json:"forward_port"`
	AccessListID          int       `json:"access_list_id"`
	CertificateID         int       `json:"certificate_id"`
	SSLForced             bool      `json:"ssl_forced"`
	CachingEnabled        bool      `json:"caching_enabled"`
	BlockExploits         bool      `json:"block_exploits"`
	AdvancedConfig        string    `json:"advanced_config"`
	Meta                  Meta      `json:"meta"`
	AllowWebsocketUpgrade bool      `json:"allow_websocket_upgrade"`
	HTTP2Support          bool      `json:"http2_support"`
	ForwardScheme         string    `json:"forward_scheme"`
	Enabled               bool      `json:"enabled"`
	Locations             []string  `json:"locations"`
	HSTSEnabled           bool      `json:"hsts_enabled"`
	HSTSSubdomains        bool      `json:"hsts_subdomains"`
}

var authResponse AuthResponse

func auth(username string, password string, url string) (string, error) {
	fmt.Println("Authenticating to Nginx Proxy Manager")

	if authResponse.Token != "" && time.Now().Before(authResponse.Expires) {
		fmt.Println("Using cached token")
		return authResponse.Token, nil
	}

	// Create The Request
	bodyObject := AuthRequest{Identity: username, Secret: password}
	jsonBody, err := json.Marshal(bodyObject)
	if err != nil {
		return "", err
	}

	jsonBodyBytes := []byte(jsonBody)
	bodyReader := bytes.NewReader(jsonBodyBytes)

	req, err := http.NewRequest("POST", url+"/api/tokens", bodyReader)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	//Make The Request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", errors.New("Failed to authenticate to Nginx Proxy Manager")
	}

	// Parse The Response
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var content AuthResponse
	err = json.Unmarshal(resBody, &content)
	if err != nil {
		return "", err
	}

	fmt.Printf("Authenticated to Nginx Proxy Manager with token, expiring %s\n", content.Expires)

	authResponse = content

	return content.Token, nil
}

func GetProxyHosts(username string, password string, url string) ([]string, error) {

	token, err := auth(username, password, url)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url+"/api/nginx/proxy-hosts", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	//Make The Request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("Failed to get proxy hosts")
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var content []Config
	err = json.Unmarshal(resBody, &content)
	if err != nil {
		return nil, err
	}
	var hosts []string
	for _, host := range content {
		hosts = append(hosts, host.DomainNames...)
	}

	return hosts, nil
}
