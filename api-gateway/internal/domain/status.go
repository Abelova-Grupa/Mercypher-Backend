package domain

import "time"

// Status struct holds status update of a message.
// Valid statuses are:
//		0 -> UNKNOWN
//		1 -> DELIVERED
//		2 -> SEEN
// RecipientId could be fetched from message service,
// for one message is tied to one recipient (or more 
// for group chats) but it is faster to eliminate 
// that request.
type Status struct {
	MessageId	string		`json:"message_id"`
	RecipientId	string		`json:"recipient_id"`
	Status		uint8		`json:"status"`
	Timestamp	time.Time	`json:"timestamp"`
}
