package course

import (
	"log"
	"time"
)

type (
	Service interface {
		Create(name, startDate, endDate string) (*Course, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(l *log.Logger, repo Repository) Service {
	return &service{
		log:  l,
		repo: repo,
	}
}

func (s service) Create(name, startDate, endDate string) (*Course, error) {
	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	course := &Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}

	course, err = s.repo.Create(course)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	return course, nil
}
