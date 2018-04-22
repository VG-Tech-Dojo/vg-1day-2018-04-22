package bot

import (
	"context"

	"github.com/VG-Tech-Dojo/vg-1day-2018-04-22/ihaya/model"
)

// Multicaster は1つのチャンネルで複数botを動かすためのヘルパーです
//
// msgInで受け取ったmessageをbotsに登録された全botに渡します
//
// botsへの登録はBotInで行います
//
//   fields
// 	   BotIn chan *Bot
// 	   bots  map[*Bot]bool
// 	   msgIn chan *model.Message
type Multicaster struct {
	BotIn chan *Bot
	bots  []*Bot
	msgIn chan *model.Message
}

// Run はMulticasterを起動します
func (mc *Multicaster) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(mc.msgIn)
			return
		case bot := <-mc.BotIn:
			mc.bots = append(mc.bots, bot)
		case msg := <-mc.msgIn:
			for _, bot := range mc.bots {
				bot.in <- msg
			}
		}
	}
}

// NewMulticaster は新しいMulticaster構造体のポインタを返します
func NewMulticaster(msgIn chan *model.Message) *Multicaster {
	memberIn := make(chan *Bot)
	return &Multicaster{
		BotIn: memberIn,
		bots:  []*Bot{},
		msgIn: msgIn,
	}
}
