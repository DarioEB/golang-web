package user

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) (*User, error)
	FindAll() ([]User, error)
	FindOne(id string) (*User, error)
	Delete(id string) (*User, error)
	Update(id string, firstname, lastname, email, phone *string) (*User, error)
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(user *User) (*User, error) {
	user.ID = uuid.New().String()
	query := repo.db.Create(user)
	if query.Error != nil {
		repo.log.Println(query.Error)
		return user, query.Error
	}
	repo.log.Println("User created successfully with id: ", user.ID)
	return user, nil
}

func (repo *repo) FindAll() ([]User, error) {
	var u []User
	result := repo.db.Model(&u).Order("created_at DESC").Find(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return u, nil
}

func (repo *repo) FindOne(id string) (*User, error) {
	user := User{ID: id}
	result := repo.db.First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *repo) Delete(id string) (*User, error) {
	user := User{ID: id}
	result := repo.db.Delete(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *repo) Update(id string, firstname, lastname, email, phone *string) (*User, error) {
	user := User{}
	values := make(map[string]interface{})
	if firstname != nil {
		values["firstname"] = *firstname
	}
	if lastname != nil {
		values["lastname"] = *lastname
	}
	if email != nil {
		values["email"] = *email
	}
	if phone != nil {
		values["phone"] = *phone
	}

	result := repo.db.Model(&user).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
