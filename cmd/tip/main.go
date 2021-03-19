package main

import (
	"context"
	"github.com/tenhan/tip/configs"
	"github.com/tenhan/tip/internal/tip"
	"os"
	"strings"
)

func main() {

	// init client
	client := tip.NewClient(configs.GetHandlerWrapperList())
	if len(os.Args) == 1 {
		client.ShowHelpInfo()
		return
	}

	// ready go
	if err := client.Run(context.Background(), strings.Join(os.Args[1:], " ")); err != nil {
		panic(err)
	}
}
