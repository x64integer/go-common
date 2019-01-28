package api

import "log"

// Service is layer between route handler and database access
type Service struct {
}

// Register user account
func (svc *Service) Register(fields []*entityField) ([]byte, error) {
	for _, field := range fields {
		log.Println(field.Key, field.Value)
	}
	return nil, nil
}

// Login user
func (svc *Service) Login(fields []*entityField) ([]byte, error) {
	for _, field := range fields {
		log.Println(field.Key, field.Value)
	}
	return nil, nil
}

// Logout user
func (svc *Service) Logout(fields []*entityField) ([]byte, error) {
	for _, field := range fields {
		log.Println(field.Key, field.Value)
	}
	return nil, nil
}
