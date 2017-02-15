package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"reflect"

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

	var m map[interface{}]interface{}

	data, err := ioutil.ReadFile(inFile)
	if err != nil {
		exitError(err)
	}
	err = yaml.Unmarshal(data, &m)

	nm := make(map[string]interface{})
	recurseMapInterface(m, nm)

	data, err = json.MarshalIndent(nm, "", "	")
	if err != nil {
		exitError(err)
	}

	if outFile != "" {
		ioutil.WriteFile(outFile, data, 0644)
	} else {
		fmt.Print(string(data))
	}

}

func recurseMapInterface(m map[interface{}]interface{}, newMap map[string]interface{}) {
	nm := make(map[string]interface{})
	for k, v := range m {
		nk := ""
		switch t := k.(type) {
		case string:
			nk = t
		case int:
			nk = strconv.Itoa(t)
		case float64:
			nk = strconv.FormatFloat(t, 'E', -1, 64)
		case map[interface{}]interface{}:
			recurseMapInterface(t, newMap)
		}

		var nv interface{}
		switch t := v.(type) {
		case map[interface{}]interface{}:
			m := make(map[string]interface{})
			nm[nk] = m
			recurseMapInterface(t, m)
		case []interface{}:
			recurseArray(nk, t, nm)
		default:
			nv = t
		}
		if nv != nil {
			nm[nk] = nv
		}
	}

	for k, v := range nm {
		newMap[k] = v
	}

}

func recurseArray(k string, slc []interface{}, container interface{}) {
	nslc := make([]interface{}, 0)
	for _, i := range slc {
		var v interface{}
		switch i := i.(type) {
		case []interface{}:
			recurseArray(k, i, &nslc)
		case map[interface{}]interface{}:
			m := make(map[string]interface{})
			nslc = append(nslc, m)
			recurseMapInterface(i, m)
		default:
			v = i
		}
		if v != nil {
			nslc = append(nslc, v)
		}
	}

	rv := reflect.ValueOf(container)
	if rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}
	switch rv.Kind() {
	case reflect.Slice:
		*container.(*[]interface{}) = append(*container.(*[]interface{}), nslc)
	case reflect.Map:
		m := container.(map[string]interface{})
		m[k] = nslc
	}
}

func exitError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(-1)
}
