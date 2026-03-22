package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/kartikx04/chat/internal/models"
)

func CreateChat(c *models.Chat) (uuid.UUID, error) {
	chatId := NewChatID()
	fmt.Printf("chat id:%s\n", chatId)

	c.Id = chatId

	by, err := json.Marshal(c)
	if err != nil {
		return uuid.Nil, err
	}
	res, err := redisClient.Do(
		context.Background(),
		"JSON.SET",
		fmt.Sprint(chatKeyPrefix, chatId.String()),
		"$",
		string(by),
	).Result()

	if err != nil {
		log.Println("error while setting chat json: ", err)
		return uuid.Nil, err
	}

	log.Println("chat successfully set: ", res)

	return chatId, nil
}

func CreateChatIndex() {
	res, err := redisClient.Do(context.Background(),
		"FT.CREATE",
		ChatIndexKey(),
		"ON", "JSON",
		"PREFIX", "1", "chat#",
		"SCHEMA",
		"$.from_id", "AS", "from_id", "TAG",
		"$.to_id", "AS", "to_id", "TAG",
		"$.created_at", "AS", "created_at", "NUMERIC", "SORTABLE",
	).Result()
	if err != nil && !strings.Contains(err.Error(), "Index already exists") {
		log.Println(err)
	}

	fmt.Println(res, err)
}

func FetchChatBetween(id1, id2, fromTS, toTS string) ([]models.Chat, error) {
	query := fmt.Sprintf("@from_id:{%s|%s} @to_id:{%s|%s} @created_at:[%s %s]",
		id1, id2, id1, id2, fromTS, toTS)

	res, err := redisClient.Do(context.Background(),
		"FT.SEARCH",
		ChatIndexKey(),
		query,
		"SORTBY", "created_at", "DESC",
		"LIMIT", "0", "50",
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
