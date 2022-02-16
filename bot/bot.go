package main

import (
	"github.com/antchfx/htmlquery"
	"log"
	"os"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle(tele.OnText, func(c tele.Context) error {
		var (
			text   = c.Text()
			prefix = "https://open.spotify.com/track/"
		)

		if strings.HasPrefix(text, prefix) {
			doc, err := htmlquery.LoadURL(text)
			if err != nil {
				log.Print(err)
			}

			metaList, err := htmlquery.QueryAll(doc, "/html/head/meta")
			if err != nil {
				log.Print(err)
			}

			var title, info string
			for _, meta := range metaList {
				if meta.Attr[0].Val == "og:title" {
					title = meta.Attr[1].Val
					continue
				}
				if meta.Attr[0].Val == "og:description" {
					info = meta.Attr[1].Val
				}
			}
			return c.Send(title + " " + info)
		}

		return nil
	})

	b.Start()
}
