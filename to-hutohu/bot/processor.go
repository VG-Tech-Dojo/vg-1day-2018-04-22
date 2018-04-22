package bot

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strings"

	"fmt"

	"net/http"
	"net/url"

	"github.com/VG-Tech-Dojo/vg-1day-2018-04-22/to-hutohu/env"
	"github.com/VG-Tech-Dojo/vg-1day-2018-04-22/to-hutohu/model"
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

	// GachaProcessor はガチャの結果を返すprocessorの構造体です
	GachaProcessor struct{}

	// TalkProcessor はなんか話してくれるprocerrorの構造体です
	TalkProcessor struct{}
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

// Process はガチャの結果を返します
func (p *GachaProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	result := []string{"SSレア", "Sレア", "レア", "ノーマル"}[rand.Intn(4)]
	return &model.Message{
		Username: "Gacha bot",
		Body:     fmt.Sprintf("ガチャの結果は%sです！！", result),
	}, nil
}

// Process は会話の返答を返します
func (p *TalkProcessor) Process(msgIn *model.Message) (*model.Message, error) {

	r := regexp.MustCompile("\\Atalk (.*)\\z")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	text := matchedStrings[1]

	form := url.Values{}
	form.Set("apikey", env.RecruitTalkAPIToken)
	form.Set("query", text)
	res, err := http.PostForm("https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk", form)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	resultData := &struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Results []struct {
			Reply string `json:"reply"`
		} `json:"results"`
	}{}
	err = json.Unmarshal(respBody, resultData)
	if err != nil {
		return nil, err
	}

	if resultData.Status != 0 {
		return nil, errors.New(resultData.Message)
	}

	reply := ""
	if len(resultData.Results) == 0 {
		reply = "返答はない"
	} else {
		reply = resultData.Results[0].Reply
	}

	return &model.Message{
		Username: "Talk bot",
		Body:     reply,
	}, nil
}
