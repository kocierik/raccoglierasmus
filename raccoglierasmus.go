package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"strconv"
	"strings"
)

const (
	erasmusUrl = "https://almarm.unibo.it/almarm/outgoing/offerteDiScambio.htm;jsessionid=1B386E89CF0D846D3D3C923002746E3D.prod-formazione-almarm-llpp-java-11?execution=e1s1"
)

type Erasmus struct {
	id               int
	universita       string
	nazione          string
	areaDisciplinare string
	docente          string
	posti            int
	lingue           []string
}

func collyError(r *colly.Response, err error) {
	fmt.Fprintln(os.Stderr, "Request URL:", r.Request.URL, "failed with response:", r.StatusCode, "\nError:", err)
}

func main() {
	c := colly.NewCollector()

	c.OnHTML("a[title='dettaglio offerta'][href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println("Visiting", link)

		erasmus := Erasmus{}
		erasmus.id = // Assign a unique ID for each entry, you can use a counter or another method

			e.Request.Visit(link)
		c.OnHTML("div.rigaPari, div.rigaDispari", func(e *colly.HTMLElement) {
			text := strings.TrimSpace(e.Text)
			parts := strings.Split(text, ":")

			if len(parts) != 2 {
				return
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {
			case "Università":
				erasmus.universita = value
			case "Nazione":
				erasmus.nazione = value
			case "Area disciplinare":
				erasmus.areaDisciplinare = value
			case "Docente":
				erasmus.docente = value
			case "Posti":
				posti, err := strconv.Atoi(value)
				if err == nil {
					erasmus.posti = posti
				}
			case "Lingue":
				lingue := strings.Split(value, ",")
				for i := range lingue {
					lingue[i] = strings.TrimSpace(lingue[i])
				}
				erasmus.lingue = lingue
			}
		})

		// Do something with the populated Erasmus struct here
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(collyError)

	c.Visit(erasmusUrl)
}

