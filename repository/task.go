package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, id int) ([]entity.Task, error)
	StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTasks(ctx context.Context, id int) ([]entity.Task, error) {
	// return nil, nil // TODO: replace this
	var tasks []entity.Task
	err := r.db.WithContext(ctx).Where("user_id = ?", id).Find(&tasks).Error
	if err != nil {
		return []entity.Task{}, err
	}
	return tasks, nil
}

func (r *taskRepository) StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error) {
	// return 0, nil // TODO: replace this
	err = r.db.WithContext(ctx).Create(&task).Error
	if err != nil {
		return 0, err
	} else {
		return task.ID, nil
	}
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	// return entity.Task{}, nil // TODO: replace this
	var task entity.Task
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&task).Error
	if err != nil {
		return entity.Task{}, err
	}
	return task, nil
}

func (r *taskRepository) GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error) {
	// return nil, nil // TODO: replace this
	var tasks []entity.Task
	err := r.db.WithContext(ctx).Where("category_id = ?", catId).Find(&tasks).Error
	if err != nil {
		return []entity.Task{}, err
	}
	return tasks, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task) error {
	// return nil // TODO: replace this
	err := r.db.WithContext(ctx).Model(&entity.Task{}).Where("id = ?", task.ID).Updates(&task).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	// return nil // TODO: replace this
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Task{}).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}
