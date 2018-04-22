package model

import (
	"database/sql"
)

// Message はメッセージの構造体です
type Message struct {
	ID       int64  `json:"id"`
	Body     string `json:"body"`
	Username string `json:"username"`
}

// MessagesAll は全てのメッセージを返します
func MessagesAll(db *sql.DB) ([]*Message, error) {

	rows, err := db.Query(`select id, username, body from message`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ms []*Message
	for rows.Next() {
		m := &Message{}
		if err := rows.Scan(&m.ID, &m.Username, &m.Body); err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ms, nil
}

// MessageByID は指定されたIDのメッセージを1つ返します
func MessageByID(db *sql.DB, id string) (*Message, error) {
	m := &Message{}

	if err := db.QueryRow(`select id, username, body from message where id = ?`, id).Scan(&m.ID, &m.Username, &m.Body); err != nil {
		return nil, err
	}

	return m, nil
}

// Insert はmessageテーブルに新規データを1件追加します
func (m *Message) Insert(db *sql.DB) (*Message, error) {
	res, err := db.Exec(`insert into message (username, body) values (?, ?)`, m.Username, m.Body)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Message{
		ID:       id,
		Body:     m.Body,
		Username: m.Username,
	}, nil
}

// Update はmessageテーブルのデータを更新します
func (m *Message) Update(db *sql.DB) error {
	res, err := db.Exec(`update message set body = ? where id = ?`, m.Body, m.ID)
	if err != nil {
		return err
	}
	_, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

// 1-4. メッセージを削除しよう
// ...
