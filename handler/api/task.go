package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type TaskAPI interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)

	MarkTask(w http.ResponseWriter, r *http.Request)
	UnMarkTask(w http.ResponseWriter, r *http.Request)
	Reminder(w http.ResponseWriter, r *http.Request)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) MarkTask(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("task_id")
	err := t.taskService.MarkTask(r.Context(), taskId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}
	// http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"task_id": taskId,
		"message": "success mark task",
	})
}

func (t *taskAPI) UnMarkTask(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("task_id")

	err := t.taskService.UnMarkTask(r.Context(), taskId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}
	// http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"task_id": taskId,
		"message": "success unmark task",
	})
}

func (t *taskAPI) Reminder(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest
	id := r.Context().Value("id").(string)
	idLogin, _ := strconv.Atoi(id)

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
	}

	rmdr, err := t.taskService.UpdateTaskReminder(r.Context(), &entity.Task{
		ID:       task.ID,
		Reminder: task.Reminder,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": idLogin,
		"task_id": rmdr.ID,
		"message": "success update task reminder",
	})
}

func (t *taskAPI) GetTask(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	userId := r.Context().Value("id")
	taskId := r.URL.Query().Get("task_id")

	if userId == nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("get task", "invalid user id")
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	} else if taskId == "" {
		idLogin, err := strconv.Atoi(userId.(string))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("get category", err.Error())
			json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
			return
		}

		tasks, err := t.taskService.GetTasks(r.Context(), int(idLogin))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)
	} else {
		tid, _ := strconv.Atoi(taskId)

		tasks, err := t.taskService.GetTaskByID(r.Context(), int(tid))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)
	}
}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest            // membuat variabel task dengan tipe data TaskRequest
	id := r.Context().Value("id").(string) // mengambil data id dari kuki
	idLogin, _ := strconv.Atoi(id)         // mengubah string menjadi int

	err := json.NewDecoder(r.Body).Decode(&task) // mengubah data json menjadi struct
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	// TODO: answer here
	if task.Title == "" || task.Description == "" || task.CategoryID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
	}

	tsk, err := t.taskService.StoreTask(r.Context(), &entity.Task{ // memasukkan data ke database
		// ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		// Completed:   task.Completed,
		CategoryID: task.CategoryID,
		UserID:     idLogin,
	})
	if err != nil { // jika terjadi error pada penginputan data
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": idLogin,
		"task_id": tsk.ID,
		"message": "success create new task",
	})
}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	userId := r.Context().Value("id").(string)
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("delete task", "invalid user id")
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	taskID := r.URL.Query().Get("task_id")
	tid, _ := strconv.Atoi(taskID)

	err := t.taskService.DeleteTask(r.Context(), tid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}

	uid, _ := strconv.Atoi(userId)
	result := map[string]interface{}{
		"user_id": uid,
		"task_id": tid,
		"message": "success delete task",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest
	id := r.Context().Value("id").(string)
	idLogin, _ := strconv.Atoi(id)

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	// TODO: answer here

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
	}

	tsk, err := t.taskService.UpdateTask(r.Context(), &entity.Task{
		ID:          task.ID,
		Title:       task.Title,
		Reminder:    task.Reminder,
		Description: task.Description,
		CategoryID:  task.CategoryID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":     idLogin,
		"task_id":     tsk.ID,
		"title":       tsk.Title,
		"description": tsk.Description,
		"category_id": tsk.CategoryID,
		"message":     "success update task",
	})
}

func (t *taskAPI) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = entity.Task{
		ID:         task.ID,
		CategoryID: task.CategoryID,
		UserID:     int(idLogin),
	}

	_, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userId,
		"task_id": task.ID,
		"message": "success update task category",
	})
}
