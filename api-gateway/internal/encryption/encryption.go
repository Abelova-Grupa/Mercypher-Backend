package encryption

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// KeyBundle stores the public cryptographic information for a user
type KeyBundle struct {
	Username  string `json:"username"`
	PublicKey string `json:"public_key"` // Base64 SPKI
}

// WrappedKey represents a Group Key encrypted for a specific user
type WrappedKey struct {
	GroupID   string `json:"group_id"`
	Recipient string `json:"recipient_id"`
	Sender    string `json:"sender_id"`
	KeyData   string `json:"key_data"` // Base64 wrapped key
}

type KeyServer struct {
	mu           sync.RWMutex
	publicKeys   map[string]string       // username -> public_key
	groupSecrets map[string][]WrappedKey // group_id -> list of wrapped keys
}

func NewKeyServer() *KeyServer {
	return &KeyServer{
		publicKeys:   make(map[string]string),
		groupSecrets: make(map[string][]WrappedKey),
	}
}

// StartKeyServer initializes the routes and runs the server
func StartKeyServer(port string) {
	ks := NewKeyServer()

	mux := http.NewServeMux()

	// 1. Register/Update a user's Public Key
	mux.HandleFunc("/keys/register", ks.handleRegisterKey)

	// 2. Fetch a user's Public Key (to start DH)
	mux.HandleFunc("/keys/fetch", ks.handleFetchKey)

	// 3. Store wrapped group keys for members to pick up
	mux.HandleFunc("/keys/group/distribute", ks.handleDistributeGroupKey)

	// 4. Members fetch their wrapped key for a specific group
	mux.HandleFunc("/keys/group/fetch", ks.handleFetchGroupKey)

	fmt.Printf("Encryption Key Server running on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Printf("Key Server failed: %s\n", err)
	}
}

func (ks *KeyServer) handleRegisterKey(w http.ResponseWriter, r *http.Request) {
	var bundle KeyBundle
	if err := json.NewDecoder(r.Body).Decode(&bundle); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ks.mu.Lock()
	ks.publicKeys[bundle.Username] = bundle.PublicKey
	ks.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func (ks *KeyServer) handleFetchKey(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	ks.mu.RLock()
	pubKey, uuid := ks.publicKeys[username]
	ks.mu.RUnlock()

	if !uuid {
		http.Error(w, "User keys not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"public_key": pubKey})
}

func (ks *KeyServer) handleDistributeGroupKey(w http.ResponseWriter, r *http.Request) {
	var keys []WrappedKey
	if err := json.NewDecoder(r.Body).Decode(&keys); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ks.mu.Lock()
	for _, k := range keys {
		ks.groupSecrets[k.GroupID] = append(ks.groupSecrets[k.GroupID], k)
	}
	ks.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func (ks *KeyServer) handleFetchGroupKey(w http.ResponseWriter, r *http.Request) {
	groupID := r.URL.Query().Get("group_id")
	userID := r.URL.Query().Get("user_id")

	ks.mu.RLock()
	allKeys := ks.groupSecrets[groupID]
	ks.mu.RUnlock()

	for _, k := range allKeys {
		if k.Recipient == userID {
			json.NewEncoder(w).Encode(k)
			return
		}
	}

	http.Error(w, "No key found for this user in this group", http.StatusNotFound)
}
