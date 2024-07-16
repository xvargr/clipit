package fileReader

import (
	"encoding/json"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Read(filepath string, result interface{}) {
	file, err := os.Open(filepath)
	check(err)
	defer file.Close()

	parseJSONStream(file, &result)
}

func parseJSONStream(file *os.File, result interface{}) {
	decoder := json.NewDecoder(file)

	for {
		if err := decoder.Decode(result); err != nil {
			if err.Error() == "EOF" {
				break
			}
			check(err)
		}
	}
}
