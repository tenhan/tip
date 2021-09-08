package tip

import (
	"github.com/tenhan/tip/internal/handler"
	"github.com/tenhan/tip/internal/handler/dictionary"
)

var handlerWrapperList []handler.Wrapper

func init() {
	handlerWrapperList = append(handlerWrapperList, handler.Wrapper{
		Name:        "Jargon",
		Description: "Show the meaning of jargon.",
		Handler:     dictionary.Jargon{},
	})
}

func GetHandlerWrapperList() []handler.Wrapper {
	return handlerWrapperList
}
