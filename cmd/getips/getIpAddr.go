package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	ndata "github.com/fredhsu/nautobot-buildacl/pkg/data"
	ngql "github.com/fredhsu/nautobot-buildacl/pkg/nautobotgql"
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
	ns := ngql.NewNautobotServer(false, server, token)
	ips := ns.QueryIPAddresses(tag)
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
