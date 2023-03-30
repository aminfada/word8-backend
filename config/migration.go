package config

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	dbdomain "vocab8/domain/db"

	"github.com/go-pg/pg/v10"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	postgresInstance "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewDB(c Config) {
	err := DBMigration(c)
	if err != nil {
		log.Println(err)
		log.Fatal("db migration failed due to the above error")
	}

	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", c.DB.Host, c.DB.Port),
		User:     c.DB.User,
		Password: c.DB.Password,
		Database: c.DB.Name,
	})
	_, err = db.Exec("SELECT 1")
	if err != nil {
		log.Println(err)
		log.Fatal("connection to the database couldnt be established!")
	} else {
		log.Println("connected to the db!")
	}
	DB = db

	err = WordMigration(c)
	if err != nil {
		log.Println(err)
		log.Fatal("word migration failed due to the above error")
	}

}

func DBMigration(c Config) (err error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.DB.User, url.QueryEscape(c.DB.Password), c.DB.Host, c.DB.Port, c.DB.Name))
	if err != nil {
		log.Println(err)
		return
	}

	driver, err := postgresInstance.WithInstance(db, &postgresInstance.Config{
		MigrationsTable: "vocab8_schema_migrations",
	})
	if err != nil {
		log.Println(err)
		return
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+c.DB.MigrationPath,
		"postgres", driver)
	if err != nil {
		log.Println(err)
		return
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Println("no change since the last migartion.")
		err = nil
	}
	if err != nil && err != migrate.ErrNoChange {
		log.Println(err)
		return
	}

	log.Println("migration completed!")

	return
}

func WordMigration(c Config) (err error) {
	var tbm_word_lst []int
	var schema_model dbdomain.WordMigrations
	_ = DB.Model(&schema_model).
		Limit(1).
		Select()

	words, err := os.ReadDir(c.WordPath)
	if err != nil {
		return
	}

	for _, e := range words {
		wordno, err := strconv.Atoi(strings.Split(e.Name(), ".")[0])
		if err != nil {
			return err
		}
		if wordno > schema_model.No {
			tbm_word_lst = append(tbm_word_lst, wordno)
		}
	}

	if len(tbm_word_lst) == 0 {
		return
	}

	sort.Ints(tbm_word_lst)

	var duplicated int
	for _, tbm_word := range tbm_word_lst {
		file, err := os.ReadFile(c.WordPath + "/" + strconv.Itoa(tbm_word) + ".txt")
		if err != nil {
			return err
		}
		parsedFile := ExtractWordsFromFile(file)

		for _, word := range parsedFile {
			_, err = DB.Model(&word).Insert()
			if err != nil {
				pgErr, ok := err.(pg.Error)
				if ok && pgErr.IntegrityViolation() {
					duplicated++
					continue
				} else {
					return err
				}

			}
		}
	}

	fmt.Println("number of duplication:", duplicated)

	schema_model.No = tbm_word_lst[len(tbm_word_lst)-1]
	_, err = DB.Model(&schema_model).
		WherePK().
		Update()
	if err != nil {
		return err
	}

	return
}

func ExtractWordsFromFile(content []byte) (words []dbdomain.Word) {
	lower_index := 0
	var word dbdomain.Word
	content = append(content, 10)

	for el_index, el := range content {
		if el == 10 {
			if el_index == lower_index+1 {
				words = append(words, word)
				word = dbdomain.Word{}
				lower_index = el_index
				continue
			}

			var newline string
			if lower_index == 0 {
				newline = string(content[lower_index:el_index])
			} else {
				newline = string(content[lower_index+1 : el_index])
			}

			switch {
			case word == dbdomain.Word{} && len(newline) > 0:
				word.Word = newline
				word.CreatedAt = time.Now()
				word.UpdatedAt = time.Now()
			default:
				word.Description += fmt.Sprintf("%s\n", newline)
			}
			lower_index = el_index

		}
		if el_index == len(content)-1 {
			words = append(words, word)
			word = dbdomain.Word{}
			lower_index = el_index
			continue
		}

	}
	return
}
