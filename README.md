# Nautobot-ACL
This tool will grab tagged IP addresses from Nautobot and create an ACL based on tags.
For instance if a number of IP addresses are tagged `critical` then it can create a list
of permit ACLs for all critical IPs.

Usage:
`getips -token f6df868dfa674ff1d5fdfaac169eda87a55d2d93 -server nautobot.sjc.aristanetworks.com -tag critical -bfout configs/permit2.json -cliout configs/acls.eos`


`getips -token f6df868dfa674ff1d5fdfaac169eda87a55d2d93 -server nautobot.sjc.aristanetworks.com -bfout configs/permit2.json`