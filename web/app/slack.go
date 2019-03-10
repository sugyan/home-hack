package app

// Message type
type Message struct {
	ResponseType string        `json:"response_type,omitempty"`
	Channel      string        `json:"channel,omitempty"`
	UserName     string        `json:"username,omitempty"`
	IconEmoji    string        `json:"icon_emoji,omitempty"`
	Text         string        `json:"text"`
	Attachments  []*Attachment `json:"attachments"`
}

// Attachment type
type Attachment struct {
	AuthorIcon string `json:"author_icon"`
	AuthorName string `json:"author_name"`
	Text       string `json:"text"`
}
