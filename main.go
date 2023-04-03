package main

import (
	"log"

	"github.com/billikeu/Go-EdgeGPT/edgegpt"
)

func callback(answer *edgegpt.Answer) {
	if answer.IsDone() {
		log.Println(answer.NumUserMessages(), answer.MaxNumUserMessages(), answer.Text())
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	bot := edgegpt.NewChatBot("cookie.json", []map[string]interface{}{}, "http://127.0.0.1:10809")
	err := bot.Init()
	if err != nil {
		panic(err)
	}
	err = bot.Ask("give me a joke", edgegpt.Creative, callback)
	if err != nil {
		panic(err)
	}
	err = bot.Ask("It's not funny", edgegpt.Creative, callback)
	if err != nil {
		panic(err)
	}
}
