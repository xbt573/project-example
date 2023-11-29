package services

import (
	"github.com/xbt573/project-example/models"
	"gorm.io/gorm"
	"sync"
)

var (
	once    sync.Once
	service *TodoService
)

type TodoService struct {
	database *gorm.DB
}

func InitTodoService(db *gorm.DB) {
	once.Do(func() {
		service = &TodoService{database: db}
	})
}

func GetTodoServiceInstance() *TodoService {
	return service
}

func (t *TodoService) Create(todo models.TODO) (models.TODO, error) {
	result := t.database.Create(&todo)
	if result.Error != nil {
		return models.TODO{}, result.Error
	}

	return todo, nil
}

func (t *TodoService) List() ([]models.TODO, error) {
	var users []models.TODO

	result := t.database.Find(&users)
	if result.Error != nil {
		return []models.TODO{}, result.Error
	}

	return users, nil
}

func (t *TodoService) Search(id uint) (models.TODO, error) {
	todo := models.TODO{ID: id}

	result := t.database.Find(&todo)
	if result.Error != nil {
		return models.TODO{}, result.Error
	}

	return todo, nil
}

func (t *TodoService) Update(todo models.TODO) (models.TODO, error) {
	result := t.database.Save(&todo)
	if result.Error != nil {
		return models.TODO{}, result.Error
	}

	return todo, nil
}

func (t *TodoService) Delete(id uint) (models.TODO, error) {
	todo, err := t.Search(id)
	if err != nil {
		return models.TODO{}, err
	}

	result := t.database.Delete(&todo)
	if result.Error != nil {
		return models.TODO{}, result.Error
	}

	return todo, nil
}
