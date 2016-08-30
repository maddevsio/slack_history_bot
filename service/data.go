package service

type IndexData struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Message   string `json:"message"`
	Channel   string `json:"channel"`
	Timestamp string `json:"timestamp"`
}
