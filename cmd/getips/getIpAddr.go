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
	"github.com/fredhsu/nautobot-buildacl/ipaddress"
	ndata "github.com/fredhsu/nautobot-buildacl/pkg/data"
	ngql "github.com/fredhsu/nautobot-buildacl/pkg/nautobotgql"
)

var token, tag, server, bfout, cliout, projectID, branch, gitlab string

func init() {
	flag.StringVar(&token, "token", "", "API token")
	flag.StringVar(&tag, "tag", "", "Tags to match")
	flag.StringVar(&server, "server", "", "Nautobot server")
	flag.StringVar(&bfout, "bfout", "", "Output file for Batfish policy")
	flag.StringVar(&cliout, "cliout", "", "Output file for ACL CLI")
	flag.StringVar(&projectID, "projectid", "", "Project ID for Gitlab")
	flag.StringVar(&branch, "branch", "", "Branch for Gitlab")
	flag.StringVar(&gitlab, "gitlab", "", "Gitlab server")
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var wh ndata.Webhook
	err := decoder.Decode(&wh)

	if err != nil {
		panic(err)
	}

	fmt.Println(wh.Data)
	ns := ngql.NewNautobotServer(false, server, token)

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
	if gitlab != "" {
		commitFiles(bfp.ToJSONString(), bfFilePath, cli, cliFilePath)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post called"}`))
}

// TODO get string of cli back
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

// TODO get string of bfp back
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

func commitFiles(bfp, bfFilePath, cli, cliFilePath string) {
	commit := ndata.Commit{
		ID:            projectID,
		Branch:        branch,
		CommitMessage: "Nautobot IP Changed",
	}
	actions := []ndata.CommitAction{}
	if bfFilePath != "" {
		actions = append(actions, ndata.CommitAction{
			Action:   "update",
			FilePath: bfFilePath,
			Content:  bfp,
		})
	}
	if cliFilePath != "" {
		actions = append(actions, ndata.CommitAction{
			Action:   "update",
			FilePath: cliFilePath,
			Content:  cli,
		})
	}
	commit.Actions = actions
	gitlabURL := "http://" + gitlab + "/api/v4/projects/" + projectID + "/repository/commits"
	client := &http.Client{}
	req, err := http.NewRequest("POST", gitlabURL, nil)
	req.Header.Add("PRIVATE-TOKEN", token)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	http.HandleFunc("/nautobot", WebhookHandler)

	fmt.Printf("Starting webhook receiver - nautobot server %s, token %s\n", server, token)
	log.Fatal(http.ListenAndServe(":8080", nil))

	// fmt.Println(acl.GenerateCLI())
}
