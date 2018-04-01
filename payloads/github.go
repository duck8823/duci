package payloads

type GitHubUser struct {
	Id      int64  `json:"id"`
	Login   string `json:"login"`
	Url     string `json:"url"`
	HtmlUrl string `json:"html_url"`
}

type GitHubCommitComment struct {
	Action string `json:"action"`
	Comment struct {
		Id        int64      `json:"id"`
		Url       string     `json:"url"`
		HtmlUrl   string     `json:"html_url"`
		User      GitHubUser `json:"user"`
		Position  int        `json:"position"`
		Line      int        `json:"line"`
		Path      string     `json:"path"`
		CommitId  string     `json:"commit_id"`
		CreatedAt string     `json:"created_at"`
		UpdatedAt string     `json:"updated_at"`
		Body      string     `json:"body"`
	} `json:"comment"`
	Repository struct {
		Id       int64      `json:"id"`
		FullName string     `json:"full_name"`
		Owner    GitHubUser `json:"owner"`
	} `json:"repository"`
}
