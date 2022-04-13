package main

import "testing"

func TestWordCount(t *testing.T) {
	
	input := "happiness.txt"
	expectedOutput := "success"

	output := WordCount(input)

	if output != expectedOutput {
		t.Errorf("got %q, wanted %q", output, expectedOutput)
	}

}
