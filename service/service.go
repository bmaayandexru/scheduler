package service

import (
	"database/sql"
	"time"

	"github.com/bmaayandexru/scheduler/nextdate"
	"github.com/bmaayandexru/scheduler/storage"
	"github.com/go-pg/pg/v10/orm"
)

type TaskService struct {
	store storage.TaskStore
}

const limit = 50

var Service TaskService

func NewTaskService(store storage.TaskStore) TaskService {
	return TaskService{store: store}
}

func (ts TaskService) Add(task storage.Task) (orm.Result, error) {
	return ts.store.Add(task)
}

/*
	func (ts TaskService) Delete(id string) (orm.Result, error) {
		return ts.store.Delete(id)
	}
*/
func (ts TaskService) Delete(id int) (orm.Result, error) {
	return ts.store.Delete(id)
}

func (ts TaskService) Find(search string) ([]storage.Task, orm.Result, error) {
	return ts.store.Find(search)
}

func (ts TaskService) Get(id int) (storage.Task, error) {
	return ts.store.Get(id)
}

func (ts TaskService) Update(task storage.Task) (sql.Result, error) {
	return ts.store.Update(task)
}

func (ts TaskService) Done(id int) error {
	// выполнение задачи - это перенос либо удаление
	var task storage.Task
	var err error
	// запросить задаче по id
	if task, err = ts.Get(id); err != nil {
		return err
	}
	// если правила повторения нет удалить
	if len(task.Repeat) == 0 {
		_, err = ts.Delete(id)
		return err
	} else {
		// если есть вызвать nextdate и перенести (update)
		if task.Date, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat); err != nil {
			return err
		}
		if _, err = ts.Update(task); err != nil {
			return err
		}
		return nil
	}
}
