package controllers

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/kartikx04/chat/internal/database"
	"github.com/kartikx04/chat/internal/models"
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
	userId := r.URL.Query().Get("id")
	contactId := r.URL.Query().Get("contact")

	if userId == "" || contactId == "" {
		slog.WarnContext(r.Context(), "chat-history: missing params",
			"id", userId, "contact", contactId,
		)
		http.Error(w, "missing id or contact", http.StatusBadRequest)
		return
	}

	fromTS := "0"
	toTS := "+inf"
	if r.URL.Query().Get("from-ts") != "" && r.URL.Query().Get("to-ts") != "" {
		fromTS = r.URL.Query().Get("from-ts")
		toTS = r.URL.Query().Get("to-ts")
	}

	// Try Redis first
	chats, err := redisrepo.FetchChatBetween(userId, contactId, fromTS, toTS)
	if err != nil || len(chats) == 0 {
		slog.DebugContext(r.Context(), "chat-history: redis miss, trying postgres",
			"user", userId, "contact", contactId,
		)
		chats, err = fetchChatsFromDB(userId, contactId, 50)
		if err != nil {
			slog.ErrorContext(r.Context(), "chat-history: db fetch failed",
				"error", err,
				"user", userId,
				"contact", contactId,
			)
			http.Error(w, "failed to fetch history", http.StatusInternalServerError)
			return
		}
	}

	// Tag which messages belong to the requesting user
	for i := range chats {
		chats[i].IsSelf = chats[i].FromId.String() == userId
	}

	slog.DebugContext(r.Context(), "chat-history: returning messages",
		"count", len(chats),
		"user", userId,
		"contact", contactId,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}

func fetchChatsFromDB(userId, contactId string, limit int) ([]models.Chat, error) {
	var chats []models.Chat
	err := database.DB.
		Where(
			"(from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?)",
			userId, contactId, contactId, userId,
		).
		Order("created_at_unix ASC").
		Limit(limit).
		Find(&chats).Error
	return chats, err
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
