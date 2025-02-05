package telegram

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type Chat struct {
	ID int `json:"id"`
}

type User struct {
	UserName string `json:"username"`
}

type IncomingMessage struct {
	From User   `json:"from"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}
