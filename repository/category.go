package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error)
	StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error)
	StoreManyCategory(ctx context.Context, categories []entity.Category) error
	GetCategoryByID(ctx context.Context, id int) (entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error) {
	// return nil, nil // TODO: replace this
	var categories []entity.Category
	err := r.db.WithContext(ctx).Where("user_id = ?", id).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil

}

func (r *categoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error) {
	// return 0, nil // TODO: replace this
	err = r.db.WithContext(ctx).Create(category).Error
	if err != nil {
		return 0, err
	} else {
		return category.ID, nil
	}
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []entity.Category) error {
	// return nil // TODO: replace this
	err := r.db.WithContext(ctx).Create(&categories).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.Category, error) {
	// return entity.Category{}, nil // TODO: replace this
	var category entity.Category
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error
	if err != nil {
		return entity.Category{}, err
	}
	return category, nil

}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Table("categories").Where("id = ?", category.ID).Updates(&category).Error
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Table("categories").Where("id = ?", id).Delete(&entity.Category{}).Error
}
