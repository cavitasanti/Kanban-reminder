package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type CategoryAPI interface {
	GetCategory(w http.ResponseWriter, r *http.Request)
	CreateNewCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryWithTasks(w http.ResponseWriter, r *http.Request)
}

type categoryAPI struct {
	categoryService service.CategoryService
}

func NewCategoryAPI(categoryService service.CategoryService) *categoryAPI {
	return &categoryAPI{categoryService}
}

func (c *categoryAPI) GetCategory(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	userId := r.Context().Value("id")

	if userId == nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("get category", "invalid user id")
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	} else {
		idLogin, err := strconv.Atoi(userId.(string))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("get category", err.Error())
			json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
			return
		}

		categories, err := c.categoryService.GetCategories(r.Context(), int(idLogin))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(categories)
	}
}

func (c *categoryAPI) CreateNewCategory(w http.ResponseWriter, r *http.Request) {

	var category entity.CategoryRequest
	id := r.Context().Value("id").(string)
	idLogin, _ := strconv.Atoi(id)
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}
	if category.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
	}
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
	}
	categoryy, err := c.categoryService.StoreCategory(r.Context(), &entity.Category{
		Type:   category.Type,
		UserID: idLogin,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
	}

	result := map[string]interface{}{
		"user_id":     idLogin,
		"category_id": categoryy.ID,
		"message":     "success create new category",
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(result)
}

func (c *categoryAPI) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	userId := r.Context().Value("id").(string)
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("delete category", "invalid user id")
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categoryID := r.URL.Query().Get("category_id")
	cid, _ := strconv.Atoi(categoryID)

	err := c.categoryService.DeleteCategory(r.Context(), cid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}

	uid, _ := strconv.Atoi(userId)
	result := map[string]interface{}{
		"user_id":     uid,
		"category_id": cid,
		"message":     "success delete category",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)

}

func (c *categoryAPI) GetCategoryWithTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("get category task", err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categories, err := c.categoryService.GetCategoriesWithTasks(r.Context(), int(idLogin))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)

}
