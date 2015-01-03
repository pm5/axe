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
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Referer", url)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.9; rv:36.0) Gecko/20100101 Firefox/36.0")
	resp, err := client.Do(req)
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
		return nil, err
	}
	return GetData(b)
}

func GetAllData(url string) (output []Chief, err error) {
	var d []Chief
	for p := 1; p <= 24; p++ {
		d, err = GetSinglePageData(fmt.Sprintf(url+"?page=%d", p))
		if err != nil {
			return nil, err
		}
		output = append(output, d...)
	}
	return
}

func main() {
	d, err := GetAllData("http://axe-level-4.herokuapp.com/lv4/")
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
