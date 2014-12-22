package main

import (
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strconv"
)

type Grade struct {
	Name   string         `json:"name"`
	Grades map[string]int `json:"grades"`
}

func main() {
	doc, err := goquery.NewDocument("http://axe-level-1.herokuapp.com/")
	if err != nil {
		log.Fatal(err)
	}

	fields := []string{"國語", "數學", "自然", "社會", "健康教育"}
	lines := doc.Find("tr")
	output := make([]Grade, 0, lines.Length())
	lines.Slice(1, lines.Length()).Each(func(i int, s *goquery.Selection) {
		g := Grade{Name: "", Grades: map[string]int{}}
		g.Name = s.Find("td").First().Text()
		s.Find("td").Slice(1, 6).Each(func(i int, t *goquery.Selection) {
			g.Grades[fields[i]], err = strconv.Atoi(t.Text())
		})
		output = append(output, g)
	})

	b, err := json.Marshal(output)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	out.WriteTo(os.Stdout)
}
