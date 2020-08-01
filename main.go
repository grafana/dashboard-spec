package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Spec map[string]interface{}

func main() {
	s := loadSpec("bundle/7.0/spec.json")
	fmt.Println(s)
}

func loadSpec(file string) (s Spec) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	if err := json.Unmarshal(data, &s); err != nil {
		log.Fatalln(err)
	}
	return
}
