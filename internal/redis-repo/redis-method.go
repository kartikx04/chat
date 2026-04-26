package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kartikx04/chat/internal/database"
	"github.com/kartikx04/chat/internal/models"
	"github.com/redis/go-redis/v9"
)

func CreateChat(c *models.Chat) (uuid.UUID, error) {
	chatId := NewChatID()
	c.Id = chatId

	by, err := json.Marshal(c)
	if err != nil {
		return uuid.Nil, err
	}

	res, err := redisClient.Do(
		context.Background(),
		"JSON.SET",
		ChatKey(chatId),
		"$",
		string(by),
	).Result()
	if err != nil {
		slog.Error("redis: failed to set chat json", "chat_id", chatId, "error", err)
		return uuid.Nil, err
	}

	slog.Debug("redis: chat stored", "chat_id", chatId, "result", res)
	return chatId, nil
}

func CreateChatIndex() {
	_, err := redisClient.Do(context.Background(),
		"FT.CREATE",
		ChatIndexKey(),
		"ON", "JSON",
		"PREFIX", "1", "chat#",
		"SCHEMA",
		"$.from_id", "AS", "from_id", "TAG",
		"$.to_id", "AS", "to_id", "TAG",
		"$.created_at_unix", "AS", "created_at_unix", "NUMERIC", "SORTABLE",
	).Result()

	if err != nil {
		if strings.Contains(err.Error(), "Index already exists") {
			slog.Debug("redis: chat index already exists")
			return
		}
		slog.Error("redis: failed to create chat index", "error", err)
	} else {
		slog.Info("redis: chat index created")
	}
}

func FetchChatBetween(id1, id2, fromTS, toTS string) ([]models.Chat, error) {
	query := fmt.Sprintf(
		"@from_id:{%s|%s} @to_id:{%s|%s} @created_at_unix:[%s %s]",
		id1, id2, id1, id2, fromTS, toTS,
	)

	res, err := redisClient.Do(context.Background(),
		"FT.SEARCH",
		ChatIndexKey(),
		query,
		"SORTBY", "created_at_unix", "DESC",
		"LIMIT", "0", "50",
	).Result()
	if err != nil {
		slog.Error("redis: failed to fetch chats", "id1", id1, "id2", id2, "error", err)
		return nil, err
	}

	data := Deserialise(res)
	chats := DeserialiseChat(data)
	return chats, nil
}

func IsUserExist(username string) bool {
	var count int64
	database.DB.Model(&models.Users{}).Where("username = ?", username).Count(&count)
	return count > 0
}

func UpdateContactList(id uuid.UUID, contact string) error {
	zs := redis.Z{Score: float64(time.Now().Unix()), Member: contact}

	err := redisClient.ZAdd(context.Background(),
		contactListZKey(id),
		zs,
	).Err()
	if err != nil {
		slog.Error("redis: failed to update contact list",
			"user_id", id,
			"contact", contact,
			"error", err,
		)
		return err
	}

	return nil
}

func FetchContactList(id uuid.UUID) ([]models.ContactList, error) {
	zRangeArg := redis.ZRangeArgs{
		Key:   contactListZKey(id),
		Start: 0,
		Stop:  -1,
		Rev:   true,
	}

	res, err := redisClient.ZRangeArgsWithScores(context.Background(), zRangeArg).Result()
	if err != nil {
		slog.Error("redis: failed to fetch contact list", "user_id", id, "error", err)
		return nil, err
	}

	var contactList []models.ContactList
	for _, z := range res {
		contactId, err := uuid.Parse(z.Member.(string))
		if err != nil {
			slog.Warn("redis: invalid uuid in contact list", "member", z.Member, "error", err)
			continue
		}

		username, err := GetUsernameById(contactId)
		if err != nil {
			slog.Warn("redis: failed to resolve username", "contact_id", contactId, "error", err)
			continue
		}

		contactList = append(contactList, models.ContactList{
			Id:           contactId,
			Username:     username,
			LastActivity: int64(z.Score),
		})
	}

	return contactList, nil
}

func GetUsernameById(id uuid.UUID) (string, error) {
	return redisClient.Get(context.Background(),
		fmt.Sprintf("user:%s", id),
	).Result()
}

func GetIdByUsername(username string) (uuid.UUID, error) {
	val, err := redisClient.Get(context.Background(),
		fmt.Sprintf("username:%s", username),
	).Result()
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(val)
}

func SetUsernameLookup(id uuid.UUID, username string) error {
	err := redisClient.Set(context.Background(),
		fmt.Sprintf("user:%s", id),
		username,
		0,
	).Err()
	if err != nil {
		slog.Error("redis: failed to set username lookup", "user_id", id, "error", err)
	}
	return err
}

func SetIdLookup(username string, id uuid.UUID) error {
	if redisClient == nil {
		slog.Error("redis: client not initialized")
		return fmt.Errorf("redis client not initialized")
	}

	key := fmt.Sprintf("username:%s", username)

	if err := redisClient.Set(context.Background(), key, id.String(), 0).Err(); err != nil {
		slog.Error("redis: failed to set id lookup", "key", key, "error", err)
		return err
	}

	slog.Debug("redis: id lookup set", "username", username, "user_id", id)
	return nil
}
