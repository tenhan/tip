package tip

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/tenhan/tip/configs"
	"github.com/tenhan/tip/internal/handler"
	"sync"
)

const VersionName string = "1.0.0"

type Client struct {
	handlerWrapperList []handler.Wrapper
	resultList         []handler.Result
	resultLock         sync.Mutex
}

func NewClient(wrappers []handler.Wrapper) *Client {
	return &Client{
		handlerWrapperList: wrappers,
	}
}

// SetSpiderWrapper
func (c *Client) SetSpiderWrapper(wrappers []handler.Wrapper) (err error) {
	c.handlerWrapperList = wrappers
	return
}

// Run
func (c *Client) Run(ctx context.Context, keyword string) (err error) {
	if c.handlerWrapperList == nil || len(c.handlerWrapperList) == 0 {
		return fmt.Errorf("no spider found")
	}
	group := sync.WaitGroup{}
	group.Add(len(c.handlerWrapperList))
	for _, s := range c.handlerWrapperList {
		go func(sp handler.Wrapper) {
			defer group.Done()
			r, err := sp.Handler.Handle(ctx, keyword)
			if err != nil {
			} else {
				err := c.AppendResult(ctx, sp.Name, r)
				if err != nil {

				}
			}
		}(s)
	}
	group.Wait()
	return
}

// AppendResult
func (c *Client) AppendResult(ctx context.Context, name string, results []handler.Result) (err error) {
	c.resultLock.Lock()
	num := len(c.resultList)
	c.resultList = append(c.resultList, results...)
	defer c.resultLock.Unlock()
	for k, v := range results {
		co := color.New(color.FgGreen)
		_, err = co.Printf("[%d][%s]: %s\n", num+k+1, name, v.Title)
		if err != nil {
			return
		}
		co.DisableColor()
		_, err = co.Printf("%s\n", v.Body)
		if err != nil {
			return
		}
	}
	return
}

// ShowHelpInfo
func (c *Client) ShowHelpInfo() {
	fmt.Printf("Tip(%s): A command line tool for searching tips.\n", VersionName)
	fmt.Printf("Usage: tip [keywords]\n")
	fmt.Printf("Available handler(s):\n")
	for i, s := range configs.GetHandlerWrapperList() {
		fmt.Printf("[%d]%s: %s\n", i+1, s.Name, s.Description)
	}
}
