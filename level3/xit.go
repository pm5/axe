package main

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	re "regexp"
)

func GetSessionID(r *http.Response) (sessID string) {
	sessRE, _ := re.Compile("PHPSESSID=(.*?);")
	return sessRE.FindStringSubmatch(r.Cookies()[0].String())[1]
}

type Chief map[string]string

func GetData(body []byte) (output []Chief, err error) {
	colNames := []string{"town", "village", "name"}
	rowRE, _ := re.Compile("(?sU)<tr>(.*)</tr>")
	colRE, _ := re.Compile("(?sU)<td>(.*)</td>")
	for index, row := range rowRE.FindAllSubmatch(body, -1) {
		if index == 0 {
			continue
		}
		c := Chief{}
		for i, col := range colRE.FindAllSubmatch(row[1], -1) {
			c[colNames[i]] = string(col[1])
		}
		output = append(output, c)
	}
	return
}

func GetNextPage(url string, sessID string) (output []Chief, err error) {
	url = url + "?page=next"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "PHPSESSID="+sessID)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	return GetData(b)
}

func GetAllData(url string) (output []Chief, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	sessID := GetSessionID(resp)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	output, err = GetData(b)
	if err != nil {
		return
	}

	for i := 2; i <= 76; i++ {
		d, err := GetNextPage(url, sessID)
		if err != nil {
			return output, err
		}
		output = append(output, d...)
	}
	return
}

func main() {
	lv3URL := "http://axe-level-1.herokuapp.com/lv3/"
	d, err := GetAllData(lv3URL)
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
