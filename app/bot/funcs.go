package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// funcs bot, returns numArticles last articles in MD format from https://funcs.radio-t.com/api/v1/funcs/lastmd/5
type Funcs struct {
	client       HTTPClient
	funcsAPI     string
	funcName     string
}

type funcsArticle struct {
	Title string    `json:"title"`
	Link  string    `json:"link"`
	Ts    time.Time `json:"ats"`
}

// NewFuncs makes new funcs bot
func NewFuncs(client HTTPClient, api string, name string) *Funcs {
	log.Printf("[INFO] funcs bot with api %s", api)
	return &Funcs{client: client, funcsAPI: api, funcName: name}
}

// Help returns help message
func (n Funcs) Help() string {
	return genHelpMsg(n.ReactOn(), "Описание функции Krista BI")
}

// OnMessage returns N last funcs articles
func (n Funcs) OnMessage(msg Message) (response Response) {
	if !contains(n.ReactOn(), msg.Text) {
		return Response{}
	}

	reqURL := fmt.Sprintf("%s/engine/function//%s", n.funcsAPI, n.funcName)
	log.Printf("[DEBUG] request %s", reqURL)

	req, err := makeHTTPRequest(reqURL)
	if err != nil {
		log.Printf("[WARN] failed to make request %s, error=%v", reqURL, err)
		return Response{}
	}

	resp, err := n.client.Do(req)
	if err != nil {
		log.Printf("[WARN] failed to send request %s, error=%v", reqURL, err)
		return Response{}
	}
	defer resp.Body.Close()

	articles := []funcsArticle{}
	if err = json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		log.Printf("[WARN] failed to parse response, error %v", err)
		return Response{}
	}

	var lines []string
	for _, a := range articles {
		lines = append(lines, fmt.Sprintf("- [%s](%s) %s", a.Title, a.Link, a.Ts.Format("2006-01-02")))
	}
	return Response{
		Text: strings.Join(lines, "\n") + "\n- [все функции](https://help.krista.ru/kb/4217)",
		Send: true,
	}
}

// ReactOn keys
func (n Funcs) ReactOn() []string {
	return []string{"func!", "функция!"}
}
