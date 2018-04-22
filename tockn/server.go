package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/VG-Tech-Dojo/vg-1day-2018-04-22/tockn/bot"
	"github.com/VG-Tech-Dojo/vg-1day-2018-04-22/tockn/controller"
	"github.com/VG-Tech-Dojo/vg-1day-2018-04-22/tockn/db"
	"github.com/VG-Tech-Dojo/vg-1day-2018-04-22/tockn/model"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// Server はAPIサーバーが実装された構造体です
type Server struct {
	db          *sql.DB
	Engine      *gin.Engine
	multicaster *bot.Multicaster
	poster      *bot.Poster
	bots        []*bot.Bot
}

// NewServer は新しいServerの構造体のポインタを返します
func NewServer() *Server {
	return &Server{
		Engine: gin.Default(),
	}
}

// Init はサーバーを初期化します
func (s *Server) Init(dbconf, env string) error {
	cs, err := db.NewConfigsFromFile(dbconf)
	if err != nil {
		return err
	}

	db, err := cs.Open(env)
	if err != nil {
		return err
	}
	s.db = db

	// routing
	s.Engine.LoadHTMLGlob("./templates/*")

	s.Engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	s.Engine.Static("/assets", "./assets")

	// tutorial. 自己紹介を追加する
	s.Engine.GET("/tockn", func(c *gin.Context) {
		c.HTML(http.StatusOK, "tockn.html", gin.H{})
	})

	// api
	api := s.Engine.Group("/api")
	api.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	msgStream := make(chan *model.Message)
	mctr := &controller.Message{DB: db, Stream: msgStream}
	api.GET("/messages", mctr.All)
	api.GET("/messages/:id", mctr.GetByID)
	api.POST("/messages", mctr.Create)
	api.PUT("/messages/:id", mctr.UpdateByID)
	api.DELETE("/messages/:id", mctr.DeleteByID)

	// bot
	mc := bot.NewMulticaster(msgStream)
	s.multicaster = mc

	poster := bot.NewPoster(10)
	s.poster = poster

	helloWorldBot := bot.NewHelloWorldBot(s.poster.In)
	s.bots = append(s.bots, helloWorldBot)
	omikujiBot := bot.NewOmikujiBot(s.poster.In)
	s.bots = append(s.bots, omikujiBot)
	keywordBot := bot.NewKeywordBot(s.poster.In)
	s.bots = append(s.bots, keywordBot)

	return nil
}

// Close はDBとの接続を閉じてサーバーを終了します
func (s *Server) Close() error {
	return s.db.Close()
}

// Run はサーバーを起動します
func (s *Server) Run(port string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// botを起動
	go s.multicaster.Run(ctx)
	go s.poster.Run(ctx, fmt.Sprintf("http://0.0.0.0:%s", port))

	for _, b := range s.bots {
		go b.Run(ctx)
		s.multicaster.BotIn <- b
	}

	s.Engine.Run(fmt.Sprintf(":%s", port))
}

func main() {
	var (
		dbconf = flag.String("dbconf", "dbconfig.yml", "database configuration file.")
		env    = flag.String("env", "development", "application envirionment (production, development etc.)")
		port   = flag.String("port", "8080", "listening port.")
	)
	flag.Parse()

	s := NewServer()
	if err := s.Init(*dbconf, *env); err != nil {
		log.Fatalf("fail to init server: %s", err)
	}
	defer s.Close()

	s.Run(*port)
}
