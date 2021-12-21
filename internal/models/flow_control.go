package models

import (
	"github.com/brenordv/go-request/internal/db"
	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"net/http"
	"sync"
)

type FlowControl struct {
	WaitGroup *sync.WaitGroup
	GuardChannel chan struct{}
	SessionId string
	Request *http.Request
	ProgressBar *progressbar.ProgressBar
	Db *db.DatabaseClient
}

func (fc *FlowControl) GenerateSessionId() {
	u := uuid.New()
	fc.SessionId = u.String()
}
