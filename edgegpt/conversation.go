package edgegpt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

var HEADERS_INIT_CONVER = map[string]string{
	"authority":                   "edgeservices.bing.com",
	"accept":                      "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	"accept-language":             "en-US,en;q=0.9",
	"cache-control":               "max-age=0",
	"sec-ch-ua":                   `"Chromium";v="110", "Not A(Brand";v="24", "Microsoft Edge";v="110"`,
	"sec-ch-ua-arch":              `"x86"`,
	"sec-ch-ua-bitness":           `"64"`,
	"sec-ch-ua-full-version":      `"110.0.1587.69"`,
	"sec-ch-ua-full-version-list": `"Chromium";v="110.0.5481.192", "Not A(Brand";v="24.0.0.0", "Microsoft Edge";v="110.0.1587.69"`,
	"sec-ch-ua-mobile":            "?0",
	"sec-ch-ua-model":             `""`,
	"sec-ch-ua-platform":          `"Windows"`,
	"sec-ch-ua-platform-version":  `"15.0.0"`,
	"sec-fetch-dest":              "document",
	"sec-fetch-mode":              "navigate",
	"sec-fetch-site":              "none",
	"sec-fetch-user":              "?1",
	"upgrade-insecure-requests":   "1",
	"user-agent":                  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.69",
	"x-edge-shopping-flag":        "1",
	"x-forwarded-for":             "1.1.1.1",
}

// Conversation API
type Conversation struct {
	Struct     map[string]interface{}
	Session    *resty.Request
	cookiePath string
	cookies    []map[string]interface{}
	proxy      string
}

func NewConversation(cookiePath string, cookies []map[string]interface{}, proxy string) *Conversation {
	con := &Conversation{
		Struct: map[string]interface{}{
			"conversationId":        "",
			"clientId":              "",
			"conversationSignature": "",
			"result": map[string]string{
				"value":   "Success",
				"message": "",
			},
		},
		cookiePath: cookiePath,
		cookies:    cookies,
		proxy:      proxy,
	}
	return con
}

func (con *Conversation) Init() error {
	// session
	client := resty.New()
	if con.proxy != "" {
		client.SetProxy(con.proxy)
	}
	client.SetTimeout(time.Second * 30)
	con.Session = client.R().
		SetHeaders(HEADERS_INIT_CONVER)

	// set cookies
	if con.cookies == nil {
		con.cookies = []map[string]interface{}{}
	}
	if len(con.cookies) == 0 {
		b, err := ioutil.ReadFile(con.cookiePath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, &con.cookies)
		if err != nil {
			return err
		}
	}
	for _, item := range con.cookies {
		con.Session.SetCookie(&http.Cookie{
			Name:  item["name"].(string),
			Value: item["value"].(string),
		})
	}
	// Send GET request
	url := os.Getenv("BING_PROXY_URL")
	if url == "" {
		url = "https://edgeservices.bing.com/edgesvc/turing/conversation/create"
	}
	resp, err := con.Session.Get(url)
	if err != nil {
		return fmt.Errorf("init conversation err:%s", err.Error())
	}
	if resp.StatusCode() != 200 {
		log.Println(resp.String())
		return fmt.Errorf("conversation authentication failed: %d", resp.StatusCode())
	}

	// log.Println("conversation info:", resp.String())
	j := gjson.Parse(resp.String())
	if j.Get("result.value").String() == "UnauthorizedRequest" {
		return fmt.Errorf("conversation err:%s", j.Get("result.message").String())
	}
	err = json.Unmarshal([]byte(resp.String()), &con.Struct)
	if err != nil {
		return fmt.Errorf("init conversation err:%s", err.Error())
	}
	return nil

}
