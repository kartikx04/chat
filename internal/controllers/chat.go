package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	redisrepo "github.com/kartikx04/chat/internal/redis-repo"
)

type userReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Client   string `json:"client"`
}

type response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int         `json:"total,omitempty"`
}

func verifyContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := r.URL.Query().Get("username")
	log.Printf("verifyContactHandler: username=%s", username)

	if username == "" {
		log.Printf("verifyContactHandler: username is empty")
		json.NewEncoder(w).Encode(response{Status: false, Message: "username is required"})
		return
	}

	exists := redisrepo.IsUserExist(username)
	log.Printf("verifyContactHandler: IsUserExist returned %v for %s", exists, username)

	res := response{Status: exists}
	if !exists {
		res.Message = "user not found"
	}
	json.NewEncoder(w).Encode(res)
}

func chatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// user1 user2
	u1 := r.URL.Query().Get("u1")
	u2 := r.URL.Query().Get("u2")

	// chat between timerange fromTS toTS
	// where TS is timestamp
	// 0 to positive infinity
	fromTS, toTS := "0", "+inf"

	if r.URL.Query().Get("from-ts") != "" && r.URL.Query().Get("to-ts") != "" {
		fromTS = r.URL.Query().Get("from-ts")
		toTS = r.URL.Query().Get("to-ts")
	}

	res := chatHistory(u1, u2, fromTS, toTS)
	json.NewEncoder(w).Encode(res)
}

func contactListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	log.Printf("contactListHandler: id=%s", idStr)

	if idStr == "" {
		json.NewEncoder(w).Encode(response{Status: false, Message: "id is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("contactListHandler: invalid id format - %v", err)
		json.NewEncoder(w).Encode(response{Status: false, Message: "invalid id"})
		return
	}

	res := contactList(id)
	json.NewEncoder(w).Encode(res)
}

func verifyContact(username string) *response {
	// if invalid username and password return error
	// if valid user create new session
	res := &response{Status: true}

	status := redisrepo.IsUserExist(username)
	if !status {
		res.Status = false
		res.Message = "invalid username"
	}

	return res
}
func chatHistory(username1, username2, fromTS, toTS string) *response {
	// if invalid usernames return error
	// if valid users fetch chats
	res := &response{}

	// check if user exists
	if !redisrepo.IsUserExist(username1) || !redisrepo.IsUserExist(username2) {
		res.Message = "incorrect username"
		return res
	}

	id1, err := redisrepo.GetIdByUsername(username1)
	if err != nil {
		res.Message = "could not resolve user1"
		return res
	}

	id2, err := redisrepo.GetIdByUsername(username2)
	if err != nil {
		res.Message = "could not resolve user2"
		return res
	}

	chats, err := redisrepo.FetchChatBetween(id1.String(), id2.String(), fromTS, toTS)
	if err != nil {
		log.Println("error in fetch chat between", err)
		res.Message = "unable to fetch chat history. please try again later."
		return res
	}

	res.Status = true
	res.Data = chats
	res.Total = len(chats)
	return res
}

func contactList(id uuid.UUID) *response {
	res := &response{}

	username, err := redisrepo.GetUsernameById(id)
	if err != nil {
		log.Printf("contactList: GetUsernameById error for %s: %v", id, err)
		res.Message = "incorrect id"
		return res
	}

	// check if user exists
	if !redisrepo.IsUserExist(username) {
		log.Printf("contactList: IsUserExist returned false for %s", username)
		res.Message = "incorrect id"
		return res
	}

	contactList, err := redisrepo.FetchContactList(id)
	if err != nil {
		log.Printf("contactList: FetchContactList error: %v", err)
		res.Message = "unable to fetch contact list. please try again later."
		return res
	}

	res.Status = true
	res.Data = contactList
	res.Total = len(contactList)
	return res
}

func addContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	contact := r.URL.Query().Get("contact")

	log.Printf("addContactHandler: id=%s, contact=%s", id, contact)

	if id == "" || contact == "" {
		json.NewEncoder(w).Encode(response{Status: false, Message: "id and contact are required"})
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Printf("addContactHandler: invalid id - %v", err)
		json.NewEncoder(w).Encode(response{Status: false, Message: "invalid id"})
		return
	}

	// Verify contact username exists
	if !redisrepo.IsUserExist(contact) {
		log.Printf("addContactHandler: IsUserExist returned false for %s", contact)
		json.NewEncoder(w).Encode(response{Status: false, Message: "contact does not exist"})
		return
	}

	log.Printf("addContactHandler: IsUserExist passed, resolving username to ID")

	// Get contact's UUID from their username
	contactId, err := redisrepo.GetIdByUsername(contact)
	if err != nil {
		log.Printf("addContactHandler: GetIdByUsername error for %s - %v", contact, err)
		json.NewEncoder(w).Encode(response{Status: false, Message: "contact not found"})
		return
	}

	log.Printf("addContactHandler: resolved %s to %s", contact, contactId.String())

	// Add bidirectional
	err = redisrepo.UpdateContactList(uid, contactId.String())
	if err != nil {
		log.Printf("addContactHandler: UpdateContactList error (1) - %v", err)
		json.NewEncoder(w).Encode(response{Status: false, Message: "failed to add contact"})
		return
	}

	err = redisrepo.UpdateContactList(contactId, uid.String())
	if err != nil {
		log.Printf("addContactHandler: UpdateContactList error (2) - %v", err)
		json.NewEncoder(w).Encode(response{Status: false, Message: "failed to add contact"})
		return
	}

	log.Printf("addContactHandler: success, added contact %s for user %s", contact, id)
	json.NewEncoder(w).Encode(response{Status: true, Message: "contact added"})
}
