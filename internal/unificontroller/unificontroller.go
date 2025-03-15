package unificontroller

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/unpoller/unifi"
)

type Client struct {
	Name string
	Ip   string
}

func GetFixedIpClients(username string, password string, url string, siteName string) ([]Client, error) {
	fmt.Println("Fetching Unifi Clients")

	unifiConfig := &unifi.Config{
		User: username,
		Pass: password,
		URL:  url,
	}

	invalidChars := regexp.MustCompile(`\s|'|:|,|_|’`)

	uClient, err := unifi.NewUnifi(unifiConfig)
	if err != nil {
		return nil, err
	}
	sites, err := uClient.GetSites()
	if err != nil {
		return nil, err
	}
	fmt.Println("The Following Sites have been found configured on the Unifi Controller:")
	var targetSite *unifi.Site
	for _, site := range sites {
		fmt.Println(site.Name)
		fmt.Println(site.SiteName)
		if strings.HasPrefix(site.SiteName, siteName) {
			targetSite = site
		}
	}
	if targetSite == nil {
		fmt.Println("The target site was not found - Unable to continue")
		return nil, errors.New("target site not found")
	}
	fmt.Println("Target Site Found")
	fmt.Println("Fetching Clients")
	clients, err := uClient.GetUsers([]*unifi.Site{targetSite}, 87600)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%d Clients Fetched\n", len(clients))
	var fixedIps []Client
	for _, client := range clients {
		if client.UseFixedIp.Val {
			name := strings.ToLower(client.Name)
			if len(client.Note) != 0 {
				name = strings.ToLower(client.Note)
			}
			name = invalidChars.ReplaceAllString(name, "-")
			fixedIps = append(fixedIps, Client{Name: name, Ip: client.FixedIp})
		}
	}
	fmt.Printf("%d Fixed IP Clients Found\n", len(fixedIps))
	return fixedIps, nil
}
