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
	handlerWrapperList = append(handlerWrapperList, handler.Wrapper{
		Name:        "百度翻译",
		Description: "From www.baidu.com",
		Handler:     dictionary.BaiduTranslate{},
	})
}

func GetHandlerWrapperList() []handler.Wrapper {
	return handlerWrapperList
}
