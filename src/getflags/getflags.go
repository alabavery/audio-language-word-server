package getflags

import (
	"flag"
	"fmt"
)

type FlagValues struct {
	DB string
}

func GetFlags() *FlagValues {
	dbConnectionStrPtr := flag.String("db", "", "the connection string for the words database")
	flag.Parse()

	dbConnectionStr := *dbConnectionStrPtr

	if dbConnectionStr == "" {
		fmt.Println("Must provide the following flags:")
		flag.PrintDefaults()
		panic("missing flags")
	}
	return &FlagValues{
		DB: dbConnectionStr,
	}
}
