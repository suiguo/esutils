package main

import (
	"log"
	"time"

	esconvert "github.com/suiguo/esutils"
)

type TestMapping struct {
	Mid      string    `json:"mid"`
	Date     time.Time `json:"date"`
	Number   int       `json:"number"`
	Float32  float32   `json:"float32"`
	Float64  float64   `json:"float64"`
	Number64 int64     `json:"number64"`
}
type TestMapping2 struct {
	Mid      string    `json:"mid"`
	Date     time.Time `json:"date"`
	Number   int       `json:"number"`
	Float32  float32   `json:"float32"`
	Float64  float64   `json:"float64"`
	Number64 int64     `json:"number64"`
	NewParam int64     `json:"new_param"`
}

func main() {
	tool, err := esconvert.NewConver(esconvert.V8, esconvert.WithHost("http://127.0.0.1:9200"), esconvert.WithIgnoreAbove(256))
	if err != nil {
		panic(err)
	}
	resp, err := tool.Create("index-test-create3", &TestMapping{})
	log.Printf("Create======>data[%s] error[%v]\n", resp, err)
	resp, err = tool.GetMapping("index-test-create3")
	log.Printf("GetMapping======>data[%s] error[%v]\n", resp, err)
	resp, err = tool.Put("index-test-create3", &TestMapping2{})
	log.Printf("Put======>data[%s] error[%v]\n", resp, err)
	resp, err = tool.GetMapping("index-test-create3")
	log.Printf("GetMapping======>data[%s] error[%v]\n", resp, err)
}
