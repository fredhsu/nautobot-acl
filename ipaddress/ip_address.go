package ipaddress

import "strings"

type IPAddresses struct {
	IPAddresses []IPAddressType
}

type IPAddressType struct {
	ID            string
	DNSName       string
	TenantGroupID string
	Status        string
	Family        int
	Address       string
	MaskLength    string
	Tag           []string
}

func (iptype *IPAddressType) GetIPWithMask() string {
	return iptype.Address + "/" + iptype.MaskLength
}

func (iptype *IPAddressType) GenerateIPFromAny(action string) string {
	if strings.HasSuffix(iptype.Address, "/32") {
		return action + " ip any host " + strings.TrimSuffix(iptype.Address, "/32")
	} else {
		return action + " ip any " + iptype.Address
	}

}
