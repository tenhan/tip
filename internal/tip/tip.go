package tip

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/tenhan/tip/configs"
	"github.com/tenhan/tip/internal/handler"
	"os"
	"strings"
	"sync"
	"time"
)



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

type DefaultFieldsHook struct {}

func (df *DefaultFieldsHook) Fire(entry *log.Entry) error {
	entry.Data["handler"] = entry.Context.Value("name")
	return nil
}

func (df *DefaultFieldsHook) Levels() []log.Level {
	return log.AllLevels
}

// Run
func (c *Client) Run(ctx context.Context) (err error) {
	startAt := time.Now()
	var verbose = flag.Bool("v", false, "show debug log")
	flag.Parse()
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.AddHook(&DefaultFieldsHook{})
	if *verbose{
		log.SetLevel(log.DebugLevel)
	}else{
		log.SetLevel(log.WarnLevel)
	}
	keyword := strings.Join(flag.Args()," ")
	if c.handlerWrapperList == nil || len(c.handlerWrapperList) == 0 {
		err = fmt.Errorf("no spider found")
		return
	}
	log.WithContext(ctx).Debugf("client start, keyword: %s",keyword)
	group := sync.WaitGroup{}
	group.Add(len(c.handlerWrapperList))
	for _, s := range c.handlerWrapperList {
		go func(sp handler.Wrapper) {
			defer group.Done()
			valueCtx := context.WithValue(ctx,"name",sp.Name)
			cancelCtx,cancel := context.WithCancel(valueCtx)
			defer cancel()
			log.WithContext(cancelCtx).Debug("Handle start")
			r, err1 := sp.Handler.Handle(cancelCtx, keyword)
			if err1 != nil {
				log.WithContext(cancelCtx).Errorf("Handle failed, err: %v", err1)
			} else {
				log.WithContext(cancelCtx).Debugf("Handle finished, result count: %d",len(r))
				err2 := c.AppendResult(cancelCtx, sp.Name, r)
				if err2 != nil {
					log.WithContext(cancelCtx).Errorf("AppendResult failed, err: %v",err2)
				}
			}
		}(s)
	}
	group.Wait()
	endAt := time.Now()
	// save to history file
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal( err )
	}
	filename := fmt.Sprintf("%s/.tip_history",dirname)
	fd,err :=os.OpenFile(filename,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	if err != nil {
		log.WithContext(ctx).Errorf("can not open history file %s, err: %s",filename,err)
		return
	}

	history := handler.HistoryObject{
		Keyword:  keyword,
		Results:  c.GetResult(),
		StartAt:  startAt,
		EndAt:    endAt,
		Duration: endAt.Sub(startAt).Milliseconds(),
	}
	historyJson,err := json.Marshal(&history)
	if err != nil {
		log.WithContext(ctx).Errorf("json.Marshal history fail: %s",err)
		return
	}
	historyJson = append(historyJson,'\n')
	_,err = fd.Write(historyJson)
	if err != nil {
		log.WithContext(ctx).Errorf("write buff to history file fail: %s",filename)
		return
	}
	return fd.Close()
}

func (c *Client) GetResult() []handler.Result  {
	return c.resultList
}
// AppendResult
func (c *Client) AppendResult(ctx context.Context, name string, results []handler.Result) (err error) {
	c.resultLock.Lock()
	num := len(c.resultList)
	c.resultList = append(c.resultList, results...)
	defer c.resultLock.Unlock()
	prefix := ""
	if num > 0{
		prefix = "\n"
	}
	for _, v := range results {
		co := color.New(color.FgGreen)
		_, err = co.Printf("%s[%s]: %s\n", prefix, name, v.Title)
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
	fmt.Printf("Tip(%s): A command line tool for searching tips.\n", configs.VersionName)
	fmt.Printf("Usage: tip [keywords]\n")
	fmt.Printf("Available handler(s):\n")
	for i, s := range GetHandlerWrapperList() {
		fmt.Printf("[%d]%s: %s\n", i+1, s.Name, s.Description)
	}
}
