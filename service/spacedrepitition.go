package service

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Topic struct {
	Title   string
	Times   int
	NextRun time.Time
	Id      int64
}

type SpacedRepetition struct {
	SqlDataBase *sql.DB
}

func (sr SpacedRepetition) Init()  {
	sqlStmt := `
	create table if not exists sr_data 
	(id integer not null primary key, title text, times int,
	next_run datetime);
	`

	_, createErr := sr.SqlDataBase.Exec(sqlStmt)

	if createErr != nil {
		panic(createErr)
	}
}

//Add Adds a topic into persistence layer
func (sr SpacedRepetition) Add(topic *Topic)  {
	insertStatementTemplate := `
insert into sr_data (title, times, next_run) values ("%s", 0, "%s")
`
	_, err := sr.SqlDataBase.Exec(
		fmt.Sprintf(insertStatementTemplate, topic.Title, time.Now().Add(
			time.Hour * 1).Format(time.RFC3339)))
	if err != nil {
		panic(err)
	}
}

//GetTopicNow Get a topic from persistence layer
func (sr SpacedRepetition) GetTopicNow() (topic *Topic) {
	getOldestTemplate := `SELECT id, title, times, next_run FROM sr_data
WHERE next_run <= "%s" ORDER BY next_run ASC LIMIT 1
`
	rows, err := sr.SqlDataBase.Query(fmt.Sprintf(getOldestTemplate,
		time.Now().Format(time.RFC3339)))

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		var title string
		var times int
		var nextRun string

		err := rows.Scan(&id, &title, &times, &nextRun)
		if err != nil {
			panic(err)
		}

		nextRunTime, _ := time.Parse(time.RFC3339, nextRun)
		topic = &Topic{
			Id: id,
			Title: title,
			Times: times,
			NextRun: nextRunTime,
		}
	}

	return topic
}

//RescheduleTopic Updates the next run time in the persistence layer
func (sr SpacedRepetition) RescheduleTopic(topic *Topic) {
	timeFactor := topic.Times + 1

	nextRunTime := topic.NextRun.Add(
		time.Hour * time.Duration(24 * 3 * timeFactor))

	updateStatement := `UPDATE sr_data set next_run = "%s" WHERE id = %d`

	_, _ = sr.SqlDataBase.Exec(fmt.Sprintf(updateStatement, nextRunTime.Format(
		time.RFC3339), topic.Id))
}
