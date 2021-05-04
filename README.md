# Nautobot-ACL
This tool will grab tagged IP addresses from Nautobot and create an ACL based on tags.
For instance if a number of IP addresses are tagged `critical` then it can create a list
of permit ACLs for all critical IPs.

Usage:
`getips -token f6df868dfa674ff1d5fdfaac169eda87a55d2d93 -server nautobot.sjc.aristanetworks.com -tag critical -bfout configs/permit2.json -cliout configs/acls.eos`


```getips -token f6df868dfa674ff1d5fdfaac169eda87a55d2d93 \
-server nautobot.sjc.aristanetworks.com \
-bfpath permit.json \
-branch nautobot \
-clipath /networks/configs/acl.cfg \
-gitlab dmz-gitlab.sjc.aristanetworks.com \
-projectid 5 \
-gitlabtoken NnnPwyihFTVRsnqk_dfi
```


# TODO
- [ ] Use Viper to do CLI and env
- [ ] Check for diff between current files and newly generated ones, only commit and push if there is a change
