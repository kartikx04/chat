package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kartikx04/chat/internal/models"
)

func CreateChat(c *models.Chat) (string, error) {
	chatKey := chatKey()
	fmt.Printf("chat key:%s\n", chatKey)

	by, _ := json.Marshal(c)

	res, err := redisClient.Do(
		context.Background(),
		"JSON.SET",
		chatKey,
		"$",
		string(by),
	).Result()

	if err != nil {
		log.Println("error while setting chat json: ", err)
	}

	log.Println("chat successfully set: ", res)

	return chatKey, nil
}

func CreateFetchChatBetweenIndex() {
	res, err := redisClient.Do(context.Background(),
		"FT.CREATE",
		chatIndex(),
		"ON", "JSON",
		"PREFIX", "1", "chat#",
		"SCHEMA", "$.from", "AS", "from", "TAG",
		"$.to", "AS", "to", "TAG",
		"$.timestamp", "AS", "timestamp", "NUMERIC", "SORTABLE",
	).Result()

	fmt.Println(res, err)
}

func FetchChatBetween(id1, id2, fromTS, toTS string) ([]models.Chat, error) {
	// redis-cli
	// SYNTAX: FT.SEARCH index query
	// FT.SEARCH idx#chats '@from:{id2|id1} @to:{id1|id2} @timestamp:[0 +inf]'
	query := fmt.Sprintf("@from:{%s|%s} @to:{%s|%s} @timestamp:[%s %s]",
		id1, id2, id1, id2, fromTS, toTS)

	res, err := redisClient.Do(context.Background(),
		"FT.SEARCH",
		chatIndex(),
		query,
		"SORTBY", "timestamp", "DESC",
	).Result()

	if err != nil {
		return nil, err
	}

	// deserialise redis data to map
	data := Deserialise(res)

	// deserialise data map to chat
	chats := DeserialiseChat(data)
	return chats, nil
}
