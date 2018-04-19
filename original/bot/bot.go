package bot

import (
	"context"
	"log"

	"github.com/VG-Tech-Dojo/vg-1day-2018-04-22/original/model"
)

type (
	// Bot はinで受け取ったmessageがcheckerの条件を満たした場合、processorが投稿用messageを作り、outに渡します
	//
	//   fields
	//     name      string
	//     in        chan *model.Message
	//     out       chan *model.Message
	//     checker   Checker
	//     processor Processor
	Bot struct {
		name      string
		in        chan *model.Message
		out       chan *model.Message
		checker   Checker
		processor Processor
	}
)

// Run はBotを起動します
func (b *Bot) Run(ctx context.Context) {
	// メッセージ監視
	for {
		select {
		case <-ctx.Done():
			close(b.in)
			return
		case m := <-b.in:
			if b.checker.Check(m) {
				nm, err := b.processor.Process(m)
				if err != nil {
					log.Printf("%s: %#v\n", b.name, err)
					b.out <- &model.Message{
						Body: "気が乗らないパカ",
					}
					// selectから抜ける
					break
				}
				b.out <- nm
			}
		}
	}
}

// NewHelloWorldBot は"hello"を受け取ると"hello, world!"を返す新しいBotの構造体のポインタを返します
func NewHelloWorldBot(out chan *model.Message) *Bot {
	in := make(chan *model.Message)

	checker := NewRegexpChecker("\\Ahello\\z")

	processor := &HelloWorldProcessor{}

	return &Bot{
		name:      "helloworldbot",
		in:        in,
		out:       out,
		checker:   checker,
		processor: processor,
	}
}

// NewOmikujiBot は"大吉", "吉", "中吉", "小吉", "末吉", "凶"のいずれかをランダムで返す新しいBotの構造体のポインタを返します
func NewOmikujiBot(out chan *model.Message) *Bot {
	in := make(chan *model.Message)

	checker := NewRegexpChecker("\\Aomikuji\\z")

	processor := &OmikujiProcessor{}

	return &Bot{
		name:      "omikujibot",
		in:        in,
		out:       out,
		checker:   checker,
		processor: processor,
	}
}

// NewKeywordBot はメッセージ本文からキーワードを抽出して返す新しいBotの構造体のポインタを返します
func NewKeywordBot(out chan *model.Message) *Bot {
	in := make(chan *model.Message)

	checker := NewRegexpChecker("\\Akeyword .*")

	processor := &KeywordProcessor{}

	return &Bot{
		name:      "keywordbot",
		in:        in,
		out:       out,
		checker:   checker,
		processor: processor,
	}
}
