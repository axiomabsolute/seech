package sqlite

import (
	"database/sql"
	"log"
	"os"

	"github.com/gchaincl/dotsql"
	_ "github.com/mattn/go-sqlite3"
)

func createDb(indexPath string) {
	db, err := sql.Open("sqlite3", indexPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dot, err := dotsql.LoadFromFile("sql/trigrams.sql")
	if err != nil {
		log.Fatal("Unable to load dotsql SQL file: ", err)
	}
	_, err = dot.Exec(db, "create-table")
	if err != nil {
		log.Fatal("Unable to create table", err)
	}
	log.Println("Created base table")

	_, err = dot.Exec(db, "create-trigram-table")
	if err != nil {
		log.Fatal("Unable to create trigram index", err)
	}
	log.Println("Created trigram index")

	_, err = dot.Exec(db, "create-insert-trigger")
	if err != nil {
		log.Fatal("Unable to create insert trigger", err)
	}
	log.Println("Created insert trigger")

	_, err = dot.Exec(db, "create-update-trigger")
	if err != nil {
		log.Fatal("Unable to create update trigger", err)
	}
	log.Println("Created update trigger")

	_, err = dot.Exec(db, "create-delete-trigger")
	if err != nil {
		log.Fatal("Unable to create delete trigger", err)
	}
	log.Println("Created delete trigger")
}

func CheckAndCreate(indexPath string) {
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		createDb(indexPath)
	}
}

func TrigramAddToIndex(indexPath string, filePath string, lineNumber int, doc string) {
	db, err := sql.Open("sqlite3", indexPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := "INSERT INTO sources (file_path, line_number, doc) VALUES (?, ?, ?);"
	_, err = db.Exec(query, filePath, lineNumber, doc)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Print succeeded")
	}
}

func TrigramSearch(indexPath string, searchTerm string) {
	db, err := sql.Open("sqlite3", indexPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := "SELECT sources.file_path, sources.line_number, sources.doc FROM trigrams(?) JOIN sources ON sources.id = trigrams.rowid"
	rows, err := db.Query(query, searchTerm)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var (
			filePath   string
			lineNumber int
			doc        string
		)
		if err := rows.Scan(&filePath, &lineNumber, &doc); err != nil {
			log.Fatal(err)
		}
		log.Printf("Result(%s, %d, %s)\n", filePath, lineNumber, doc)
	}

}

func TrigramRemoveFromIndex(indexPath string, filePath string) {
}

func TrigramClear(indexPath string) {
}
