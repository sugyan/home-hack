package app

const apiBaseURL = "https://slack.com/api"

type message struct {
	ResponseType string        `json:"response_type,omitempty"`
	Channel      string        `json:"channel,omitempty"`
	UserName     string        `json:"username,omitempty"`
	IconEmoji    string        `json:"icon_emoji,omitempty"`
	Text         string        `json:"text"`
	Attachments  []*attachment `json:"attachments"`
}

type attachment struct {
	AuthorIcon string `json:"author_icon"`
	AuthorName string `json:"author_name"`
	Text       string `json:"text"`
}

type apiResponse struct {
	OK       bool              `json:"ok"`
	Messages []*historyMessage `json:"messages"`
	HasMore  bool              `json:"has_more"`
	Error    string            `json:"error"`
}

type historyMessage struct {
	MessageID string      `json:"client_msg_id"`
	Text      string      `json:"text"`
	Reactions []*reaction `json:"reactions"`
}

type reaction struct {
	Name string `json:"name"`
}
