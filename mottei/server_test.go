package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	dbconf = "dbconfig.yml"
	env    = "test"
	port   = "50000"
)

var tsURL = "http://localhost:" + port

func TestMain(m *testing.M) {
	os.Exit(realMain(m))
}

func realMain(m *testing.M) int {
	s := NewServer()
	if err := s.Init(dbconf, env); err != nil {
		panic(fmt.Sprintf("failed to init server: %v", err))
	}
	go s.Run(port)
	defer s.Close()

	return m.Run()
}

func TestTopページが200を返す(t *testing.T) {
	resp, err := http.Get(tsURL + "/")
	if err != nil {
		t.Fatalf("failed to get response: %s", err)
	}
	defer resp.Body.Close()

	if expected := 200; resp.StatusCode != expected {
		t.Fatalf("status code expected %d but not, actual %d", expected, resp.StatusCode)
	}
}

func TestAPIがpingに応答する(t *testing.T) {
	resp, err := http.Get(tsURL + "/api/ping")
	if err != nil {
		t.Fatalf("failed to get response: %s", err)
	}
	defer resp.Body.Close()

	if expected := 200; resp.StatusCode != expected {
		t.Fatalf("status code expected %d but not, actual %d", expected, resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read http response, %s", err)
	}

	if expected := "pong"; string(b) != expected {
		t.Fatalf("response body expected %s but not, actual %s", expected, string(b))
	}
}

func TestAPIがメッセージを全て返す(t *testing.T) {
	resp, err := http.Get(tsURL + "/api/messages")
	if err != nil {
		t.Fatalf("failed to get response: %s", err)
	}
	defer resp.Body.Close()

	if expected := 200; resp.StatusCode != expected {
		t.Fatalf("status code expected %d but not, actual %d", expected, resp.StatusCode)
	}

	if expected := "application/json; charset=utf-8"; resp.Header.Get("Content-Type") != expected {
		t.Fatalf("response header expected %s but not, actual: %s", expected, resp.Header.Get("Content-Type"))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read http response, %s", err)
	}

	expected := `{"error":null,"result":[{"id":1,"body":"hoge"},{"id":2,"body":"fuga"},{"id":3,"body":"piyo"}]}`
	// http responseの末尾に改行が含まれるので除去して比較します
	actual := strings.TrimRight(string(b), "\n")
	if actual != expected {
		t.Fatalf("response body expected %s, but %s", expected, string(b))
	}
}

func TestAPIが指定したIDのメッセージを返す(t *testing.T) {
	resp, err := http.Get(tsURL + "/api/messages/1")
	if err != nil {
		t.Fatalf("failed to get response: %s", err)
	}
	defer resp.Body.Close()

	if expected := 200; resp.StatusCode != expected {
		t.Fatalf("status code expected %d but not, actual %d", expected, resp.StatusCode)
	}

	if expected := "application/json; charset=utf-8"; resp.Header.Get("Content-Type") != expected {
		t.Fatalf("response header expected %s but not, actual: %s", expected, resp.Header.Get("Content-Type"))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read http response, %s", err)
	}

	expected := `{"error":null,"result":{"id":1,"body":"hoge"}}`
	// http responseの末尾に改行が含まれるので除去して比較します
	actual := strings.TrimRight(string(b), "\n")
	if actual != expected {
		t.Fatalf("response body expected %s, but %s", expected, string(b))
	}
}

func TestAPIが新しいメッセージを作成する(t *testing.T) {
	tm := "testmessage"
	resp, err := http.Post(tsURL+"/api/messages", "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(`{"body": "%s"}`, tm))))
	if err != nil {
		t.Fatalf("failed to post request: %s", err)
	}
	defer resp.Body.Close()

	if expected := 201; resp.StatusCode != expected {
		t.Fatalf("status code expected %d but not, actual %d", expected, resp.StatusCode)
	}

	if expected := "application/json; charset=utf-8"; resp.Header.Get("Content-Type") != expected {
		t.Fatalf("response header expected %s but not, actual: %s", expected, resp.Header.Get("Content-Type"))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read http response, %s", err)
	}

	expected := fmt.Sprintf(`{"error":null,"result":{"id":4,"body":"%s"}}`, tm)
	// http responseの末尾に改行が含まれるので除去して比較します
	actual := strings.TrimRight(string(b), "\n")
	if actual != expected {
		t.Fatalf("response body expected %s, but %s", expected, string(b))
	}
}

func TestHelloWorldBotが反応する(t *testing.T) {
	// botが反応するキーワードを投稿する
	r, err := http.Post(tsURL+"/api/messages", "application/json", bytes.NewBuffer([]byte(`{"body": "hello"}`)))
	if err != nil {
		t.Fatalf("failed to post request: %s", err)
	}
	defer r.Body.Close()

	// 最新メッセージを取得する
	time.Sleep(1 * time.Second)
	resp, err := http.Get(tsURL + "/api/messages/6")
	if err != nil {
		t.Fatalf("failed to get response: %s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read http response, %s", err)
	}

	expected := `{"error":null,"result":{"id":6,"body":"hello, world!"}}`
	// http responseの末尾に改行が含まれるので除去して比較します
	actual := strings.TrimRight(string(b), "\n")
	if actual != expected {
		t.Fatalf("response body expected %s, but %s", expected, string(b))
	}
}

func TestAPIが指定したIDのメッセージを更新する(t *testing.T) {}

func TestAPIが指定したIDのメッセージを削除する(t *testing.T) {}
