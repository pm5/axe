package main

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	re "regexp"
	"strconv"
)

type Grade struct {
	Name   string         `json:"name"`
	Grades map[string]int `json:"grades"`
}

func GetHTMLBody(url string) (output []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func GetData(b []byte) (output []Grade, err error) {
	rowRE, err := re.Compile("(?sU)<tr>(.*)</tr>")
	colRE, err := re.Compile("(?sU)<td>(.*)</td>")
	colNames := []string{"姓名", "國語", "數學", "自然", "社會", "健康教育"}
	for index, row := range rowRE.FindAllSubmatch(b, -1) {
		if index == 0 {
			continue
		}
		g := Grade{Name: "", Grades: make(map[string]int, len(colNames)-1)}
		for index, col := range colRE.FindAllSubmatch(row[1], -1) {
			if index == 0 {
				g.Name = string(col[1])
			} else {
				g.Grades[colNames[index]], err = strconv.Atoi(string(col[1]))
			}
		}
		output = append(output, g)
	}
	return
}

func main() {
	b, err := GetHTMLBody("http://axe-level-1.herokuapp.com/")
	if err != nil {
		log.Fatal(err)
	}
	//os.Stdout.Write(b)

	d, err := GetData(b)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(o)

	j, err := json.Marshal(d)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("output.json", j, 0600)
	if err != nil {
		log.Fatal(err)
	}
	//os.Stdout.Write(j)
}
