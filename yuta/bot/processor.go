package bot

import (
	"regexp"
	"strings"

	"fmt"

	"encoding/json"
	"github.com/VG-Tech-Dojo/vg-1day-2018-04-22/yuta/env"
	"github.com/VG-Tech-Dojo/vg-1day-2018-04-22/yuta/model"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	keywordAPIURLFormat = "https://jlp.yahooapis.jp/KeyphraseService/V1/extract?appid=%s&sentence=%s&output=json"
)

type (
	// Processor はmessageを受け取り、投稿用messageを作るインターフェースです
	Processor interface {
		Process(message *model.Message) (*model.Message, error)
	}

	// HelloWorldProcessor は"hello, world!"メッセージを作るprocessorの構造体です
	HelloWorldProcessor struct{}

	// OmikujiProcessor は"大吉", "吉", "中吉", "小吉", "末吉", "凶"のいずれかをランダムで作るprocessorの構造体です
	OmikujiProcessor struct{}

	// KeywordProcessor はメッセージ本文からキーワードを抽出するprocessorの構造体です
	KeywordProcessor struct{}

	GachaProcessor struct{}

	ChatProcessor struct{}
)

// Process は"hello, world!"というbodyがセットされたメッセージのポインタを返します
func (p *HelloWorldProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	return &model.Message{
		Body: msgIn.Body + ", world!",
	}, nil
}

// Process は"大吉", "吉", "中吉", "小吉", "末吉", "凶"のいずれかがbodyにセットされたメッセージへのポインタを返します
func (p *OmikujiProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	fortunes := []string{
		"大吉",
		"吉",
		"中吉",
		"小吉",
		"末吉",
		"凶",
	}
	result := fortunes[randIntn(len(fortunes))]
	return &model.Message{
		Body: result,
	}, nil
}

// Process はメッセージ本文からキーワードを抽出します
func (p *KeywordProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	r := regexp.MustCompile("\\Akeyword (.*)\\z")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	text := matchedStrings[1]

	url := fmt.Sprintf(keywordAPIURLFormat, env.KeywordAPIAppID, url.QueryEscape(text))

	type keywordAPIResponse map[string]interface{}
	var json keywordAPIResponse
	get(url, &json)

	keywords := []string{}
	for k, v := range json {
		if k == "Error" {
			return nil, fmt.Errorf("%#v", v)
		}
		keywords = append(keywords, k)
	}

	return &model.Message{
		Body: "キーワード：" + strings.Join(keywords, ", "),
	}, nil
}

func (p *GachaProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	rare := []string{
		"SSレア",
		"Sレア",
		"レア",
		"ノーマル",
	}
	result := rare[randIntn(len(rare))]
	return &model.Message{
		Body: result,
	}, nil
}

type TalkResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Results []struct {
		Perplexity float64 `json:"perplexity"`
		Reply      string  `json:"reply"`
	} `json:"results"`
}

func (p *ChatProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	r := regexp.MustCompile("\\Atalk (.*)\\z")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	text := matchedStrings[1]

	values := url.Values{}
	values.Add("apikey", env.ChatAPIAppID)
	values.Add("query", text)

	resp, err := http.PostForm("https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk", values)
	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var talkResponse TalkResponse
	err = json.Unmarshal(body, &talkResponse)
	if err != nil {
		return nil, err
	}

	return &model.Message{
		Body: talkResponse.Results[0].Reply,
	}, nil
}
