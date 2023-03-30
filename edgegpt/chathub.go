package edgegpt

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

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

	// Buffered channel of outbound messages.
	sendChan chan []byte
	recvChan chan []byte
	answer   chan *Answer
}

func NewChatHub(addr, path string, conversation *Conversation) *ChatHub {
	chathub := &ChatHub{
		addr:     addr,
		path:     path,
		done:     make(chan struct{}),
		sendChan: make(chan []byte, 9999),
		recvChan: make(chan []byte, 9999),
		answer:   make(chan *Answer, 9999),
		request: NewChatHubRequest(
			conversation.Struct["conversationSignature"].(string),
			conversation.Struct["clientId"].(string),
			conversation.Struct["conversationId"].(string),
			0,
		),
	}
	return chathub
}

func (chathub *ChatHub) Init() error {
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
	// defer c.Close()
	chathub.ws = c
	return nil
}

func (chathub *ChatHub) Start() {
	// reveive
	go func() {
		defer close(chathub.done)
		for {
			_, message, err := chathub.ws.ReadMessage()
			if err != nil {
				log.Println("read err:", err)
				return
			}
			// log.Printf("recv: %s", message)
			chathub.recvChan <- message
		}
	}()
	// send
	for {
		select {
		case <-chathub.done:
			log.Println("send done")
			return
		case msg := <-chathub.sendChan:
			err := chathub.ws.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("write err:", err)
				return
			}
		}
	}
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
	// await self.wss.recv()
	return nil
}

// Ask a question to the bot
func (chathub *ChatHub) askStream(prompt string, conversationStyle ConversationStyle, timeout time.Duration) error {
	defer close(chathub.answer)

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
		select {
		case <-chathub.done:
			return nil
		case msg := <-chathub.recvChan:
			objects := strings.Split(string(msg), DELIMITER)
			for _, obj := range objects {
				if obj == "" {
					continue
				}
				answer := NewAnswer(obj)
				mType := answer.Type()
				switch mType {
				case 1:
					chathub.answer <- answer
				case 2:
					chathub.answer <- answer
					answer.Done()
					return nil
				default:
					log.Println(obj)
				}
			}
		case <-time.After(timeout):
			return fmt.Errorf("timeout")
		}
	}
	return nil
}

func (chathub *ChatHub) Answer() chan *Answer {
	return chathub.answer
}
