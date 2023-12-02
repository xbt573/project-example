package services

import (
	"errors"
	"github.com/xbt573/project-example/models"
	"gorm.io/gorm"
)

var (
	ErrTodoServiceNotFound = errors.New("todo not found")
)

type TodoService interface {
	Create(models.TODO) (models.TODO, error)
	List(userid uint) ([]models.TODO, error)
	Find(userid, id uint) (models.TODO, error)
	Update(todo models.TODO) (models.TODO, error)
	Delete(userid, id uint) (models.TODO, error)
}

type concreteTodoService struct {
	database *gorm.DB
}

func NewTodoService(database *gorm.DB) TodoService {
	return &concreteTodoService{database}
}

func (c *concreteTodoService) Create(todo models.TODO) (models.TODO, error) {
	result := c.database.Create(&todo)
	return todo, result.Error
}

func (c *concreteTodoService) List(userid uint) ([]models.TODO, error) {
	var todos []models.TODO

	result := c.database.Where("user_id = ?", userid).Find(&todos)
	return todos, result.Error
}

func (c *concreteTodoService) Find(userid, id uint) (models.TODO, error) {
	var todo models.TODO
	result := c.database.Where("user_id = ? AND id = ?", userid, id).First(&todo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.TODO{}, ErrTodoServiceNotFound
		}

		return models.TODO{}, result.Error
	}

	return todo, nil
}

func (c *concreteTodoService) Update(todo models.TODO) (models.TODO, error) {
	var oldTodo models.TODO
	result := c.database.Where("user_id = ? AND id = ?", todo.UserID, todo.ID).First(&oldTodo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.TODO{}, ErrTodoServiceNotFound
		}

		return models.TODO{}, result.Error
	}

	result = c.database.Save(&todo)
	return todo, result.Error
}

func (c *concreteTodoService) Delete(userid, id uint) (models.TODO, error) {
	var todo models.TODO
	result := c.database.Where("user_id = ? AND id = ?", userid, id).First(&todo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.TODO{}, ErrTodoServiceNotFound
		}

		return models.TODO{}, result.Error
	}

	result = c.database.Delete(&todo)
	return todo, result.Error
}
