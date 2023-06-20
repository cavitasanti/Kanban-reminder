package service

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/repository"
	"context"
	"strconv"
)

type TaskService interface {
	GetTasks(ctx context.Context, id int) ([]entity.Task, error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	StoreTask(ctx context.Context, task *entity.Task) (entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) (entity.Task, error)
	UpdateTaskReminder(ctx context.Context, task *entity.Task) (entity.Task, error)
	DeleteTask(ctx context.Context, id int) error

	MarkTask(ctx context.Context, id string) error
	UnMarkTask(ctx context.Context, id string) error
}

type taskService struct {
	taskRepo     repository.TaskRepository
	categoryRepo repository.CategoryRepository
}

func NewTaskService(taskRepo repository.TaskRepository, categoryRepo repository.CategoryRepository) TaskService {
	return &taskService{taskRepo, categoryRepo}
}

func (s *taskService) MarkTask(ctx context.Context, id string) error {
	idn, _ := strconv.Atoi(id)
	// return s.taskRepo.UpdateMarkTask(ctx, map[string]interface{}{"completed": true, "ID": idn})
	err := s.taskRepo.UpdateMarkTask(ctx, map[string]interface{}{"completed": true, "ID": idn})
	if err != nil {
		return err
	}
	return nil
}

func (s *taskService) UnMarkTask(ctx context.Context, id string) error {
	idn, _ := strconv.Atoi(id)
	// return s.taskRepo.UpdateMarkTask(ctx, map[string]interface{}{"completed": false, "ID": idn})
	err := s.taskRepo.UpdateMarkTask(ctx, map[string]interface{}{"completed": false, "ID": idn})
	if err != nil {
		return err
	}
	return nil
}

func (s *taskService) GetTasks(ctx context.Context, id int) ([]entity.Task, error) {
	return s.taskRepo.GetTasks(ctx, id)
}

func (s *taskService) StoreTask(ctx context.Context, task *entity.Task) (entity.Task, error) {
	_, err := s.taskRepo.StoreTask(ctx, task) // memanggil fungsi StoreTask di repository
	if err != nil {
		return entity.Task{}, err
	}
	return *task, nil
}

func (s *taskService) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	return s.taskRepo.GetTaskByID(ctx, id)
}

func (s *taskService) UpdateTask(ctx context.Context, task *entity.Task) (entity.Task, error) {
	if task.CategoryID != 0 {
		cat, err := s.categoryRepo.GetCategoryByID(ctx, task.CategoryID)
		if err != nil {
			return entity.Task{}, err
		}

		if cat.ID == 0 || cat.Type == "" || cat.UserID != task.UserID {
			return entity.Task{}, err
		}
	}

	err := s.taskRepo.UpdateTask(ctx, task)
	if err != nil {
		return entity.Task{}, err
	}
	return *task, nil
}

func (s *taskService) UpdateTaskReminder(ctx context.Context, task *entity.Task) (entity.Task, error) {
	if task.CategoryID != 0 {
		cat, err := s.categoryRepo.GetCategoryByID(ctx, task.CategoryID)
		if err != nil {
			return entity.Task{}, err
		}

		if cat.ID == 0 || cat.Type == "" || cat.UserID != task.UserID {
			return entity.Task{}, err
		}
	}

	err := s.taskRepo.UpdateTaskReminder(ctx, task)
	if err != nil {
		return entity.Task{}, err
	}
	return *task, nil
}

func (s *taskService) DeleteTask(ctx context.Context, id int) error {
	return s.taskRepo.DeleteTask(ctx, id)
}
