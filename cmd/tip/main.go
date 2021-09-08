package main

import (
	"context"
	"github.com/tenhan/tip/internal/tip"
	"os"
)

func main() {
	client := tip.NewClient(tip.GetHandlerWrapperList())
	if len(os.Args) == 1 {
		client.ShowHelpInfo()
		return
	}

	if err := client.Run(context.Background()); err != nil {
		panic(err)
	}
}