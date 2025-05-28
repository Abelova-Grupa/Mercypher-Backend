package domain

import (
	"encoding/json"
)

// Envelope stores data as RawMessage so it can be unmarshaled according 
// to the given type.
type Envelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"` // defer decoding of data
}

