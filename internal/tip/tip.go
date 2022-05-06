package tip

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/tenhan/tip/configs"
	"github.com/tenhan/tip/internal/handler"
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

type DefaultFieldsHook struct{}

func (df *DefaultFieldsHook) Fire(entry *log.Entry) error {
	entry.Data["handler"] = entry.Context.Value("name")
	return nil
}

func (df *DefaultFieldsHook) Levels() []log.Level {
	return log.AllLevels
}

// Run
func (c *Client) Run(ctx context.Context) (err error) {
	var verbose = flag.Bool("v", false, "show debug log")
	var veryVerbose = flag.Bool("vv", false, "show trace log")
	var showHistory = flag.Bool("l", false, "show history")
	flag.Parse()
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.AddHook(&DefaultFieldsHook{})
	if *veryVerbose {
		log.SetLevel(log.TraceLevel)
	} else if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
	if *showHistory {
		return c.ShowHistory(ctx)
	}
	keyword := strings.Join(flag.Args(), " ")
	if c.handlerWrapperList == nil || len(c.handlerWrapperList) == 0 {
		err = fmt.Errorf("no spider found")
		return
	}
	log.WithContext(ctx).Debugf("client start, keyword: %s", keyword)
	group := sync.WaitGroup{}
	group.Add(len(c.handlerWrapperList))
	for _, s := range c.handlerWrapperList {
		go func(sp handler.Wrapper) {
			defer group.Done()
			startAt := time.Now()
			valueCtx := context.WithValue(ctx, "name", sp.Name)
			cancelCtx, cancel := context.WithCancel(valueCtx)
			defer cancel()
			log.WithContext(cancelCtx).Debug("Handle start")
			r, err1 := sp.Handler.Handle(cancelCtx, keyword)
			if err1 != nil {
				log.WithContext(cancelCtx).Errorf("Handle failed, err: %v", err1)
			} else {
				endAt := time.Now()
				result := &handler.HandlerResult{
					Name:     sp.Name,
					Keyword:  keyword,
					Results:  r,
					StartAt:  startAt,
					EndAt:    endAt,
					Duration: endAt.Sub(startAt).Milliseconds(),
				}
				if err := c.AddHistory(ctx, result); err != nil {
					log.WithContext(ctx).Errorf("save history failed: %v, result: %v", err, result)
				}
				log.WithContext(cancelCtx).Debugf("Handle finished, result count: %d", len(r))
				if err := c.AppendResult(cancelCtx, result); err != nil {
					log.WithContext(cancelCtx).Errorf("AppendResult failed: %v", err)
				}
			}
		}(s)
	}
	group.Wait()
	return nil
}

func (c *Client) GetResultList() []handler.Result {
	return c.resultList
}

func (c *Client) AddHistory(ctx context.Context, result *handler.HandlerResult) error {
	// save to history file
	filename, err := c.GetHistoryFilename()
	if err != nil {
		return err
	}
	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	resultJson, err := json.Marshal(&result)
	if err != nil {
		return err
	}
	resultJson = append(resultJson, '\n')
	_, err = fd.Write(resultJson)
	if err != nil {
		return err
	}
	return fd.Close()
}

// AppendResult
func (c *Client) AppendResult(ctx context.Context, result *handler.HandlerResult) (err error) {
	c.resultLock.Lock()
	num := len(c.resultList)
	c.resultList = append(c.resultList, result.Results...)
	defer c.resultLock.Unlock()
	prefix := ""
	if num > 0 {
		prefix = "\n"
	}
	for k, v := range result.Results {
		co := color.New(color.FgGreen)
		_, err = co.Printf("%s[%d][%s][%s]: %s\n", prefix, num+1+k, result.StartAt.Format(time.RFC3339), result.Name, v.Title)
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

func (c *Client) GetHistoryFilename() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/.tip_history", dir), nil
}

func (c *Client) ShowHistory(ctx context.Context) error {
	filename, err := c.GetHistoryFilename()
	if err != nil {
		return err
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	r := bufio.NewReader(file)
	obj := &handler.HandlerResult{}
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		err = json.Unmarshal(line, obj)
		if err != nil {
			return err
		}
		err = c.AppendResult(ctx, obj)
		if err != nil {
			return err
		}
	}
	return nil
}
