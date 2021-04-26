package data

import (
	"bufio"
	"strings"
	"testing"
)

var cliExample = `ip access-list demo
   10 deny ip any host 178.165.72.177
   20 deny ip any host 195.206.105.217
   30 permit ip any host 89.234.157.254
   180 permit ip any any
`

var highestSeq = 180

var aclExample = ACL{
	Name: "demo",
	Actions: map[int]string{
		10:  "deny ip any host 178.165.72.177",
		20:  "deny ip any host 195.206.105.217",
		30:  "permit ip any host 89.234.157.254",
		180: "permit ip any any",
	},
}

func TestNewACLFromCLI(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader(cliExample))
	acl := NewACLFromCLI(scanner)
	aclLength := len(aclExample.Actions)
	if acl.Name != aclExample.Name {
		t.Errorf("Expeced %s to be demo but got %s\n", aclExample.Name, acl.Name)
	}
	if len(acl.Actions) != aclLength {
		t.Errorf("Should have gotten %d lines but got %d\n", aclLength, len(acl.Actions))
	}
	for seq, action := range aclExample.Actions {
		if acl.Actions[seq] != action {
			t.Errorf("Expected ACL line: %d to be: %s, but got: %s\n", seq, aclExample.Actions[seq], acl.Actions[seq])
		}
	}
	for seq, action := range acl.Actions {
		if aclExample.Actions[seq] != action {
			t.Errorf("Expected ACL line: %d to be: %s, but got: %s\n", seq, aclExample.Actions[seq], acl.Actions[seq])
		}
	}
}

func TestGetHighestSeq(t *testing.T) {
	seq := aclExample.GetHighestSeq()
	if seq != highestSeq {
		t.Errorf("Highest seq was : %d, but should have been %d", seq, highestSeq)
	}
}

func TestCopy(t *testing.T) {
	copyTest := aclExample.Copy()
	if copyTest.Name != aclExample.Name {
		t.Errorf("Name of copy is %s, but should be %s\n", copyTest.Name, aclExample.Name)
	}
	for k, v := range copyTest.Actions {
		if v != aclExample.Actions[k] {
			t.Errorf("Copied entry should be %d : %s, but got %s", k, aclExample.Actions[k], v)
		}
	}
}

func TestAppendAction(t *testing.T) {
	// test1 assumes an 'any any' entry at the end
	test1 := aclExample.Copy()
	oldHigh := test1.GetHighestSeq()
	oldAction := test1.Actions[oldHigh]
	newAction := "permit ip host 10.1.2.3 172.22.22.0/24"
	test1.AppendAction(newAction)
	if test1.GetHighestSeq() != oldHigh+ACL_INCREMENT {
		t.Errorf("After append the highest sequence should be: %d, but got %d", oldHigh+ACL_INCREMENT, test1.GetHighestSeq())
	}
	if test1.Actions[test1.GetHighestSeq()] != oldAction {
		t.Errorf("After append the last entry should be %s, but is %s", oldAction, test1.Actions[test1.GetHighestSeq()])
	}
	if test1.Actions[oldHigh] != newAction {
		t.Errorf("After append the new entry should be %s, but is %s", newAction, test1.Actions[oldHigh])
	}
	// Now remove any any and append
	newAction2 := "deny ip 10.1.1.0/24 10.1.2.0/24"
	delete(test1.Actions, test1.GetHighestSeq())
	test1.AppendAction(newAction2)
	if test1.GetHighestSeq() != oldHigh+ACL_INCREMENT {
		t.Errorf("After append the highest sequence should be: %d, but got %d", oldHigh+ACL_INCREMENT, test1.GetHighestSeq())
	}
	if test1.Actions[test1.GetHighestSeq()] != newAction2 {
		t.Errorf("After append the last entry should be %s, but is %s", newAction2, test1.Actions[test1.GetHighestSeq()])
	}

}

func TestGenerateAVD(t *testing.T) {

}
func TestGenerateCLI(t *testing.T) {
	gen := aclExample.GenerateCLI()
	if gen != cliExample {
		t.Errorf("GenerateCLI did not produce the expected output:\n%s expected:\n%s", gen, cliExample)
	}

}

func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
