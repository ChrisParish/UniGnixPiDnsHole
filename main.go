package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"unipidns/internal/nginxproxymanager"
	"unipidns/internal/pihole"
	"unipidns/internal/unificontroller"
)

type Config struct {
	Unifi             *Unifi             `json:"unifi"`
	PiHole            []PiHole           `json:"pihole"`
	NginxProxyManager *NginxProxyManager `json:"nginxProxyManager"`
	Domain            string             `json:"domain"`
	WebEdge           string             `json:"webEdge"`
	Local             string             `json:"local"`
}

type Unifi struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Url      string `json:"url"`
	Site     string `json:"site"`
}

type PiHole struct {
	Url      string `json:"url"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type NginxProxyManager struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("****UniPiDns****")
	fmt.Println()
	rawConfig, err := os.ReadFile("config.json")
	check(err)
	fmt.Println("Config file read successfully")
	var config *Config
	err = json.Unmarshal(rawConfig, &config)
	check(err)
	fmt.Println()
	fmt.Printf("Unifi Controller Url: %s\n", config.Unifi.Url)
	fmt.Printf("Nginx Proxy Manager Url: %s\n", config.NginxProxyManager.Url)
	for _, pihole := range config.PiHole {
		fmt.Printf("PiHole %s Url: %s\n", pihole.Name, pihole.Url)
	}
	fmt.Println()
	fixedIps, err := unificontroller.GetFixedIpClients(config.Unifi.Username, config.Unifi.Password, config.Unifi.Url, config.Unifi.Site)

	for idx, client := range fixedIps {
		fixedIps[idx].Name = fmt.Sprintf("%s.%s", strings.ToLower(client.Name), config.Local)
	}

	var fixedIpClients []string
	for _, client := range fixedIps {
		fixedIpClients = append(fixedIpClients, fmt.Sprintf("%s %s", client.Ip, client.Name))
	}
	fmt.Printf("%d Fixed IP Clients Found\n", len(fixedIpClients))

	check(err)

	hosts, err := nginxproxymanager.GetProxyHosts(config.NginxProxyManager.Username, config.NginxProxyManager.Password, config.NginxProxyManager.Url)
	check(err)
	var cnameHosts []string
	for _, host := range hosts {
		if strings.HasSuffix(host, config.Domain) {
			cnameHosts = append(cnameHosts, host)
		}
	}
	for idx, cnameHost := range cnameHosts {
		cnameHosts[idx] = fmt.Sprintf("%s,%s.%s", cnameHost, config.WebEdge, config.Local)
	}

	fmt.Printf("%d CNAME Hosts Found\n", len(cnameHosts))

	fmt.Println()
	for _, piHoleConfig := range config.PiHole {
		fmt.Println("Processing DNS on PiHole: " + piHoleConfig.Name)
		aRecords, cnames, err := pihole.GetLocalDns(piHoleConfig.Url, piHoleConfig.Password)
		check(err)
		fmt.Printf("	%d A Records Found\n", len(aRecords))
		fmt.Printf("	%d CNAME Records Found\n", len(cnames))

		check(err)
		var hostsToAdd []string
		var cnamesToAdd []string
		var hostsToRemove []string
		var cnamesToRemove []string

		// Add Fixed IP Clients
		for _, clientIp := range fixedIpClients {
			if !slices.Contains(aRecords, clientIp) {
				hostsToAdd = append(hostsToAdd, clientIp)
			}

		}
		// Remove Fixed IP Clients
		for _, aRecord := range aRecords {
			if !slices.Contains(fixedIpClients, aRecord) {
				hostsToRemove = append(hostsToRemove, aRecord)
			}
		}
		// Add CNAME Hosts
		for _, cnameHost := range cnameHosts {
			if !slices.Contains(cnames, cnameHost) {
				cnamesToAdd = append(cnamesToAdd, cnameHost)
			}
		}
		// Remove CNAME Hosts
		for _, cname := range cnames {
			if !slices.Contains(cnameHosts, cname) {
				cnamesToRemove = append(cnamesToRemove, cname)
			}
		}

		fmt.Printf("	%d A Records to Add\n", len(hostsToAdd))
		fmt.Printf("	%d A Records to Remove\n", len(hostsToRemove))
		fmt.Printf("	%d CNAME Records to Add\n", len(cnamesToAdd))
		fmt.Printf("	%d CNAME Records to Remove\n", len(cnamesToRemove))

		for _, host := range hostsToAdd {
			err = pihole.AddLocalDns(piHoleConfig.Url, piHoleConfig.Password, strings.Split(host, " ")[1], strings.Split(host, " ")[0])
			check(err)
		}

		for _, host := range hostsToRemove {
			err = pihole.RemoveLocalDns(piHoleConfig.Url, piHoleConfig.Password, strings.Split(host, " ")[1], strings.Split(host, " ")[0])
			check(err)
		}

		for _, cname := range cnamesToAdd {
			err = pihole.AddCname(piHoleConfig.Url, piHoleConfig.Password, strings.Split(cname, ",")[0], strings.Split(cname, ",")[1])
			check(err)
		}

		for _, cname := range cnamesToRemove {
			err = pihole.RemoveCname(piHoleConfig.Url, piHoleConfig.Password, strings.Split(cname, ",")[0], strings.Split(cname, ",")[1])
			check(err)
		}

	}
}
