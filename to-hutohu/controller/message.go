package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/VG-Tech-Dojo/vg-1day-2018-04-22/to-hutohu/httputil"
	"github.com/VG-Tech-Dojo/vg-1day-2018-04-22/to-hutohu/model"
	"github.com/gin-gonic/gin"
)

// Message is controller for requests to messages
type Message struct {
	DB     *sql.DB
	Stream chan *model.Message
}

// All は全てのメッセージを取得してJSONで返します
func (m *Message) All(c *gin.Context) {
	msgs, err := model.MessagesAll(m.DB)
	if err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	if len(msgs) == 0 {
		c.JSON(http.StatusOK, make([]*model.Message, 0))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": msgs,
		"error":  nil,
	})
}

// GetByID はパラメーターで受け取ったidのメッセージを取得してJSONで返します
func (m *Message) GetByID(c *gin.Context) {
	msg, err := model.MessageByID(m.DB, c.Param("id"))

	switch {
	case err == sql.ErrNoRows:
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusNotFound, resp)
		return
	case err != nil:
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": msg,
		"error":  nil,
	})
}

// Create は新しいメッセージ保存し、作成したメッセージをJSONで返します
func (m *Message) Create(c *gin.Context) {
	var msg model.Message

	if c.Request.ContentLength == 0 {
		resp := httputil.NewErrorResponse(errors.New("body is missing"))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := c.BindJSON(&msg); err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// 1-2. ユーザー名を追加しよう
	// できる人は、ユーザー名が空だったら`anonymous`等適当なユーザー名で投稿するようにしてみよう

	inserted, err := msg.Insert(m.DB)
	if err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// bot対応
	m.Stream <- inserted

	c.JSON(http.StatusCreated, gin.H{
		"result": inserted,
		"error":  nil,
	})
}

// UpdateByID は...
func (m *Message) UpdateByID(c *gin.Context) {
	var msg model.Message

	if c.Request.ContentLength == 0 {
		resp := httputil.NewErrorResponse(errors.New("body is missing"))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := c.BindJSON(&msg); err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	err := msg.Update(m.DB)
	switch {
	case err == sql.ErrNoRows:
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusNotFound, resp)
		return
	case err != nil:
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteByID は...
func (m *Message) DeleteByID(c *gin.Context) {
	messageID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	msg := &model.Message{
		ID: messageID,
	}
	if err := msg.Delete(m.DB); err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (m *Message) SaveImage(c *gin.Context) {
	uploadedFile, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	file, err := os.Create(`./assets/imgs/` + uploadedFile.Filename)
	if err != nil {
		c.Error(err)
		return
	}
	defer file.Close()

	tmp, err := uploadedFile.Open()
	if err != nil {
		c.Error(err)
		return
	}

	data, err := ioutil.ReadAll(tmp)
	if err != nil {
		c.Error(err)
		return
	}
	file.Write(data)

	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully.", uploadedFile.Filename))
}
