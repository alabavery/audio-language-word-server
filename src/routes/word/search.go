package word

import (
	"audio-language/words/server/rediscli"
	"audio-language/words/server/words"
	"encoding/json"
	"net/http"
)

// Returns a route that searches for the queried "string".
// It matches anything that starts with the searched string.
// It returns matches in alphabetical order.
func getSearchWordsRoute(wordList *[]string, cli *rediscli.WordRedisCli) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searched := r.FormValue("string")
		words := words.SearchWords(searched, wordList, cli)
		var res json.RawMessage
		if len(*words) == 0 {
			res = json.RawMessage("[]")
		} else {
			res, _ = json.Marshal(*words)
		}
		w.Write(res)
	}
}
