package data

import (
	"bufio"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const ACL_INCREMENT = 10

// ACL defines an access list with the Actions map using the sequence number
// as a key, and the action as the string value. The Name is the name of the
// ACL. For example:
// ip access-list demo
//   10 permit ip any any
// Would have a Name of `demo` and the entry in the Actions map would be keyed with
// `10` and the value would be `permit ip any any`
type ACL struct {
	Name    string
	Actions map[int]string
}

// ACLAction is for future use to create more flexible and type checked ACLs
type ACLAction struct {
	Action string
}

func NewACLFromCLI(scanner *bufio.Scanner) ACL {
	scanner.Scan()
	aclTitle := scanner.Text()
	name := strings.Fields(aclTitle)[2]
	acl := ACL{
		Name:    name,
		Actions: map[int]string{},
	}
	for scanner.Scan() {
		ace := strings.TrimSpace(scanner.Text())
		line := strings.Fields(ace)
		rest := strings.TrimPrefix(ace, line[0]+" ")
		seq, err := strconv.Atoi(line[0])
		if err != nil {
			panic(err)
		}
		// action := ACLAction{Action: rest}
		acl.Actions[seq] = rest
	}
	return acl
}

func (acl *ACL) Copy() *ACL {
	newCopy := ACL{Name: acl.Name, Actions: make(map[int]string)}
	for k, v := range acl.Actions {
		newCopy.Actions[k] = v
	}
	return &newCopy
}

func (acl *ACL) GetHighestSeq() int {
	max := 0
	for k := range acl.Actions {
		if k > max {
			max = k
		}
	}
	return max
}
func (acl *ACL) GetLowestSeq() int {
	min := acl.GetHighestSeq()
	for k := range acl.Actions {
		if k < min {
			min = k
		}
	}
	return min
}

func (acl *ACL) AppendAction(action string) {
	newHigh := acl.GetHighestSeq() + ACL_INCREMENT
	// Check if hitting any any at the end, if so then increment and append
	// with sequence increased by 10
	if strings.Contains(acl.Actions[acl.GetHighestSeq()], "any any") {
		tmp := acl.Actions[acl.GetHighestSeq()]
		acl.Actions[acl.GetHighestSeq()] = action
		acl.Actions[newHigh] = tmp
	} else {
		// Otherwise just append
		acl.Actions[newHigh] = action
	}
}

func (acl *ACL) RemoveAction(seq int) error {
	if _, ok := acl.Actions[seq]; ok {
		delete(acl.Actions, seq)
		return nil
	} else {
		return fmt.Errorf("Sequence number %d not found", seq)
	}
}

func (acl *ACL) GenerateAVD() string {
	m := make(map[interface{}]interface{})
	action := make(map[string]string)
	action["action"] = "permit ip any any"
	ace := make(map[int]interface{})
	ace[10] = action
	ace[20] = action
	seqNum := make(map[string]interface{})
	seqNum["sequence_numbers"] = ace
	aclName := make(map[string]interface{})
	aclName["demo"] = seqNum
	m["standard_access_list"] = aclName
	fmt.Println(m)
	return ""
}

func (acl *ACL) GenerateCLI() string {
	output := "ip access-list " + acl.Name + "\n"

	keys := make([]int, len(acl.Actions))
	i := 0
	for k := range acl.Actions {
		keys[i] = k
		i++
	}
	sort.Ints(keys)

	for _, k := range keys {
		output += fmt.Sprintf("   %d %s\n", k, acl.Actions[k])
	}
	return output
}
