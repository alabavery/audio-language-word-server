package dbwrapper

import (
	"database/sql"
	"log"

	// lets sql recognize pg driver
	_ "github.com/lib/pq"
)

// DBWrapper around sql.DB. Should pass around pointer.
type DBWrapper struct {
	db *sql.DB
}

type RawAll struct {
	WordID         int
	WordString     string
	PartOfSpeechID int
	PartOfSpeech   string
	LemmaID        int
	LemmaString    string
	Definitions    string
}

// Open the initial db connection
func Open(connection string) *DBWrapper {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	return &DBWrapper{db: db}
}

// type wordRecord struct {
// 	ID     int
// 	String string
// }

// type partOfSpeechRecord struct {
// 	ID           int
// 	Word         int
// 	PartOfSpeech string
// }

// type lemmaRecord struct {
// 	ID           int
// 	PartOfSpeech int
// 	Word         int
// 	Definitions  string
// }

func (db *DBWrapper) GetAll() []RawAll {
	rows, err := db.db.Query(`
		SELECT
			w.id,
			w.string,
			p.id,
			p.part_of_speech,
			l.id,
			w2.string,
			l.definitions
		FROM words w
		JOIN parts_of_speech p ON w.id = p.word
		JOIN lemmas l ON p.id = l.part_of_speech
		LEFT JOIN words w2 ON l.word = w2.id;
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	res := make([]RawAll, 0)
	for rows.Next() {
		var raw RawAll
		if err = rows.Scan(
			&raw.WordID,
			&raw.WordString,
			&raw.PartOfSpeechID,
			&raw.PartOfSpeech,
			&raw.LemmaID,
			&raw.LemmaString,
			&raw.Definitions,
		); err != nil {
			log.Fatal(err)
		}
		res = append(res, raw)
	}
	return res
}
