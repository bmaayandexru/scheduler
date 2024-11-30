package storage

// DB - PostgreSQL
// driver - pq
import (
	"database/sql"
	"fmt"
	"log"
	"time"

	//_ "modernc.org/sqlite"
	_ "github.com/lib/pq" // Импорт драйвера
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"` // omitempty
	Title   string `json:"title"`
	Comment string `json:"comment"` // omitempty
	Repeat  string `json:"repeat"`  // omitempty
}

const (
	templ = "20060102"
	limit = 50
)

type TaskStore struct {
	DB *sql.DB
}

func NewTaskStore(db *sql.DB) TaskStore {
	return TaskStore{DB: db}
}

func (ts TaskStore) Add(task Task) error {
	_, err := ts.DB.Exec("INSERT INTO scheduler(date, title, comment, repeat) VALUES ($1, $2, $3, $4) ",
		task.Date, task.Title, task.Comment, task.Repeat)
	return err
}

func (ts TaskStore) Delete(id string) error {
	_, err := ts.DB.Exec("DELETE FROM scheduler WHERE id = $1", id)
	return err
}

func (ts TaskStore) Find(search string) ([]Task, error) {
	var (
		err  error
		rows *sql.Rows
	)
	tasks := make([]Task, 0)
	// возвращаем всё если строка пустая
	if len(search) == 0 {
		rows, err = ts.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT $1", limit)
	} else {
		// парсим строку на дату
		if date, err := time.Parse("02-01-2006", search); err == nil {
			// дата есть
			rows, err = ts.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = $1 LIMIT $2",
				date.Format(templ),
				limit)
		} else {
			// даты нет
			search = "%" + search + "%"
			rows, err = ts.DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE UPPER(title) LIKE UPPER($1) OR UPPER(comment) LIKE UPPER($1) ORDER BY date LIMIT $2",
				search,
				limit)
		}
	}
	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return tasks, err
	}
	return tasks, err
}

func (ts TaskStore) Get(id string) (Task, error) {
	row := ts.DB.QueryRow("SELECT * FROM scheduler WHERE id = $1", id)
	task := Task{}
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	return task, err
}

func (ts TaskStore) Update(task Task) error {
	_, err := ts.DB.Exec("UPDATE scheduler SET  date = $2, title = $3, comment = $4, repeat = $5 WHERE id = $1",
		task.ID,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat)
	return err
}

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

func InitDBase(connStr string) (*sql.DB, error) {
	var SqlDB *sql.DB
	fmt.Println("Init Data Base...")
	SqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}

	// Проверяем соединение
	err = SqlDB.Ping()
	if err != nil {
		log.Fatalf("Ошибка проверки соединения: %v", err)
	}
	if _, err = SqlDB.Exec(schemaSQL); err != nil {
		log.Fatalf("InitDB err: %v", err)
		return SqlDB, err
	}

	fmt.Println("Успешное подключение к базе данных")
	return SqlDB, err
}
