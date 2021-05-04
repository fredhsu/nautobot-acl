package gitlab

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Gitlab struct {
	Host    string `json:"host"`
	Project string `json:"project"`
	Branch  string `json:"branch"`
	Token   string `json:"token"`
}

type Commit struct {
	ID            string         `json:"id"`
	Branch        string         `json:"branch"`
	CommitMessage string         `json:"commit_message"`
	Actions       []CommitAction `json:"actions"`
}

type CommitAction struct {
	Action   string `json:"action"`
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

func (g *Gitlab) CommitFiles(actions []CommitAction, message string) error {
	commit := Commit{
		ID:            g.Project,
		Branch:        g.Branch,
		CommitMessage: message,
		Actions:       actions,
	}
	log.Printf("Commiting: %+v", commit)
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
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	log.Printf("Commit response:\n %s", responseBody)
	return nil
}
