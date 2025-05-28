package domain

// User struct should be used for retreiving information about the user
// in context of getting a UserId (for sending messages) while knowing
// their username or email.
//
// Note: Only for logged (websocket) users!
type User struct {
	UserId		string	`json:"user_id"`	
	Username	string  `json:"username"`
	Email		string  `json:"email"`
}
