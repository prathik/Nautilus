package service

import (
	"database/sql"
	"fmt"
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

	sr.Add(&Topic{
		Title: "Test",
	})

	rows, err := db.Query(`select id, title, times, next_run
from sr_data`)

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

		nr, _ := time.Parse(time.RFC3339,
			nextRun)

		if nr.Before(time.Now().Add(time.Minute * 59)) == true {
			t.Errorf("%s is before 1 hour. Base time %s", nr.String(),
				nextRun)
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
	sr.Add(&Topic{
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

func TestSpacedRepetition_GivenTopicInDb_WhenGetCalled_ReturnTopTopic(t *testing.T)  {
	db := getDb(t)
	defer db.Close()
	defer func(){
		_ = os.Remove("./sr.db")
	}()

	sr := SpacedRepetition{
		SqlDataBase:db,
	}

	// Create table
	sr.Init()

	insertStatementTemplate := `
insert into sr_data (title, times, next_run) values ("%s", 0, "%s")
`
	_, _ = db.Exec(fmt.Sprintf(insertStatementTemplate, "Test", time.Now().Format(time.RFC3339)))
	_, _ = db.Exec(fmt.Sprintf(insertStatementTemplate, "Test Old", time.Now().Add(
		time.Hour * time.Duration(-1)).Format(time.RFC3339)))

	topic := sr.GetTopicNow()

	if topic.Title != "Test Old" {
		t.Error("Invalid topic: " + topic.Title)
	}
}

func TestSpacedRepetition_GivenNoTopicInDb_WhenGetCalled_ReturnTopTopic(t *testing.T)  {
	db := getDb(t)
	defer db.Close()
	defer func(){
		_ = os.Remove("./sr.db")
	}()

	sr := SpacedRepetition{
		SqlDataBase:db,
	}

	// Create table
	sr.Init()

	topic := sr.GetTopicNow()

	if topic != nil {
		t.Error("Topic is present.")
	}
}

func TestSpacedRepetition_GivenTopicId_WhenRescheduleCalled_UpdateNextRunTime(t *testing.T)  {
	db := getDb(t)
	defer db.Close()
	defer func(){
		_ = os.Remove("./sr.db")
	}()

	sr := SpacedRepetition{
		SqlDataBase:db,
	}

	// Create table
	sr.Init()

	insertStatementTemplate := `
insert into sr_data (title, times, next_run) values ("%s", 0, "%s")
`
	_, _ = db.Exec(fmt.Sprintf(insertStatementTemplate, "Test", time.Now().Format(
		time.RFC3339)))

	topic := sr.GetTopicNow()

	if topic == nil {
		t.Errorf("No topic fetched.")
	}

	if topic.Title != "Test" {
		t.Error("Invalid topic.")
	}

	sr.RescheduleTopic(topic)

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

		nr, _ := time.Parse(time.RFC3339,
			nextRun)

		// Scheduled 3 days ahead from now.
		if nr.Before(time.Now().Add(time.Hour * 24 * 2)) == true {
			t.Errorf("Next run is before two days %s", nr.String())
		}

		if nr.After(time.Now().Add(time.Hour * 24 * 4)) == true {
			t.Fail()
		}

		count++
	}

	if count == 0 {
		t.Error("No rows")
	}
}

func TestSpacedRepetition_GivenTopics_WhenGetAllCalled_ThenReturnAllTopicSlice(t *testing.T) {
	db := getDb(t)
	defer db.Close()
	defer func(){
		_ = os.Remove("./sr.db")
	}()

	sr := SpacedRepetition{
		SqlDataBase:db,
	}

	// Create table
	sr.Init()

	sr.Add(&Topic{Title: "Test"})
	sr.Add(&Topic{Title: "Test One"})
	sr.Add(&Topic{Title: "Test Two"})

	topics := sr.GetAll()

	if len(topics) != 3 {
		t.Errorf("All topics have not been fetched.")
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

