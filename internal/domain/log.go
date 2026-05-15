package domain

import "time"

type LogStatus string

const (
	LogStatusDone  LogStatus = "done"
	LogStatusError LogStatus = "error"
)

type Log struct {
	ID           int       `json:"id"`
	Filename     string    `json:"filename"`
	Status       LogStatus `json:"status"`
	NodeCount    int       `json:"node_count"`
	UploadedAt   time.Time `json:"uploaded_at"`
	ErrorMessage *string   `json:"error_message"`
}
