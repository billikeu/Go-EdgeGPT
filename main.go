package main

import (
	"log"
	"time"

	"github.com/billikeu/Go-EdgeGPT/edgegpt"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	bot := edgegpt.NewChatBot("cookie.json", []map[string]interface{}{}, "http://127.0.0.1:10809")
	err := bot.Init()
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		err := bot.Ask("hi", edgegpt.Creative)
		if err != nil {
			log.Println(err)
		}
	}()
	for {
		answer := <-bot.Answer()
		if answer == nil {
			log.Println("answer nil, end")
			break
		}
		text := answer.Text()
		if text == "" {
			log.Println(answer.Raw())
			continue
		}
		log.Println(text)
	}
	time.Sleep(time.Second * 5)
}
