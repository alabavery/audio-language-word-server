package word

import (
	"audio-language/words/server/rediscli"
	"net/http"

	"github.com/go-chi/chi"
)

func getFindPartOfSpeechRoute(cli *rediscli.WordRedisCli) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "partOfSpeechId")

		val, exists := cli.Get(id)
		if !exists {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		w.Write([]byte(val))
	}
}
