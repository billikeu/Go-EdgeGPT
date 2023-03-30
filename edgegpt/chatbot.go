package edgegpt

import (
	"log"
	"time"
)

// Combines everything to make it seamless
type ChatBot struct {
	cookiePath string
	cookies    []map[string]interface{}
	proxy      string
	chatHub    *ChatHub
	addr       string
	path       string
	Err        error
}

func NewChatBot(cookiePath string, cookies []map[string]interface{}, proxy string) *ChatBot {
	addr := "sydney.bing.com"
	path := "/sydney/ChatHub"
	bot := &ChatBot{
		cookiePath: cookiePath,
		cookies:    cookies,
		proxy:      proxy,
		addr:       addr,
		path:       path,
		chatHub:    nil,
		Err:        nil,
	}
	return bot
}

func (bot *ChatBot) Init() error {
	conversation := NewConversation(bot.cookiePath, bot.cookies, bot.proxy)
	err := conversation.Init()
	if err != nil {
		return err
	}
	log.Println("init conversation success")
	bot.chatHub = NewChatHub(bot.addr, bot.path, conversation)
	err = bot.chatHub.Init()
	if err != nil {
		return err
	}
	log.Println("init chathub success")
	go bot.chatHub.Start()
	log.Println("init chatbot success")
	return nil
}

// Ask a question to the bot
func (bot *ChatBot) Ask(prompt string, conversationStyle ConversationStyle) error {
	defer bot.chatHub.Close()
	bot.Err = bot.chatHub.askStream(prompt, conversationStyle, time.Minute*5)
	return bot.Err
}

func (bot *ChatBot) Close() error {
	return bot.chatHub.Close()
}

// Reset the conversation
func (bot *ChatBot) Reset() error {
	bot.chatHub.Close()
	err := bot.Init()
	if err != nil {
		return err
	}
	return nil
}

func (bot *ChatBot) Answer() chan *Answer {
	return bot.chatHub.Answer()
}
