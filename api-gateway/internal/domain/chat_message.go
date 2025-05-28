package domain

// ChatMessage stores data of various contents of Envelope.Data json
type ChatMessage struct {
	MessageId  	string `json:"message_id"`
	SenderId   	string `json:"sender_id"`
	Receiver_id string `json:"receiver_id"`
	Timestamp  	int64  `json:"timestamp"`
	Body       	string `json:"body"`
}


