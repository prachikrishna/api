package core

//28-04-22 leaving after adding list type here and in mysvc,everything is fine till here
//add the main.go and client part

import (
	word_pb "prachi/grpc"
	"prachi/mysvc"
	"strings"

	"sort"
)

type service struct {
}

func NewService() mysvc.Service {
	return &service{}
}

//func (s *service) GetWCount(text string) (result map[string]uint32, err error) {
func (s *service) GetWCount(text string) []*word_pb.Word {
	words1 := strings.Fields(text)

	frequency_map := make(map[string]uint64)
	//return_value := make(map[string]uint64)

	for _, words := range words1 {
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
	var words []*word_pb.Word
	for key, value := range frequency_map {
		words = append(words, &word_pb.Word{Word: key, Count: value})
	}

	sort.Slice(words, func(i, j int) bool {
		return words[i].Count > words[j].Count
	})

	if len(words) > 10 {
		return words[:10]
	}

	//after sorting,get the first 10 words with their count

	return words

}
