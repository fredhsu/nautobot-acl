package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	ndata "github.com/fredhsu/nautobot-buildacl/pkg/data"
	ngql "github.com/fredhsu/nautobot-buildacl/pkg/nautobotgql"
)

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

	// file, err := os.Open("../../configs/acls.eos")
	// if err != nil {
	// 	log.Fatalf("failed to open")

	// }
	// scanner := bufio.NewScanner(file)
	// acl := ndata.NewACLFromCLI(scanner)

	acl := ndata.NewACL("demo")
	bfp := ndata.BatfishPolicy{}
	for _, ip := range ips {
		bfp.AppendPermit(ndata.BatfishEntry{
			Name:        ip.DNSName,
			DstIPs:      ip.Address,
			IPProtocols: []string{"tcp"},
			// DstPorts:    []string{"80"},
		})
		acl.AppendAction(ip.GenerateIPFromAny("permit"))
	}
	json, err := json.MarshalIndent(bfp, "", "  ")
	err = ioutil.WriteFile("../../configs/permit2.json", json, 0644)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println(acl.GenerateCLI())
}
