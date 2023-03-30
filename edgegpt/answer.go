package edgegpt

import "github.com/tidwall/gjson"

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

func (answer *Answer) Done() {
	answer.done = true
}
