CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY,
    sender_id TEXT NOT NULL,
    receiver_id TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Crucial for GetChatHistory performance
CREATE INDEX idx_messages_participants_time 
ON messages (sender_id, receiver_id, created_at DESC);