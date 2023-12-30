package user

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(user *User) (*User, error)
		FindAll(filters Filters, offset, limit int) ([]User, error)
		FindOne(id string) (*User, error)
		Delete(id string) (*User, error)
		Update(id string, firstname, lastname, email, phone *string) (*User, error)
		Count(filters Filters) (int, error)
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

func (repo *repo) Create(user *User) (*User, error) {
	query := repo.db.Create(user)
	if query.Error != nil {
		repo.log.Println(query.Error)
		return user, query.Error
	}
	repo.log.Println("User created successfully with id: ", user.ID)
	return user, nil
}

func (repo *repo) FindAll(filters Filters, offset, limit int) ([]User, error) {
	var u []User
	tx := repo.db.Model(&u)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at DESC").Find(&u)
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

func (repo *repo) Count(filters Filters) (int, error) {
	var count int64
	tx := repo.db.Model(User{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

/* Mover a un Helper */
func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.Firstname != "" {
		filters.Firstname = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Firstname))
		tx = tx.Where("lower(firstname) like ?", filters.Firstname)
	}
	if filters.Lastname != "" {
		filters.Lastname = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Lastname))
		tx = tx.Where("lower(lastname) like ?", filters.Lastname)
	}
	return tx
}
