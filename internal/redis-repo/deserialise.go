package redisrepo

type Document struct {
	ID      string `json:"id"`
	Payload []byte `json:"payload"`
	Total   int64  `json:"total"`
}
