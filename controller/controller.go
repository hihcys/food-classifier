package controller

import (
	"budeze/food-classifier/dao"
	"budeze/food-classifier/model"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
)

// LogFormat is the structure any formatter will be handed when time to log comes
type LogFormat struct {
	TimeStamp  time.Time
	StatusCode int
	Latency    time.Duration
	ClientIP   string
	Keys       map[string]interface{}
}

// InitController init controller
func InitController(router *gin.Engine) {
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logFormat := LogFormat{
			TimeStamp:  param.TimeStamp,
			StatusCode: param.StatusCode,
			Latency:    param.Latency,
			ClientIP:   param.ClientIP,
			Keys:       param.Keys,
		}
		rep, _ := json.Marshal(logFormat)
		return fmt.Sprintf("%v\n", string(rep))
	}))

	v1 := router.Group("/v1")
	{
		v1.GET("/food/foursex", foursexHandlerGet)
		v1.POST("/food/foursex", foursexHandlerPost)
	}
}

// foursexHandlerGet classifier of food cold and heat
func foursexHandlerGet(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")
	sign := model.WxSignature(timestamp, nonce)
	if sign == signature {
		c.String(http.StatusOK, echostr)
	} else {
		c.String(http.StatusOK, "")
	}
}

// foursexHandlerPost classifier of food cold and heat
func foursexHandlerPost(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	sign := model.WxSignature(timestamp, nonce)
	if sign == signature {
		fmt.Println(c.Query("signature"))
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Errorf("get request body err: %s", err.Error())
			return
		}
		requestBody := &dao.TextRequestBody{}
		xml.Unmarshal(body, requestBody)
		fmt.Println(requestBody)
		responseBody, err := model.WxTextResponseBody(requestBody.ToUserName,
			requestBody.FromUserName, requestBody.Content)
		if err != nil {
			log.Errorf("marshal response body err: %s", err.Error())
			return
		}
		c.Writer.Header().Set("Content-Type", "text/xml")
		fmt.Fprintf(c.Writer, string(responseBody))
	} else {
		c.String(200, "")
	}
}
