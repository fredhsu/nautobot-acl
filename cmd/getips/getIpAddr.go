package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	ndata "github.com/fredhsu/nautobot-buildacl/pkg/data"
)

type Query struct {
	Query string `json:"query"`
}

type Response struct {
	Data map[string]interface{} `json:"data"`
}

var token, tag, server string

func init() {
	flag.StringVar(&token, "token", "", "API token")
	flag.StringVar(&tag, "tag", "", "Tags to match")
	flag.StringVar(&server, "server", "", "Nautobot server")
}

func main() {
	flag.Parse()
	query := Query{Query: `query { ip_addresses(tag:"` + tag + `"){address}}`}
	jsonStr, err := json.Marshal(query)

	url := "https://" + server + "/api/graphql/"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Token "+token)
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
	foo := data.Data["ip_addresses"].([]interface{})
	ips := make([]ndata.IPAddressType, len(foo))
	for i, v := range foo {
		tmp := v.(map[string]interface{})
		addr := tmp["address"].(string)
		ip := ndata.IPAddressType{Address: addr}
		ips[i] = ip
	}
	fmt.Println(ips)

	file, err := os.Open("../../configs/acls.eos")
	if err != nil {
		log.Fatalf("failed to open")

	}
	scanner := bufio.NewScanner(file)
	// scanner.Split(bufio.ScanLines)
	acl := ndata.NewACLFromCLI(scanner)
	fmt.Println(acl)
	cliacl := acl.GenerateCLI()
	fmt.Println(cliacl)
	acl.GenerateAVD()
}
