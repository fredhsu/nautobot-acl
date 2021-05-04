package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/fredhsu/nautobot-buildacl/acl"
	"github.com/fredhsu/nautobot-buildacl/batfish"
	"github.com/fredhsu/nautobot-buildacl/gitlab"
	"github.com/fredhsu/nautobot-buildacl/ipaddress"
	"github.com/fredhsu/nautobot-buildacl/nautobot"
)

var gitlabtoken, token, tag, server, bfout, cliout, projectID, branch, gitlabHost, bfpath, clipath string

func init() {
	flag.StringVar(&token, "token", "", "API token")
	flag.StringVar(&tag, "tag", "", "Tags to match")
	flag.StringVar(&server, "server", "", "Nautobot server")
	flag.StringVar(&bfout, "bfout", "", "Output file for Batfish policy")
	flag.StringVar(&cliout, "cliout", "", "Output file for ACL CLI")
	flag.StringVar(&bfpath, "bfpath", "", "Repo file for Batfish policy")
	flag.StringVar(&clipath, "clipath", "", "Repo file for ACL CLI")
	flag.StringVar(&projectID, "projectid", "", "Project ID for Gitlab")
	flag.StringVar(&branch, "branch", "", "Branch for Gitlab")
	flag.StringVar(&gitlabHost, "gitlab", "", "Gitlab server")
	flag.StringVar(&gitlabtoken, "gitlabtoken", "", "Gitlab token")
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var wh nautobot.Webhook
	err := decoder.Decode(&wh)

	if err != nil {
		panic(err)
	}

	log.Printf("Received webhook: \n %+v\n", wh.Data)

	ns := nautobot.NewNautobotServer(false, server, token)

	ips := ns.QueryIPAddresses(tag)
	acl := acl.NewPermitACLFromIPs("demo", ips)
	cli := acl.GenerateCLI()
	bfp := batfish.NewFromIPs(ips)
	if bfout != "" {
		writeBF(acl, ips)
	}
	if cliout != "" {
		writeACL(acl, ips)
	}
	if gitlabHost != "" {
		log.Printf("Commiting changes to gitlab")
		gl := gitlab.Gitlab{
			Host:    gitlabHost,
			Project: projectID,
			Branch:  branch,
			Token:   gitlabtoken,
		}
		bfAction := gitlab.CommitAction{
			Action:   "update",
			FilePath: bfpath,
			Content:  bfp.ToJSONString(),
		}
		cliAction := gitlab.CommitAction{
			Action:   "update",
			FilePath: clipath,
			Content:  cli,
		}
		actions := []gitlab.CommitAction{bfAction, cliAction}

		gl.CommitFiles(actions, "Nautobot IP change")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post called"}`))
}

func writeACL(acl acl.ACL, ips []ipaddress.IPAddressType) {
	f, err := os.Create(cliout)
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(acl.GenerateCLI())
	if err != nil {
		panic(err)
	}
	fmt.Printf("wrote ACL to file %s\n", cliout)
}

func writeBF(acl acl.ACL, ips []ipaddress.IPAddressType) {
	bfp := batfish.BatfishPolicy{}
	for _, ip := range ips {
		// TODO - shift to preprend to prioritize critical
		bfp.AppendPermit(batfish.BatfishEntry{
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
	fmt.Printf("Wrote Batfish to file %s\n", bfout)
}

func main() {
	flag.Parse()

	http.HandleFunc("/nautobot", WebhookHandler)

	fmt.Printf("Starting webhook receiver - nautobot server %s, token %s\n", server, token)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
