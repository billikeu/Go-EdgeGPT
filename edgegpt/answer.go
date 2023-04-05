package edgegpt

import (
	"encoding/json"
	"log"
	"time"

	"github.com/tidwall/gjson"
)

type SuggestedResponses struct {
	Text        string    `json:"text"`
	Author      string    `json:"author"`
	CreatedAt   time.Time `json:"createdAt"`
	Timestamp   time.Time `json:"timestamp"`
	MessageID   string    `json:"messageId"`
	MessageType string    `json:"messageType"`
	Offense     string    `json:"offense"`
	Feedback    struct {
		Tag       interface{} `json:"tag"`
		UpdatedOn interface{} `json:"updatedOn"`
		Type      string      `json:"type"`
	} `json:"feedback"`
	ContentOrigin string      `json:"contentOrigin"`
	Privacy       interface{} `json:"privacy"`
}

type Answer struct {
	msg  string
	j    gjson.Result
	done bool
}

func NewAnswer(msg string) *Answer {
	answer := &Answer{
		msg:  msg,
		j:    gjson.Parse(msg),
		done: false,
	}
	return answer
}

func (answer *Answer) Raw() string {
	return answer.msg
}

func (answer *Answer) Text() string {
	if answer.Type() == 2 {
		messages := answer.j.Get("item.messages").Array()
		if len(messages) < 1 {
			return ""
		}
		lastMsg := messages[len(messages)-1].Get("text").String()
		if lastMsg == "" {
			lastMsg = messages[len(messages)-1].Get("hiddenText").String()
		}
		if lastMsg == "" {
			lastMsg = messages[len(messages)-1].Get("spokenText").String()
		}
		return lastMsg
	}
	arguments := answer.j.Get("arguments").Array()
	if len(arguments) < 1 {
		return ""
	}
	messages := arguments[0].Get("messages").Array()
	if len(messages) < 1 {
		return ""
	}
	adaptiveCards := messages[0].Get("adaptiveCards").Array()
	if len(adaptiveCards) < 1 {
		return ""
	}
	body := adaptiveCards[0].Get("body").Array()
	if len(body) < 1 {
		return ""
	}
	text := body[0].Get("text").String()
	return text
}

func (answer *Answer) Type() int64 {
	return answer.j.Get("type").Int()
}

func (answer *Answer) SetDone() {
	answer.done = true
}

func (answer *Answer) IsDone() bool {
	return answer.Type() == 2
}

func (answer *Answer) SuggestedRes() []SuggestedResponses {
	suggest := []SuggestedResponses{}
	if answer.Type() == 2 {
		messages := answer.j.Get("item.messages").Array()
		if len(messages) < 1 {
			return suggest
		}
		suggestedResponses := messages[len(messages)-1].Get("suggestedResponses").Array()
		for _, item := range suggestedResponses {
			s := SuggestedResponses{}
			err := json.Unmarshal([]byte(item.String()), &s)
			if err != nil {
				log.Println(err)
				continue
			}
			suggest = append(suggest, s)
		}
		return suggest
	}
	if answer.Type() == 1 {
		arguments := answer.j.Get("arguments").Array()
		if len(arguments) < 1 {
			return suggest
		}
		messages := arguments[0].Get("messages").Array()
		if len(messages) < 1 {
			return suggest
		}
		suggestedResponses := messages[len(messages)-1].Get("suggestedResponses").Array()
		for _, item := range suggestedResponses {
			s := SuggestedResponses{}
			err := json.Unmarshal([]byte(item.String()), &s)
			if err != nil {
				log.Println(err)
				continue
			}
			suggest = append(suggest, s)
		}
		return suggest
	}
	return suggest
}

func (answer *Answer) MaxNumUserMessages() int64 {
	return answer.j.Get("item.throttling.maxNumUserMessagesInConversation").Int()
}

func (answer *Answer) NumUserMessages() int64 {
	return answer.j.Get("item.throttling.numUserMessagesInConversation").Int()
}
