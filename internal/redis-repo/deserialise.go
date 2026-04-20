package redisrepo

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/kartikx04/chat/internal/models"
	"github.com/redis/go-redis/v9"
)

type Document struct {
	Id      uuid.UUID `json:"id"`
	Payload []byte    `json:"payload"`
	Total   int64     `json:"total"`
}

func Deserialise(res interface{}) []Document {
	switch v := res.(type) {
	case []interface{}:
		if len(v) > 1 {
			total := len(v) - 1
			docs := make([]Document, 0, total/2)

			for i := 1; i <= total; i += 2 {
				arrOfValues := v[i+1].([]interface{})
				value := arrOfValues[len(arrOfValues)-1].(string)

				doc := Document{
					Id:      v[i].(uuid.UUID),
					Payload: []byte(value),
					Total:   v[0].(int64),
				}

				docs = append(docs, doc)
			}

			return docs
		}
	default:
		log.Printf("different response type otherthan []interface{}. type: %T", res)
		return nil
	}

	return nil
}

func DeserialiseChat(docs []Document) []models.Chat {
	chats := []models.Chat{}
	for _, doc := range docs {
		var c models.Chat
		json.Unmarshal(doc.Payload, &c)

		c.Id = doc.Id
		chats = append(chats, c)
	}

	return chats
}

func DeserialiseContactList(contacts []redis.Z) []models.ContactList {
	contactList := make([]models.ContactList, 0, len(contacts))

	// improvement tip: use switch to get type of contact.Member
	// handle unknown type accordingly
	for _, contact := range contacts {
		contactList = append(contactList, models.ContactList{
			Username:     contact.Member.(string),
			LastActivity: int64(contact.Score),
		})
	}

	return contactList
}
