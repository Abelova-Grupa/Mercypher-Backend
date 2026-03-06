CREATE SCHEMA IF NOT EXISTS message_service;
CREATE TABLE IF NOT EXISTS message_service.chat_messages (
    message_id UUID PRIMARY KEY,
    sender_id UUID NOT NULL,
    receiver_id UUID NOT NULL,
    body TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
    
-- Crucial for GetChatHistory performance
CREATE INDEX IF NOT EXISTS idx_chat_messages_participants_time 
ON message_service.chat_messages (sender_id, receiver_id, timestamp DESC);