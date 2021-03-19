package handler

import (
	"context"
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
type Handler interface {
	Handle(ctx context.Context, keyword string) (results []Result, err error)
}
