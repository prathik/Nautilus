package service

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Topic struct {
	Title string
	Times int
	NextRun time.Time
}

type SpacedRepetition struct {
	SqlDataBase *sql.DB
}

func (sr SpacedRepetition) Init()  {
	sqlStmt := `
	create table if not exists sr_data 
	(id integer not null primary key, title text, times int,
	next_run datetime(9));
	`

	_, createErr := sr.SqlDataBase.Exec(sqlStmt)

	if createErr != nil {
		panic(createErr)
	}
}

//Add Adds a topic into persistence layer
func (sr SpacedRepetition) Add(topic Topic)  {
	insertStatementTemplate := `
insert into sr_data (title, times, next_run) values ("%s", 0, "%s")
`
	_, err := sr.SqlDataBase.Exec(
		fmt.Sprintf(insertStatementTemplate, topic.Title, time.Now().Add(
			time.Hour * 1).String()))
	if err != nil {
		panic(err)
	}
}

//GetTopicNow Get a topic from persistence layer
func (sr SpacedRepetition) GetTopicNow() (topic Topic) {
	return Topic{
		Title: "Hello",
		Times: 1,
		NextRun: time.Now(),
	}
}

//RescheduleTopic Updates the next run time in the persistence layer
func (sr SpacedRepetition) RescheduleTopic(topic Topic) {
}
