package batfish

import (
	"encoding/json"

	"github.com/fredhsu/nautobot-buildacl/ipaddress"
)

type BatfishPolicy struct {
	Permit []BatfishEntry `json:"permit"`
}

type BatfishEntry struct {
	Name        string   `json:"name"`
	DstIPs      string   `json:"dstIps"`
	IPProtocols []string `json:"ipProtocols"`
	DstPorts    []string `json:"dstPorts"`
}

func (b *BatfishPolicy) AppendPermit(entry BatfishEntry) {
	b.Permit = append(b.Permit, entry)
}

func NewFromIPs(ips []ipaddress.IPAddressType) BatfishPolicy {
	bfp := BatfishPolicy{}
	for _, ip := range ips {
		bfp.AppendPermit(BatfishEntry{
			Name:        ip.DNSName,
			DstIPs:      ip.Address,
			IPProtocols: []string{"tcp"},
		})
	}
	return bfp
}

func (b *BatfishPolicy) ToJSONString() string {
	json, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(json)
}
