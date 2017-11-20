package models

import (
	"github.com/rs/xid"
	"time"
)

type Article struct {
	Id         string `bson:"_id"`
	Category   string
	Title      string
	Content    string
	CreateDate string
}

func NewArticle(category string, title string, content string) *Article {
	runeTime := []rune(time.Now().String())

	return &Article{
		genId(),
		category,
		title,
		content,
		string(runeTime[0:19]),
	}
}

func genId() string {
	return xid.New().String()
}
