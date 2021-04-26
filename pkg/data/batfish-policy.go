package data

type BatfishPolicy struct {
	Permit []BatfishEntry
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