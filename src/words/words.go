package words

import (
	"audio-language/words/server/rediscli"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ninetypercentlanguage/misc/files"

	"github.com/ninetypercentlanguage/word-utils/combined"
)

// Word is the contents of a single word's file
type Word struct {
	Word          string           `json:"word"`
	PartsOfSpeech combined.Content `json:"parts_of_speech"`
}

// SearchWords searches the words directory for words that start with the searched string
func SearchWords(searched string, wordList *[]string, cli *rediscli.WordRedisCli) *[]Word {
	matches := getMatches(searched, wordList)

	var w []Word
	for _, match := range matches {
		bytes, found := cli.Get(match)
		if !found {
			panic(fmt.Sprintf("word %v present in word list but not in redis", match))
		}
		var pos []combined.ContentItem
		err := json.Unmarshal(bytes, &pos)
		if err != nil {
			panic(err)
		}
		w = append(w, Word{Word: match, PartsOfSpeech: pos})
	}
	return &w
}

/*
InitWords goes through the words directory and saves every word json to redis under the word.
It returns a sorted list of the words that will be held in memory through the server's life.
Used on server startup.
*/
func InitWords(wordsDir string, r *rediscli.WordRedisCli) *[]string {
	wordPaths := getWordFiles(wordsDir)
	var wordStrings []string
	for _, f := range wordPaths {
		word := wordFromWordPath(f)
		wordStrings = append(wordStrings, word)
		content, _ := files.ReadFileThatMayNotExist(f)
		r.Set(word, content)
	}
	fmt.Printf("\nAdded %v words to redis\n", len(wordPaths))

	sort.Strings(wordStrings)
	return &wordStrings
}

func wordFromWordPath(path string) string {
	d := strings.Split(path, "/")
	fileName := d[len(d)-1]
	return strings.Split(fileName, ".")[0]
}

func getWordFiles(wordsDir string) []string {
	var files []string
	err := filepath.Walk(wordsDir, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func getMatches(searched string, wordList *[]string) []string {
	var containSearched []string
	matchesBegun := false
	for _, s := range *wordList {
		if strings.HasPrefix(s, searched) {
			containSearched = append(containSearched, s)
			matchesBegun = true
		} else {
			if matchesBegun {
				break
			}
		}
	}
	return containSearched
}
