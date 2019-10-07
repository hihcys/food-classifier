package dao

import (
	"encoding/xml"
	"time"
)

type TextRequestBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
}

type TextResponseBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
}

type FoodMsg struct {
	Id      int64
	Name    string
	Alias   string
	FourSex string
	Suit    string
	Taboo   string
	Desc    string
}

type Config struct {
	Rqlite RqliteConfig `yaml:"rqlite"`
}

type RqliteConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Name string `yaml:"name"`
	Pass string `yaml:"pass"`
}
