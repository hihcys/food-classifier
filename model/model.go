package model

import (
	"budeze/food-classifier/dao"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	"github.com/prometheus/common/log"
	"github.com/rqlite/gorqlite"
	"gopkg.in/yaml.v2"
)

var conn gorqlite.Connection

// WxSignature Tencent wexin MP signature
func WxSignature(timestamp, nonce string) string {
	token := "a4b46fe09f974c9f9e6b2cce5d04c57f"
	comArray := []string{token, timestamp, nonce}
	sort.Strings(comArray)
	hash := sha1.New()
	hash.Write([]byte(strings.Join(comArray, "")))
	hashBytes := hash.Sum(nil)
	hexStr := hex.EncodeToString(hashBytes)
	return hexStr
}

// WxTextResponseBody Tencent weixin MP text response body
func WxTextResponseBody(fromUserName, toUserName, content string) ([]byte, error) {
	if content != "" {
		foodMsg := GetFoodMsg(content)
		if foodMsg.FourSex != "" {
			content = fmt.Sprintf("【食物】%s\n【四性】%s\n【搭配适宜】%s\n【搭配忌讳】%s",
				content, foodMsg.FourSex, foodMsg.Suit, foodMsg.Taboo)
		} else {
			content = fmt.Sprintf("暂时未找到%s！", content)
		}
	}
	textResponseBody := &dao.TextResponseBody{}
	textResponseBody.FromUserName = fromUserName
	textResponseBody.ToUserName = toUserName
	textResponseBody.MsgType = "text"
	textResponseBody.Content = content
	textResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(textResponseBody, " ", "  ")
}

// InitConfig init service config
func InitConfig() *dao.Config {
	//read yml config
	config := dao.Config{}
	data, err := ioutil.ReadFile("../config/config.yml")
	if err != nil {
		log.Errorf("config file read error: %s\n", err.Error())
		return &config
	}
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Errorf("config file unmarshal error: %v", err)
		return &config
	}
	return &config
}

// InitRqlite init rqlite
func InitRqlite(rqCfg dao.RqliteConfig) {
	// connects to localhost on 4001 without auth
	var err error
	conn, err = gorqlite.Open(fmt.Sprintf("http://%s:%s@%s:%d/", rqCfg.Name, rqCfg.Pass, rqCfg.Host, rqCfg.Port))
	if err != nil {
		log.Errorf("open connect error: %s", err.Error())
		return
	}
	return
}

// GetFoodMsg 获取食品信息
func GetFoodMsg(foodName string) *dao.FoodMsg {
	foodMsg := dao.FoodMsg{}
	var id int64
	var name, alias, fourSex, suit, taboo, desc string
	rows, err := conn.QueryOne(fmt.Sprintf("SELECT id, name, alias, fourSex, suit, taboo, desc FROM food_msg WHERE name like '%s'", foodName))
	if err != nil {
		log.Errorf("query one error: %s\n", err.Error())
		return &foodMsg
	}
	if rows.NumRows() == 0 {
		rows, err = conn.QueryOne(fmt.Sprintf("SELECT id, name, alias, fourSex, suit, taboo, desc FROM food_msg WHERE alias like '%s'", foodName))
		if err != nil {
			log.Errorf("query one error: %s\n", err.Error())
			return &foodMsg
		}
		if rows.NumRows() == 0 {
			rows, err = conn.QueryOne(fmt.Sprintf("SELECT id, name, alias, fourSex, suit, taboo, desc FROM food_msg WHERE name like '%%%s%%'", foodName))
			if err != nil {
				log.Errorf("query one error: %s\n", err.Error())
				return &foodMsg
			}
			if rows.NumRows() == 0 {
				rows, err = conn.QueryOne(fmt.Sprintf("SELECT id, name, alias, fourSex, suit, taboo, desc FROM food_msg WHERE alias like '%%%s%%'", foodName))
				if err != nil {
					log.Errorf("query one error: %s\n", err.Error())
					return &foodMsg
				}
			}
		}
	}
	for rows.Next() {
		err = rows.Scan(&id, &name, &alias, &fourSex, &suit, &taboo, &desc)
		if err != nil {
			log.Errorf("query result scan error: %s\n", err.Error())
			return &foodMsg
		}
	}
	foodMsg = dao.FoodMsg{
		Id:      id,
		Name:    name,
		Alias:   alias,
		FourSex: fourSex,
		Suit:    suit,
		Taboo:   taboo,
		Desc:    desc,
	}
	return &foodMsg
}
