package types 

import (
	"github.com/google/uuid"
)

type MessageRequest struct {
	Message string `json:"message"`
	Session_id uuid.UUID `json:"session_id"`
}