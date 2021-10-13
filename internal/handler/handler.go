package handler

import (
	"context"
	"time"
)

type Result struct {
	Title string
	Body  string
}
type Wrapper struct {
	Name        string
	Description string
	Handler     Handler
}
type HistoryObject struct {
	Keyword string `json:"keyword"`
	Results []Result `json:"results"`
	StartAt time.Time `json:"start_at"`
	EndAt time.Time `json:"end_at"`
	Duration int64 `json:"duration"`
}
type Handler interface {
	Handle(ctx context.Context, keyword string) (results []Result, err error)
}
