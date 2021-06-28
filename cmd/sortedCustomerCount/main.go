package main

import (
	"fmt"
	"github.com/IllicLanthresh/TeamworkGoTests/internal/customerimporter"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		panic("unexpected arguments")
	} else if len(os.Args) != 2 {
		panic("missing filepath")
	}
	csvPath := os.Args[1]
	importer, err := customerimporter.NewCsvCustomerImporter(csvPath, "email")
	if err != nil {
		panic(err)
	}

	customerCountByDomain, err := importer.CustomerCountByDomain()
	if err != nil {
		panic(err)
	}
	for _, domain := range customerCountByDomain {
		fmt.Printf("%s(%d)\n", domain.Domain, domain.CustomerCount)
	}
}
