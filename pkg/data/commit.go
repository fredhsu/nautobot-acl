package data

type Commit struct {
	ID string
	Branch string
	CommitMessage string
	Actions []CommitAction

}

type CommitAction struct {
	Action string
	FilePath string
	Content string
}