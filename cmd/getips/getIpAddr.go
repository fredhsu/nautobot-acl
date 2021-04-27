package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	ndata "github.com/fredhsu/nautobot-buildacl/pkg/data"
	ngql "github.com/fredhsu/nautobot-buildacl/pkg/nautobotgql"
)

var token, tag, server, bfout, cliout string

func init() {
	flag.StringVar(&token, "token", "", "API token")
	flag.StringVar(&tag, "tag", "", "Tags to match")
	flag.StringVar(&server, "server", "", "Nautobot server")
	flag.StringVar(&bfout, "bfout", "", "Output file for Batfish policy")
	flag.StringVar(&cliout, "cliout", "", "Output file for ACL CLI")
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
	if bfout != "" {
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
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		err = ioutil.WriteFile(bfout, json, 0644)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}
	if cliout != "" {
		f, err := os.Create(cliout)
		if err != nil {
			panic(err)
		}
		_, err = f.WriteString(acl.GenerateCLI())
		if err != nil {
			panic(err)
		}
	}
	// fmt.Println(acl.GenerateCLI())
}
