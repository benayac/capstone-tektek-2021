package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/bradfitz/slice"
)

type Predict struct {
	SourceText string
	ReportList []Report
}

func (p *Predict) Predict() {
	fmt.Printf("Enter Predict Here!\n")
	var wg sync.WaitGroup
	wg.Add(len(p.ReportList))
	for i, r := range p.ReportList {
		go func(r Report, index int) {
			defer wg.Done()
			accuracy := predictModel(p.SourceText, r.Keterangan)
			p.ReportList[index].Accuracy = accuracy
		}(r, i)
		fmt.Printf("Go Routine\n")
	}
	wg.Wait()
}

func (p *Predict) SortByAccuracy() {
	slice.Sort(p.ReportList[:], func(i, j int) bool {
		return p.ReportList[i].Accuracy > p.ReportList[j].Accuracy
	})
}

func predictModel(firstText string, secondText string) float64 {
	fmt.Printf("Enter Here\n")
	url := "http://ml:8001/predict"
	postBody, _ := json.Marshal(map[string]string{
		"first_text":  firstText,
		"second_text": secondText,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Text: %v Predict: %v\n", secondText, string(body))

	accuracy, _ := strconv.ParseFloat(string(body), 8)
	return accuracy
}
