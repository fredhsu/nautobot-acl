package nautobotgql

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	ndata "github.com/fredhsu/nautobot-buildacl/pkg/data"
)

type Query struct {
	Query string `json:"query"`
}

type Response struct {
	Data map[string]interface{} `json:"data"`
}

type NautobotServer struct {
	Hostname  string
	Url       string
	Token     string
	VerifyTLS bool
}

func NewNautobotServer(verifyTLS bool, hostname, token string) NautobotServer {
	url := "https://" + hostname + "/api/graphql/"
	return NautobotServer{
		Hostname:  hostname,
		Token:     token,
		Url:       url,
		VerifyTLS: verifyTLS,
	}
}

func (n *NautobotServer) RunQuery(query Query) Response {
	jsonStr, err := json.Marshal(query)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !n.VerifyTLS},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", n.Url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Token "+n.Token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var data Response
	err = json.Unmarshal(b, &data)
	if err != nil {
		panic(err)
	}
	return data
}

func (n *NautobotServer) QueryIPAddresses(tag string) []ndata.IPAddressType {
	query := Query{Query: `query { ip_addresses(tag:"` + tag + `"){address dns_name}}`}
	data := n.RunQuery(query)
	foo := data.Data["ip_addresses"].([]interface{})
	ips := make([]ndata.IPAddressType, len(foo))
	for i, v := range foo {
		tmp := v.(map[string]interface{})
		ip := ndata.IPAddressType{
			Address: tmp["address"].(string),
			DNSName: tmp["dns_name"].(string),
		}
		ips[i] = ip
	}
	return ips
}
