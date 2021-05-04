package gitlab

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Gitlab struct {
	Host    string
	Project string
	Branch  string
	Token   string
}

type Commit struct {
	ID            string
	Branch        string
	CommitMessage string
	Actions       []CommitAction
}

type CommitAction struct {
	Action   string
	FilePath string
	Content  string
}

func (g *Gitlab) CommitFiles(actions []CommitAction, message string) error {
	commit := Commit{
		ID:            g.Project,
		Branch:        g.Branch,
		CommitMessage: message,
		Actions:       actions,
	}
	json, err := json.Marshal(commit)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(json)

	gitlabURL := "http://" + g.Host + "/api/v4/projects/" + g.Project + "/repository/commits"
	client := &http.Client{}
	req, err := http.NewRequest("POST", gitlabURL, body)
	req.Header.Add("PRIVATE-TOKEN", g.Token)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
