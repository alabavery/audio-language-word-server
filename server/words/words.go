package words

import (
	"audio-language/words/server/dbwrapper"
	"audio-language/words/server/rediscli"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Lemma struct {
	ID          int    `json:"id"`
	Word        string `json:"word"`
	Definitions string `json:"definitions"`
}

type PartOfSpeech struct {
	ID           int     `json:"id"`
	PartOfSpeech string  `json:"part_of_speech"`
	Lemmas       []Lemma `json:"lemmas"`
}

type Word struct {
	ID            int            `json:"id"`
	Word          string         `json:"word"`
	PartsOfSpeech []PartOfSpeech `json:"parts_of_speech"`
}

type WordString struct {
	ID   int
	Word string
}

type tempDenormalized struct {
	ID            int
	Word          string
	PartsOfSpeech map[int]PartOfSpeech
}

/*
InitWords goes through the words directory and saves every word json to redis under the word.
It returns a sorted list of the words that will be held in memory through the server's life.
Used on server startup.
*/
func InitWords(db *dbwrapper.DBWrapper, r *rediscli.WordRedisCli) []WordString {
	raw := db.GetAll()

	// first, we must unique the raw join data from sql
	denormalizedMap := make(map[int]tempDenormalized)
	for _, r := range raw {
		wordID := r.WordID

		t, ok := denormalizedMap[wordID]
		if !ok {
			t = tempDenormalized{
				ID:            r.WordID,
				Word:          r.WordString,
				PartsOfSpeech: make(map[int]PartOfSpeech),
			}
			denormalizedMap[wordID] = t
		}
		partOfSpeechID := r.PartOfSpeechID
		pos, ok := denormalizedMap[wordID].PartsOfSpeech[partOfSpeechID]
		if !ok {
			pos = PartOfSpeech{
				ID:           partOfSpeechID,
				PartOfSpeech: r.PartOfSpeech,
			}
		}
		pos.Lemmas = append(pos.Lemmas, Lemma{
			ID:          r.LemmaID,
			Word:        r.LemmaString,
			Definitions: r.Definitions,
		})
		denormalizedMap[wordID].PartsOfSpeech[partOfSpeechID] = pos
	}
	wordStringsUnique := make(map[int]WordString, 0)
	// redis should hold { [wordID]: json(Word) }
	for _, entry := range denormalizedMap {
		wordID := entry.ID
		word := entry.Word
		wordStringsUnique[wordID] = WordString{ID: wordID, Word: word}

		cacheData := Word{ID: wordID, Word: word, PartsOfSpeech: make([]PartOfSpeech, 0)}

		for _, pos := range entry.PartsOfSpeech {
			cacheData.PartsOfSpeech = append(cacheData.PartsOfSpeech, pos)
		}
		j, err := json.Marshal(cacheData)
		if err != nil {
			log.Fatal("could not marshal json for redis")
		}
		r.Set(rediscli.WordKeyFromId(wordID), j)
	}
	wordStrings := make([]WordString, 0)
	for _, wordString := range wordStringsUnique {
		wordStrings = append(wordStrings, wordString)
	}
	sort.Slice(wordStrings, func(i, j int) bool {
		return wordStrings[i].Word < wordStrings[j].Word
	})
	return wordStrings
}

// SearchWords searches the words directory for words that start with the searched string
func SearchWords(searched string, wordList *[]WordString, cli *rediscli.WordRedisCli) *[]Word {
	matches := getMatches(searched, wordList)
	fmt.Println(matches)

	var w []Word
	for _, match := range matches {
		bytes, found := cli.Get(rediscli.WordKeyFromId(match.ID))
		if !found {
			log.Fatal(fmt.Sprintf("word %v present in word list but not in redis", match))
		}
		var word Word
		err := json.Unmarshal(bytes, &word)
		if err != nil {
			log.Fatal("could not get word from redis json")
		}
		w = append(w, word)
	}
	return &w
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

func getMatches(searched string, wordList *[]WordString) []WordString {
	var containSearched []WordString
	matchesBegun := false
	for _, s := range *wordList {
		if strings.HasPrefix(s.Word, searched) {
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
