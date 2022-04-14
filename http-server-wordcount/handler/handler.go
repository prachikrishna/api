package handler

import (
	"fmt"
	//"http"

	//"http-server-wordcount/happiness.txt"

	"net/http"
	"path"
	"sort"
	"strings"
)

func WordCountHandler(w http.ResponseWriter, r *http.Request) {
	{
		//fileName := "happiness.txt"
		//read the file
		//byteslice, err := ioutil.ReadFile(fileName)

		/*if err != nil {

			log.Fatal(err)
		}*/
		//convert slice of bytes to string
		//text := string(byteslice)
		text := path.Base(r.URL.Path)
		words := strings.Fields(text)

		//create frequency map to store word count
		frequency_map := make(map[string]int)
		return_value := make(map[string]int)

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

			return_value[key] = frequency_map[key]
			fmt.Fprintf(w, "%s %d\n", key, return_value[key])

			if index == 9 {
				break
			}

		}
		//fmt.Println("-------------------")
		//fmt.Println(return_value)
		//return "success"
	}

}

//url : http://localhost:8080/wordcount/we are interns  we are learning go and hope to finish the assignments soon.Bye,see u soon
/*output:
are 2
we 2
hope 1
to 1
u 1
interns 1
learning 1
go 1
the 1
finish 1
*/
