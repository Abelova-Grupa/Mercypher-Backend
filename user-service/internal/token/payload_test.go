package token

import (
	"testing"
	"time"
)

func TestValid(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		wantErr   bool
	}{
		{"invalid payload", time.Now().Add(time.Hour), true},
		{"valid payload", time.Now().Add(-time.Hour), false},
	}

	for i, tt := range tests {
		payload := &Payload{ExpiresAt: tt.expiresAt}
		err := payload.Valid()
		if (err != nil) == tt.wantErr {
			t.Errorf("Payload validation test[%v] has failed, err: %v wantErr: %v", i, err, tt.wantErr)
		}
	}

}
