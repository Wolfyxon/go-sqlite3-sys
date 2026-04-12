package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "sqltest/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "test.db")

	if err != nil {
		log.Fatal(err)
	}

	exec(db, "CREATE TABLE IF NOT EXISTS test (test VARCHAR(64))")
	exec(db, "INSERT INTO test (test) VALUES ('hello')")

	query(db, "SELECT * FROM test")

	db.Close()
}

func exec(db *sql.DB, command string) {
	fmt.Println(command)

	stmt, err := db.Prepare(command)

	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec()

	if err != nil {
		log.Fatal(err)
	}

	affectedRows, err := res.RowsAffected()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Affected rows: %d\n", affectedRows)
}

type testRow struct {
	test string
}

func query(db *sql.DB, command string) {
	fmt.Println(command)
	stmt, err := db.Prepare(command)

	if err != nil {
		log.Fatal("Prep err ", err)
	}

	rows, err := stmt.Query()

	if err != nil {
		log.Fatal("Query err ", err)
	}

	columns, err := rows.Columns()

	if err != nil {
		log.Fatal("Column err ", err)
	}

	for _, column := range columns {
		fmt.Printf("%s | ", column)
	}

	fmt.Println("")

	for rows.Next() {
		var testStr string

		rows.Scan(&testStr)
		fmt.Println(testStr)
	}

	rows.Close()
}
