package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kartikx04/chat/internal/models"
	"github.com/redis/go-redis/v9"
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
		ChatKey(chatId),
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
		"$.created_at_unix", "AS", "created_at_unix", "NUMERIC", "SORTABLE",
	).Result()
	if err != nil && !strings.Contains(err.Error(), "Index already exists") {
		log.Println(err)
	}

	fmt.Println(res, err)
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
		return nil, err
	}

	// deserialise redis data to map
	data := Deserialise(res)

	// deserialise data map to chat
	chats := DeserialiseChat(data)
	return chats, nil
}

func RegisterNewUser(username, password string) error {
	// redis-cli
	// SYNTAX: SET key value
	// SET username password
	// register new username:password key-value pair
	err := redisClient.Set(context.Background(), username, password, 0).Err()
	if err != nil {
		log.Println("error while adding new user", err)
		return err
	}

	// redis-cli
	// SYNTAX: SADD key value
	// SADD users username
	err = redisClient.SAdd(context.Background(), UserSetKey(), username).Err()
	if err != nil {
		log.Println("error while adding user in set", err)
		// redis-cli
		// SYNTAX: DEL key
		// DEL username
		// drop the registered user
		redisClient.Del(context.Background(), username)

		return err
	}

	return nil
}

func IsUserExist(username string) bool {
	// redis-cli
	// SYNTAX: SISMEMBER key value
	// SISMEMBER users username

	return redisClient.SIsMember(context.Background(), UserSetKey(), username).Val()
}

func IsUserAuthentic(username, password string) error {
	// redis-cli
	// SYNTAX: GET key
	// GET username
	p := redisClient.Get(context.Background(), username).Val()

	if !strings.EqualFold(p, password) {
		return fmt.Errorf("invalid username or password")
	}

	return nil
}

// UpdateContactList add contact to username's contact list
// if not present or update its timestamp as last contacted
func UpdateContactList(id uuid.UUID, contact string) error {
	// using redis.Z{} globaly and not creating instance!
	zs := redis.Z{Score: float64(time.Now().Unix()), Member: contact}

	// redis-cli SCORE is always float or int
	// SYNTAX: ZADD key SCORE MEMBER
	// ZADD contacts:username 1661360942123 contact
	err := redisClient.ZAdd(context.Background(),
		contactListZKey(id),
		zs,
	).Err()

	if err != nil {
		log.Println("error while updating contact list. username: ",
			id, "contact:", contact, err)
		return err
	}

	return nil
}

func CreateFetchChatBetweenIndex() {
	res, err := redisClient.Do(context.Background(),
		"FT.CREATE",
		ChatIndexKey(),
		"ON", "JSON",
		"PREFIX", "1", "chat#",
		"SCHEMA", "$.from", "AS", "from", "TAG",
		"$.to", "AS", "to", "TAG",
		"$.timestamp", "AS", "timestamp", "NUMERIC", "SORTABLE",
	).Result()

	fmt.Println(res, err)
}

// FetchContactList of the user. It includes all the messages sent to and received by contact
// It will return a sorted list by last activity with a contact
func FetchContactList(id uuid.UUID) ([]models.ContactList, error) {
	zRangeArg := redis.ZRangeArgs{
		Key:   contactListZKey(id),
		Start: 0,
		Stop:  -1,
		Rev:   true,
	}

	// redis-cli
	// SYNTAX: ZRANGE key from_index to_index REV WITHSCORES
	// ZRANGE contacts:username 0 -1 REV WITHSCORES
	res, err := redisClient.ZRangeArgsWithScores(context.Background(), zRangeArg).Result()

	contactList := DeserialiseContactList(res)

	if err != nil {
		log.Println("error while fetching contact list. username: ",
			contactList[1], err)
		return nil, err
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
	return redisClient.Set(context.Background(),
		fmt.Sprintf("user:%s", id),
		username,
		0,
	).Err()
}

func SetIdLookup(username string, id uuid.UUID) error {
	return redisClient.Set(context.Background(),
		fmt.Sprintf("username:%s", username),
		id.String(),
		0,
	).Err()
}
