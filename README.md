# webhook-proxy

## Example

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/duck8823/webhook-proxy/payloads"
	"github.com/duck8823/webhook-proxy/proxy/handlers"
	"net/http"
)

func main() {
	http.Handle("/", &handlers.SlackNotificator{
		Url: "https://hooks.slack.com/services/XXXX/YYYY/ZZZZ",
		ConvertFunc: func(body []byte) (*payloads.SlackMessage, error) {
			commitComment := &payloads.GitHubCommitComment{}
			if err := json.Unmarshal(body, commitComment); err != nil {
				return nil, err
			}

			return &payloads.SlackMessage{
				Text: fmt.Sprintf("repository: %s", commitComment.Repository.FullName),
			}, nil
		},
	})

	http.ListenAndServe(":8080", nil)
}
```