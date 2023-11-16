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

func generateOutput(erasmusList []Erasmus) string {
	b := strings.Builder{}

	// Inizia a costruire la stringa con intestazione
	b.WriteString("= Elenco Erasmus\n:toc:\n\n")

	// Ciclo attraverso gli elementi Erasmus
	for _, erasmus := range erasmusList {
		b.WriteString(fmt.Sprintf("\n== ID: %d\n", erasmus.id))
		b.WriteString(fmt.Sprintf("* Università: %s\n", erasmus.universita))
		b.WriteString(fmt.Sprintf("* Nazione: %s\n", erasmus.nazione))
		b.WriteString(fmt.Sprintf("* Area disciplinare: %s\n", erasmus.areaDisciplinare))
		b.WriteString(fmt.Sprintf("* Docente proponente: %s\n", erasmus.docente))
		b.WriteString(fmt.Sprintf("* Posti disponibili: %d\n", erasmus.posti))

		b.WriteString("* Lingue di accertamento linguistico: ")
		for _, lingua := range erasmus.lingue {
			b.WriteString(fmt.Sprintf("%s ", lingua))
		}
		b.WriteString("\n")
	}

	output := b.String()
	// Qui potresti eventualmente eseguire alcune operazioni di sostituzione se necessario
	// output = replaceRegexForOutput.ReplaceAllString(output, " ")
	return output
}

func main() {
	c := colly.NewCollector()

	erasmusList := []Erasmus{}

	var headers []string

	c.OnHTML("th.iceTblHeader", func(e *colly.HTMLElement) {
		headers = append(headers, strings.TrimSpace(e.Text))
	})

	c.OnHTML("tr.rigaPari, tr.rigaDispari", func(e *colly.HTMLElement) {
		erasmus := Erasmus{}

		e.ForEach("td.colonna", func(i int, td *colly.HTMLElement) {
			text := strings.TrimSpace(td.Text)
			text = strings.ReplaceAll(text, "\u00a0", "")

			header := headers[i]

			switch header {
			case "Id":
				id, err := strconv.Atoi(text)
				if err == nil {
					erasmus.id = id
				}
			case "Università":
				erasmus.universita = text
			case "Nazione":
				erasmus.nazione = text
			case "Area disciplinare":
				erasmus.areaDisciplinare = text
			case "Docente proponente":
				erasmus.docente = text
			case "Posti disponibili":
				posti, err := strconv.Atoi(text)
				if err == nil {
					erasmus.posti = posti
				}
			case "Lingue di accertamento linguistico":
				lingue := strings.Split(text, "\n")
				for i := range lingue {
					lingue[i] = strings.TrimSpace(lingue[i])
				}
				erasmus.lingue = lingue
			}
		})

		erasmusList = append(erasmusList, erasmus)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(collyError)

	c.Visit(erasmusUrl)

	fmt.Println("Collected Erasmus Data:")
	for _, erasmus := range erasmusList {
		fmt.Printf("%+v\n", erasmus)
	}

	output := generateOutput(erasmusList)
	fmt.Println(output)
}
