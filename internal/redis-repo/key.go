package redisrepo

import "github.com/google/uuid"

// Static keys
func UserSetKey() string {
	return "users"
}

func ChatIndexKey() string {
	return "chat#idx"
}

// Session key
func SessionKey(client string) string {
	return "session#" + client
}

// Contact list
func ContactListKey(id uuid.UUID) string {
	return "contacts:" + id.String()
}

// Chat keys
const chatKeyPrefix = "chat#"

func ChatKey(id uuid.UUID) string {
	return chatKeyPrefix + id.String()
}

// ID generator
func NewChatID() uuid.UUID {
	return uuid.New()
}

// Contact keys
func contactListZKey(id uuid.UUID) string {
	return "contacts:" + id.String()
}
