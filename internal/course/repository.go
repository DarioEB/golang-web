package course

import (
	"log"

	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(course *Course) (*Course, error)
	}

	repo struct {
		log *log.Logger
		db  *gorm.DB
	}
)

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(course *Course) (*Course, error) {
	query := repo.db.Create(course)
	if query.Error != nil {
		repo.log.Println(query.Error)
		return course, query.Error
	}
	repo.log.Println("Course created successfully with id: ", course.ID)
	return course, nil
}
