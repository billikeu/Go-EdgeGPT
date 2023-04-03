# Go-EdgeGPT

Go-EdgeGPT is a New Bing unofficial API developed using Golang. With most chatbot APIs being built on Python, Go-EdgeGPT is unique in its ability to be easily compiled and deployed. It's designed to work seamlessly with your current applications, providing a customizable and reliable chatbot experience.

## Setup

```
go get -u github.com/billikeu/Go-EdgeGPT/edgegpt
```

## Example bot

```golang
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

```

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=billikeu/Go-EdgeGPT&type=Date)](https://star-history.com/#billikeu/Go-EdgeGPT&Date)

## Contributors

This project exists thanks to all the people who contribute.

 <a href="github.com/billikeu/Go-EdgeGPT/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=billikeu/Go-EdgeGPT" />
 </a>

## Reference

Thanks for [EdgeGPT](https://github.com/acheong08/EdgeGPT)