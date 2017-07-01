package services

import "github.com/social-tournament-service/app/daos"

type DefaultService struct{
	defaultDao *daos.DefaultDao
}

func NewDefaultService(defaultDao *daos.DefaultDao) *DefaultService{
	return &DefaultService{
		defaultDao: defaultDao,
	}
}

func (s *DefaultService) Reset() error{
	return s.defaultDao.Reset()
}