package edgegpt

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

// http.Header{}
var HEADERS = map[string]string{
	"accept":                      "application/json",
	"accept-language":             "en-US,en;q=0.9",
	"content-type":                "application/json",
	"sec-ch-ua":                   `"Not_A Brand";v="99", "Microsoft Edge";v="110", "Chromium";v="110"`,
	"sec-ch-ua-arch":              `"x86"`,
	"sec-ch-ua-bitness":           `"64"`,
	"sec-ch-ua-full-version":      `"109.0.1518.78"`,
	"sec-ch-ua-full-version-list": `"Chromium";v="110.0.5481.192", "Not A(Brand";v="24.0.0.0", "Microsoft Edge";v="110.0.1587.69"`,
	"sec-ch-ua-mobile":            "?0",
	"sec-ch-ua-model":             "",
	"sec-ch-ua-platform":          `"Windows"`,
	"sec-ch-ua-platform-version":  `"15.0.0"`,
	"sec-fetch-dest":              "empty",
	"sec-fetch-mode":              "cors",
	"sec-fetch-site":              "same-origin",
	"x-ms-client-request-id":      GetUuidV4(),
	"x-ms-useragent":              "azsdk-js-api-client-factory/1.0.0-beta.1 core-rest-pipeline/1.10.0 OS/Win32",
	"Referer":                     "https://www.bing.com/search?q=Bing+AI&showconv=1&FORM=hpcodx",
	"Referrer-Policy":             "origin-when-cross-origin",
	"x-forwarded-for":             GetRandomIp(),
}

// Chat API
type ChatHub struct {
	// hub *Hub
	addr string
	path string
	done chan struct{}
	// The websocket connection.
	ws      *websocket.Conn
	request *ChatHubRequest
}

func NewChatHub(addr, path string, conversation *Conversation) *ChatHub {
	chathub := &ChatHub{
		addr: addr,
		path: path,
		done: make(chan struct{}),
		request: NewChatHubRequest(
			conversation.Struct["conversationSignature"].(string),
			conversation.Struct["clientId"].(string),
			conversation.Struct["conversationId"].(string),
			0,
		),
	}
	return chathub
}

func (chathub *ChatHub) newConnect() error {
	u := url.URL{
		Scheme: "wss",
		Host:   chathub.addr,
		Path:   chathub.path,
	}
	headers := http.Header{}
	for key, value := range HEADERS {
		headers.Add(key, value)
	}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		log.Fatal("dial:", err)
	}
	chathub.ws = c
	return nil
}

func (chathub *ChatHub) Close() error {
	return chathub.ws.Close()
}

func (chathub *ChatHub) initialHandshake() error {
	msg, err := appendIdentifier(map[string]interface{}{"protocol": "json", "version": 1})
	if err != nil {
		return fmt.Errorf("initialHandshake err: %s", err.Error())
	}
	err = chathub.ws.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return fmt.Errorf("initialHandshake write err: %s", err.Error())
	}
	_, message, err := chathub.ws.ReadMessage()
	if err != nil {
		return fmt.Errorf("initialHandshake err: %v %s", message, err.Error())
	}
	return nil
}

// Ask a question to the bot
func (chathub *ChatHub) askStream(prompt string, conversationStyle ConversationStyle, callback func(answer *Answer)) error {
	err := chathub.initialHandshake()
	if err != nil {
		return err
	}
	log.Println("initialHandshake success")
	// Construct a ChatHub request
	chathub.request.Update(prompt, conversationStyle)
	// Send request
	msg, err := appendIdentifier(chathub.request.Struct)
	if err != nil {
		return fmt.Errorf("appendIdentifier request struct err: %s", err.Error())
	}
	chathub.ws.WriteMessage(websocket.TextMessage, []byte(msg))
	// var final bool = false
	for {
		_, message, err := chathub.ws.ReadMessage()
		if err != nil {
			return err
		}
		if string(message) == "" {
			return nil
		}
		answer := NewAnswer(string(message))
		if callback != nil {
			callback(answer)
		}
		if answer.IsDone() {
			return nil
		}
		if answer.Type() != 1 && answer.Type() != 2 {
			log.Println(answer.Raw())
		}
	}
}
