package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	re "regexp"
)

type Chief map[string]string

func GetHTMLBody(url string) (output []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func GetData(b []byte) (output []Chief, err error) {
	rowRE, _ := re.Compile(`(?sU)<tr>(.*)</tr>`)
	colRE, _ := re.Compile(`(?sU)<td>(.*)</td>`)
	colNames := []string{"town", "village", "name"}
	for index, row := range rowRE.FindAllSubmatch(b, -1) {
		if index == 0 {
			continue
		}
		g := Chief{}
		for i, col := range colRE.FindAllSubmatch(row[1], -1) {
			g[colNames[i]] = string(col[1])
		}
		output = append(output, g)
	}
	return
}

func GetSinglePageData(url string) (output []Chief, err error) {
	b, err := GetHTMLBody(url)
	if err != nil {
		log.Fatal(err)
	}
	return GetData(b)
}

func GetAllData(url string) (output []Chief, err error) {
	var d []Chief
	for p := 1; p <= 12; p++ {
		d, err = GetSinglePageData(fmt.Sprintf("http://axe-level-1.herokuapp.com/lv2/?page=%d", p))
		if err != nil {
			log.Fatal(err)
		}
		output = append(output, d...)
	}
	return
}

func main() {
	d, err := GetAllData("http://axe-level-1.herokuapp.com/lv2/")
	if err != nil {
		log.Fatal(err)
	}

	j, err := json.Marshal(d)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("output.json", j, 0600)
	if err != nil {
		log.Fatal(err)
	}
}
