package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	var err error
	if err = godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	connectURI, connectURIexists := os.LookupEnv("PG_CONN_URI")
	if !connectURIexists {
		return nil, fmt.Errorf("no connection string")
	}
	DB, err := sql.Open("postgres", connectURI)
	if err != nil {
		return nil, fmt.Errorf("connecting to DB failed, %v", err)
	}

	pingErr := DB.Ping()
	if pingErr != nil {
		return nil, fmt.Errorf("connection to DB lost")
	}

	fmt.Printf("connected to postgres\n")
	return DB, nil
}


const initQuery = `CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    batch VARCHAR(63),
    names TEXT,
    project VARCHAR(63),
		group1 TEXT,
		group2 TEXT,
		group3 TEXT,
		group4 TEXT,
		group5 TEXT,
		group6 TEXT,
		group7 TEXT
		);


		INSERT INTO groups (batch, names, project, group1, group2, group3, group4, group5, group6, group7)
		VALUES
    ('WD38', 'Puri,Adam,Miro,Hani,Stefan,Maxine,Stephan,Daniel,Rana,Ali,Amro,Yosra,Daniel O.,Gretchel,Daniil,Anton,Foram,Rohan,Timo', '', '', '', '', '', '', '', ''),
    ('WD38', 'Adam,Puri,Hani,Miro,Stefan,Maxine,Stephan,Daniel,Rana,Ali,Amro,Yosra,Daniel O.,Gretchel,Daniil,Anton,Foram,Rohan,Timo', 'Hackerrank', 'Adam,Puri,Hani,Miro','Stefan,Maxine,Stephan,Daniel', 'Rana,Ali,Amro,Yosra', 'Daniel O.,Gretchel,Daniil', 'Foram,Rohan,Timo', '', ''),
    ('WD38', 'Miro,Stefan,Puri,Hani,Adam,Maxine,Stephan,Daniel,Rana,Ali,Amro,Yosra,Daniel O.,Gretchel,Daniil,Anton,Foram,Rohan,Timo', 'PokeFight', 'Rana,Stephan,Anton,Foram', 'Yosra,Puri,Hani,Miro', 'Rohan,Adam,Daniel O.,Daniil', 'Gretchel,Maxine,Stefan' , 'Timo,Ali,Amro,Daniel', '', '');
`

// Rana,Stephan,Anton,Foram,Yosra,Puri,Hani,Miro,Rohan,Adam,Daniel O.,Daniil,Gretchel,Maxine,Stefan,Timo,Ali,Amro,Daniel
	func InitDB(db *sql.DB) (error) {
		_, err := db.Exec(initQuery)
		if err != nil {
			return fmt.Errorf("table creation failed: %v", err)
		}
		return nil
	}