package main

import (
	"audio-language/words/server/getflags"
	"audio-language/words/server/rediscli"
	"audio-language/words/server/routes/word"
	"audio-language/words/server/words"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

/*
"/word" /{id} ?
"/user" /{id}
*/
func main() {
	flagVals := getflags.GetFlags()
	cli := rediscli.GetWordRedisCli()
	wordList := words.InitWords(flagVals.Words, cli)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/words", word.GetWordSubRoute(wordList, cli))

	log.Fatal(http.ListenAndServe(":5000", r))
}
