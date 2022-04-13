package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

//func main () {
	//get the filename
	fileName := "happiness.txt"

	//read the file
	byteslice, err := ioutil.ReadFile(fileName)

	if err != nil {

		log.Fatal(err)
	}
	//convert slice of bytes to string
	text := string(byteslice)

	words := strings.Fields(text)

	//create frequency map to store word count
	frequency_map := make(map[string]int)
	

	for _, words := range words {

		frequency_map[words]++
	}

	//get all the keys from the map
	keys := make([]string, len(frequency_map))

	for key := range frequency_map {

		keys = append(keys, key)
	}

	//and sort them
	sort.SliceStable(keys, func(i, j int) bool {

		return frequency_map[keys[i]] > frequency_map[keys[j]]
	})

	//after sorting,get the first 10 words with their count
	for index, key := range keys {
        return_value := make(map[string]int)
		return_value[key] = frequency_map[key]

		if index == 9 {
			break
		}
		return return_value
	}
}
