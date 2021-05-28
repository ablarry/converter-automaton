package mapper

import (
	"encoding/csv"
	"github.com/ablarry/converter-automaton/pkg/model"
	"log"
	"os"
	"strings"
)

// MapperFileToPA maps the csv file to Transactions struct
func MapperFileToPA(file string) ([]*model.Transition, *model.MetaData, error) {
	rows := ReadFile(file)
	transitions := make([]*model.Transition, 0)
	var metaData *model.MetaData = nil
	for _, v := range rows {
		if v[0] != "METADATA" {
			t := MapperToTransition(v)
			transitions = append(transitions, t)
		} else {
			// Validate that if value is METADATA we need to use another mapper
			metaData = MapperToMetadata(v)
		}
	}
	return transitions, metaData, nil
}

// MapperToTransition map array string to Transition struct
func MapperToTransition(s []string) *model.Transition {
	t := model.Transition{
		InitialState: s[0],
		ReadElement:  s[1],
		PullElement:  strings.Split(s[2], ";")[0],
		FinalState:   strings.Split(s[2], ";")[1],
		PushElement:  s[3],
	}
	return &t
}

// MapperToMetadata map array string to Metadata struct
func MapperToMetadata(s []string) *model.MetaData {
	m := model.MetaData{
		AcceptStates: strings.Split(s[1], ";"),
	}
	return &m
}

// ReadFile accepts file path of csv and return content in [][] string
func ReadFile(file string) [][]string {
	//Open file csv
	f, err := os.Open(file)
	if err != nil {
		log.Fatal("Couldn't open csv file", err)
	}
	defer f.Close()

	csvr := csv.NewReader(f)
	// Disable validation of number of fields
	csvr.FieldsPerRecord = -1
	lines, err := csvr.ReadAll()
	if err != nil {
		log.Fatal("Couldn't open csv file", err)
	}

	return lines
}
