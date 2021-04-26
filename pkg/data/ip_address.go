package data

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
