package websocket

type WSEvent struct {
	UserID string `json:"user_id"`
	Match bool `json:"match"`
}

var Events = make(chan WSEvent)