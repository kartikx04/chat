package controllers

import (
	"encoding/json"
	"fmt"
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

	u := &userReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, "error decoidng request object", http.StatusBadRequest)
		return
	}

	res := verifyContact(u.Username)
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

	u := r.URL.Query().Get("id")

	res := contactList(uuid.NewSHA1(uuid.NameSpaceURL, []byte(u)))
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

	fmt.Println(username1, username2)
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
		res.Message = "incorrect id"
		return res
	}

	if !redisrepo.IsUserExist(username) {
		res.Message = "incorrect id"
		return res
	}

	contactList, err := redisrepo.FetchContactList(id)
	if err != nil {
		log.Println("error in fetch contact list for id:", id, err)
		res.Message = "unable to fetch contact list. please try again later."
		return res
	}

	res.Status = true
	res.Data = contactList
	res.Total = len(contactList)
	return res
}
