package dictionary

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/tenhan/tip/configs"
	"github.com/tenhan/tip/internal/handler"
	"github.com/tenhan/tip/pkg/utils/str"
)

type BaiduTranslate struct{}

func (s BaiduTranslate) Handle(ctx context.Context, keyword string) (results []handler.Result, err error) {
	// Is keyword an English word?
	if !str.IsAlpha(keyword) {
		log.WithContext(ctx).Debugf("ignore keyword: %s", keyword)
		return
	}
	link := fmt.Sprintf("http://www.baidu.com/s?wd=%s", keyword)

	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}
	var cookies []*http.Cookie
	cookieURL, err := url.Parse(link)
	if err != nil {
		return
	}
	cookieFile := fmt.Sprintf("%s/%s.cookies.json", configs.CookiesPath, cookieURL.Host)
	// ignore cookies error
	cookies, _ = readLocalCookies(cookieFile)
	jar.SetCookies(cookieURL, cookies)
	log.WithContext(ctx).Debugf("set request cookies: %v", cookies)
	client := &http.Client{
		Jar: jar,
	}
	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return
	}
	request.Header.Add("User-Agent", configs.UserAgent)

	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	// save cookies
	responseCookies := jar.Cookies(cookieURL)
	log.WithContext(ctx).Debugf("response cookies: %v, write to file: %s", responseCookies, cookieFile)
	err = writeLocalCookies(cookieFile, responseCookies)
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return
	}
	contentLeft := doc.Find("#content_left")
	if contentLeft.Size() == 0 {
		log.WithContext(ctx).Warnf("%s not found", keyword)
		return
	}
	endWith := " - 百度翻译"
	contentLeft.Children().Each(func(i int, selection *goquery.Selection) {
		title := str.Trim(selection.Find("h3").Text())
		if str.EndWith(title, endWith) {
			title = strings.Replace(title, endWith, "", 1)
			var lines []string
			contentNode := selection.Find(".op_dict_content")
			// if .op_dict_content not found, find .op_sp_fanyi
			if contentNode.Children().Size() == 0 {
				contentNode = selection.Find(".op_sp_fanyi")
			}
			contentNode.Children().Each(func(i int, s1 *goquery.Selection) {
				line := s1.Text()
				line = strings.Replace(line, "\n", " ", -1)
				line = strings.Replace(line, "\r", " ", -1)
				line = strings.Replace(line, "\t", " ", -1)
				line = str.RemoveDuplicatedWhiteSpace(line)
				line = str.Trim(line)
				if line != "" {
					lines = append(lines, line)
				}
			})
			if len(lines) > 1 {
				lines = lines[0 : len(lines)-1]
			}
			results = append(results, handler.Result{
				Title: title,
				Body:  strings.Join(lines, "\n"),
			})
		}
	})
	return
}

func readLocalCookies(filename string) (cookies []*http.Cookie, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		//no such file or directory
		return
	}
	err = json.Unmarshal(data, &cookies)
	if err != nil {
		return
	}
	return
}
func writeLocalCookies(filename string, cookies []*http.Cookie) error {
	data, err := json.Marshal(cookies)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, data, 0664)
	if err != nil {
		return err
	}
	return nil
}
