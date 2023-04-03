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