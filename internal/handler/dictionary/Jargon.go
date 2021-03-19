package dictionary

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tenhan/tip/internal/handler"
	"golang.org/x/net/html/charset"
	"net/http"
	"strings"
)

type Jargon struct {
}

func (s Jargon) Handle(ctx context.Context, keyword string) (results []handler.Result, err error) {
	link := "http://www.catb.org/~esr/jargon/html/go01.html"
	res, err := http.Get(link)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		return
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}
	type Node struct {
		Text string
		Href string
	}
	keywordMap := make(map[string]Node)
	doc.Find(".titlepage").Children().Each(func(i int, selection *goquery.Selection) {
		selection.Find("dt").Each(func(i int, selection *goquery.Selection) {
			a := selection.Find("a")
			text := a.Text()
			if text != "" {
				href, exist := a.Attr("href")
				if exist {
					keywordMap[strings.ToLower(text)] = Node{
						Text: text,
						Href: href,
					}
				}
			}
		})
	})

	// case insensitive
	value, ok := keywordMap[strings.ToLower(keyword)]
	if !ok {
		return
	}
	nextLink := "http://www.catb.org/~esr/jargon/html/" + value.Href
	res, err = http.Get(nextLink)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		return
	}
	// encoding="ISO-8859-1"
	d, err := charset.NewReader(res.Body, "ISO-8859-1")
	if err != nil {
		return
	}
	doc, err = goquery.NewDocumentFromReader(d)
	if err != nil {
		return
	}
	var lines []string
	doc.Find("body").Children().Each(func(i int, selection *goquery.Selection) {
		lines = append(lines, selection.Text())
	})
	// ignore fist and last one
	if len(lines) > 0 {
		lines = lines[1:]
	}
	if len(lines) > 0 {
		lines = lines[:len(lines)-1]
	}

	results = append(results, handler.Result{
		Title: value.Text,
		Body:  strings.Join(lines, "\n"),
	})
	return
}
