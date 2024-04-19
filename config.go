package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

const DEFAULT_DELAY_SECS = 10
const FILE_PATH = "delaySecs.dat"

func loadDelaySecs() int {
	content, err := ioutil.ReadFile(FILE_PATH)
	if err != nil {
		fmt.Printf("E: err while reading %q: %v\n", FILE_PATH, err)
		return DEFAULT_DELAY_SECS
	}

	delaySecs, err := strconv.Atoi(string(content))
	if err != nil {
		fmt.Printf("E: err while parsing %q from %q: %v\n", content, FILE_PATH, err)
		return DEFAULT_DELAY_SECS
	}

	return delaySecs
}

func persistDelaySecs(delaySecs int) {
	content := fmt.Sprintf("%d", delaySecs)
	err := os.WriteFile(FILE_PATH, []byte(content), 0666)
	if err != nil {
		fmt.Printf("E: err while writing %d to %q: %v\n", delaySecs, FILE_PATH, err)
	}
}
