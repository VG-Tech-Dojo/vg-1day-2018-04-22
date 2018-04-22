package bot

import (
	"regexp"
	"strings"

	"fmt"

	"encoding/json"
	"github.com/VG-Tech-Dojo/vg-1day-2018-04-22/ninja/env"
	"github.com/VG-Tech-Dojo/vg-1day-2018-04-22/ninja/model"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	keywordAPIURLFormat = "https://jlp.yahooapis.jp/KeyphraseService/V1/extract?appid=%s&sentence=%s&output=json"
	talkAPIURLFormat    = "https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk"
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

	// GachaProcessor は"SSレア", "Sレア", "レア", "ノーマル"のいずれかをランダムで作るprocessorの構造体です
	GachaProcessor struct{}

	// TalkProcessor はメッセージに対する返信を返すprocessorの構造体です
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

// Process は"SSレア", "Sレア", "レア", "ノーマル"のいずれかがbodyにセットされたメッセージへのポインタを返します
func (p *GachaProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	rarelity := []string{
		"SSレア",
		"Sレア",
		"レア",
		"ノーマル",
	}
	result := rarelity[randIntn(len(rarelity))]
	return &model.Message{
		Body: result,
	}, nil
}

// Process はメッセージ本文に対する返信を返します
func (p *TalkProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	r := regexp.MustCompile("\\Atalk (.*)\\z")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	text := matchedStrings[1]

	values := url.Values{}
	values.Add("apikey", env.TalkAPIKey)
	values.Add("query", text)

	res, err := http.PostForm(
		talkAPIURLFormat,
		values,
	)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	var decoded interface{}
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, err
	}

	return &model.Message{
		Body: decoded.(map[string]interface{})["results"].([]interface{})[0].(map[string]interface{})["reply"].(string),
	}, nil
}
