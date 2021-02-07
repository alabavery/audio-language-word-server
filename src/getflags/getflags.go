package getflags

import (
	"flag"
	"fmt"
)

// FlagValues are the variables file paths necessary for the program
type FlagValues struct {
	Words string
}

// GetFlags gets command line flags
func GetFlags() *FlagValues {
	wordsPathPtr := flag.String("words", "", "the path to the directory containing words")
	flag.Parse()

	wordsPath := *wordsPathPtr

	if wordsPath == "" {
		fmt.Println("Must provide the following flags:")
		flag.PrintDefaults()
		panic("missing flags")
	}
	return &FlagValues{
		Words: wordsPath,
	}
}
