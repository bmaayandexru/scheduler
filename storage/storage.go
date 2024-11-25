package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	_ "modernc.org/sqlite"
)

/*
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"` // omitempty
	Title   string `json:"title"`
	Comment string `json:"comment"` // omitempty
	Repeat  string `json:"repeat"`  // omitempty
}
*/

type Task struct {
	ID      int    `pg:"id,pk"   json:"id"`
	Date    string `pg:"date"    json:"date"` // omitempty
	Title   string `pg:"title"   json:"title"`
	Comment string `pg:"comment" json:"comment"` // omitempty
	Repeat  string `pg:"repeat"  json:"repeat"`  // omitempty
}

const (
	templ = "20060102"
	limit = 50
)

type TaskStore struct {
	DB *pg.DB
}

func NewTaskStore(db *pg.DB) TaskStore {
	return TaskStore{DB: db}
}

func (ts TaskStore) Add(task Task) (orm.Result, error) {
	return ts.DB.Exec("INSERT INTO scheduler(date, title, comment, repeat) VALUES ( ?, ?, ?, ?) ",
		task.Date, task.Title, task.Comment, task.Repeat)
}

/*
	func (ts TaskStore) Delete(id string) (orm.Result, error) {
		//	return ts.DB.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
		return nil, nil
	}
*/
func (ts TaskStore) Delete(id int) (orm.Result, error) {
	//	return ts.DB.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	return nil, nil
}

func (ts TaskStore) Find(search string) ([]Task, orm.Result, error) {
	// возвращаем всё если строка пустая

	if len(search) == 0 {

		res, err := ts.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit",
			sql.Named("limit", 50))
		return []Task{}, res, err
		/*
			var tasks []Task
			_, err := ts.DB.Query(&tasks, "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", 50)
			if err != nil {
				return nil, err
			}
			return nil, nil\
		*/
	}
	/*
		// парсим строку на дату
		if date, err := time.Parse("02-01-2006", search); err == nil {
			// дата есть
			return ts.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date LIMIT :limit",
				sql.Named("date", date.Format(templ)),
				sql.Named("limit", limit))
		} else {
			// даты нет
			search = "%" + search + "%"
			return ts.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE UPPER(title) LIKE UPPER(:search) OR UPPER(comment) LIKE UPPER(:search) ORDER BY date LIMIT :limit",
				sql.Named("search", search),
				sql.Named("limit", limit))
		}
	*/
	return []Task{}, nil, nil
}

func (ts TaskStore) Get(id int) (Task, error) {
	/*
		row := ts.DB.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
		task := Task{}
		err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		return task, err
	*/
	task := Task{}
	err := ts.DB.Model(&task).Where("id = ?", id).Select()
	return task, err
	//return Task{}, nil
}

func (ts TaskStore) Update(task Task) (sql.Result, error) {
	/*
		return ts.DB.Exec("UPDATE scheduler SET  date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
			sql.Named("id", task.ID),
			sql.Named("date", task.Date),
			sql.Named("title", task.Title),
			sql.Named("comment", task.Comment),
			sql.Named("repeat", task.Repeat))
	*/
	return nil, nil
}

/*
var schemaSQL string = `
CREATE TABLE IF NOT EXIST scheduler (

	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "20000101",
	title VARCHAR(32) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""

);
CREATE INDEX idx_date ON scheduler (date);
CREATE INDEX idx_title ON scheduler (title);
`
*/
var schemaSQL string = `
CREATE TABLE IF NOT EXISTS scheduler (
    id SERIAL PRIMARY KEY,
    date CHAR(8) NOT NULL, 
    title VARCHAR(32) NOT NULL,
    comment TEXT NOT NULL,
    repeat VARCHAR(128) NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
CREATE INDEX IF NOT EXISTS idx_title ON scheduler (title);
`
var DBFileRun = "scheduler.db"

func InitDBase() (*pg.DB, error) {
	var pgDB *pg.DB
	fmt.Println("Init Data Base...")
	pgDB = pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "postgres",
		Password: "password",
		Database: "testdb",
	})

	_, err := pgDB.Exec(schemaSQL)

	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v", err)
		return nil, err
	}
	fmt.Println("База открыта")
	return pgDB, err
}

func createSchema(db *pg.DB, model ...interface{}) error {

	return db.Model(model...).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})

}
