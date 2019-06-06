package service

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
	"time"
)

func TestSpacedRepetition_Add(t *testing.T) {
	db := getDb(t)
	defer db.Close()
	defer func() {
		_ = os.Remove("./sr.db")
	}()


	sr := SpacedRepetition{
		SqlDataBase: db,
	}

	sr.Init()

	sr.Add(Topic{
		Title: "Test",
	})

	rows, err := db.Query("select id, title, times, next_run from sr_data")

	if err != nil {
		t.Error(err)
	}

	defer rows.Close()

	var count = 0

	for rows.Next() {
		var id int
		var title string
		var times int
		var nextRun string
		err := rows.Scan(&id, &title, &times, &nextRun)
		if err != nil {
			t.Error(err)
		}
		if title != "Test" {
			t.Fail()
		}

		if id != 1 {
			t.Fail()
		}

		if times != 0 {
			t.Fail()
		}

		nr, _ := time.Parse("2019-06-06 21:42:16.790853 +0530 IST m=+3600.003102276",
			nextRun)

		if nr.Before(time.Now().Add(time.Hour * 1)) == false {
			t.Fail()
		}

		if nr.After(time.Now().Add(time.Minute * 61)) == true {
			t.Fail()
		}

		count++
	}

	if count == 0 {
		t.Error("No rows")
	}
}

func TestSpacedRepetition_DoubleInit(t *testing.T) {
	db := getDb(t)
	defer db.Close()
	defer func() {
		_ = os.Remove("./sr.db")
	}()


	sr := SpacedRepetition{
		SqlDataBase: db,
	}

	sr.Init()
	sr.Add(Topic{
		Title: "Test",
	})
	sr.Init()
	rows, _ := db.Query("select count(*) from sr_data")
	for rows.Next() {
		var count int
		_ = rows.Scan(&count)
		if count != 1 {
			t.Fail()
		}
	}
}

func getDb(t *testing.T) *sql.DB {
	_ = os.Remove("./sr.db")

	db, dbErr := sql.Open("sqlite3", "./sr.db")
	if dbErr != nil {
		t.Fail()
	}
	return db
}

