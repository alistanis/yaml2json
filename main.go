package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	inFile  string
	outFile string
)

func init() {
	flag.StringVar(&inFile, "f", "", "The file to read from")
	flag.StringVar(&outFile, "o", "", "Optional: The output file to write to.")
}

func main() {
	flag.Parse()
	if inFile == "" {
		fmt.Fprintln(os.Stderr, "Missing required argument -f")
	}

	var i interface{}
	data, err := ioutil.ReadFile(inFile)
	if err != nil {
		exitError(err)
	}
	err = yaml.Unmarshal(data, &i)
	if err != nil {
		exitError(err)
	}

	data, err = json.Marshal(i)
	if err != nil {
		exitError(err)
	}

	if outFile != "" {
		ioutil.WriteFile(outFile, data, 0644)
	} else {
		fmt.Print(string(data))
	}

}

func exitError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(-1)
}
