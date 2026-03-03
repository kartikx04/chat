package redisrepo

import (
	"fmt"

	"github.com/google/uuid"
)

func userSetKey() string {
	return "users"
}

func sessionKey(client string) string {
	return "session#" + client
}

func chatKey() string {
	return fmt.Sprintf("chat#%s", uuid.New().String())
}

func chatIndex() string {
	return "chat#idx"
}

func contactListKey(username string) string {
	return "contacts:" + username
}
