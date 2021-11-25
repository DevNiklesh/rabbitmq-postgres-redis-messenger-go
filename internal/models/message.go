package models

// Message is used to marsha/unmarshal
type Message struct {
	Text   string `json:"text"`
	Source string `json:"source"`
	Time   int64  `json:"time"`
}
