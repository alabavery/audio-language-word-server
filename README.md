Gets all word data from some directory passed in at runtime.
Expects the directory to be full of files names `<word>.json`.
Expects each of those word files to meet interface of type `audio-language/words/server.partOfSpeech`.
Saves all file contents to redis, keyed under the word.
Using `docker run --name some-redis -d -p 6379:6379 redis:alpine` for redis