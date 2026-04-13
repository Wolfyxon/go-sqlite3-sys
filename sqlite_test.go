package sqlite

import (
	"database/sql"
	"log"
	"testing"
)

type testRow struct {
	text   string
	number int
	data   []byte
}

func TestMain(t *testing.T) {
	db, err := sql.Open("sqlite", "a.db")

	if err != nil {
		log.Fatal("Open error: ", err)
	}

	tableRes, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS test (" +
			"text VARCHAR(1024)," +
			"number INTEGER," +
			"data BLOB" +
			")",
	)

	if err != nil {
		log.Fatal("Table create error: ", err)
	}

	rowChanges, err := tableRes.RowsAffected()

	if err != nil {
		log.Fatal("Row changes error: ", err)
	}

	if rowChanges != 0 {
		log.Fatal("Table create row changes not 0: ", rowChanges)
	}

	testNum := 123
	testText := "', 321); INSERT INTO test (text) VALUES ('pwned')"
	testBlob := []byte("beep bop boop beep")
	testBlobLen := len(testBlob)

	insertRes, err := db.Exec("INSERT INTO test (text, number, data) VALUES (?, ?, ?)", testText, testNum, testBlob)

	if err != nil {
		log.Fatal("Insert error: ", err)
	}

	insertChanges, err := insertRes.RowsAffected()

	if err != nil {
		log.Fatal("Insert changes error: ", err)
	}

	if insertChanges != 1 {
		log.Fatal("Insert changes not 1: ", insertChanges)
	}

	rows, err := db.Query("SELECT * FROM test")

	if err != nil {
		log.Fatal("Query error: ", err)
	}

	gotRow := false

	for rows.Next() {
		if gotRow {
			log.Fatal("Got more than 1 rows")
		}

		gotRow = true
		row := testRow{}

		rows.Scan(&row.text, &row.number, &row.data)

		if row.number != testNum {
			log.Fatalf("Test number %d != %d", testNum, row.number)
		}

		if row.text != testText {
			log.Fatalf("Test text '%s' != '%s'", testText, row.text)
		}

		resBlobLen := len(row.data)

		if resBlobLen != testBlobLen {
			log.Fatalf("Length of test blobs don't match. Expected %d got %d", testBlobLen, resBlobLen)
		}

		for i := range testBlobLen {
			if row.data[i] != testBlob[i] {
				log.Fatalf("Blobs don't match")
			}
		}
	}

	if !gotRow {
		log.Fatal("No row returned")
	}

}
