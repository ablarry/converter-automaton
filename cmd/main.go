package main

import (
	"github.com/ablarry/converter-automaton/pkg/mapper"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		filePath := os.Args[1]
		strEval := os.Args[2]
		log.Print("FilePath:" + filePath)
		p, err := mapper.MapperFileToPA(filePath)
		if err != nil {
			log.Fatal(err)
		}
		p.Build()
		log.Println(p)
		log.Println(p.Find(strEval, ""))
	} else {
		log.Print("Add parameters")
	}
}
