package workwechat

type Message struct {
	MsgType    string           `json:"msgtype"`
	AgentID    int              `json:"agentid"`
	Text       *MessageContent  `json:"text,omitempty"`
	Markdown   *MessageContent  `json:"markdown,omitempty"`
	MarkdownV2 *MessageContent  `json:"markdown_v2,omitempty"`
	Textcard   *MessageTextCard `json:"textcard,omitempty"`
}

type MessageContent struct {
	Content string `json:"content"`
}

type MessageTextCard struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Btntxt      string `json:"btntxt"`
}
