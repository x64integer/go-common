package api

import (
	"bytes"
	"fmt"
	"strings"
)

// Service is layer between route handler and database access
type Service struct {
}

// Register user account
func (svc *Service) Register(fields []*entityField) ([]byte, error) {
	var query string
	var columns []string
	var binders []string
	var data []interface{}
	queryBuff := bytes.NewBufferString("")

	// construct query and data: "INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)"
	for i, field := range fields {
		if field.AuthType == "auto_id" || field.AuthType == "auto_gen" {
			continue
		}

		columns = append(columns, field.AuthKey)
		binders = append(binders, "$"+fmt.Sprint(i))

		data = append(data, field.AuthValue)
	}

	queryBuff.WriteString("INSERT INTO " + strings.ToLower(fields[0].AuthTable) + "s (" + strings.Join(columns, ", ") + ") VALUES (" + strings.Join(binders, ", ") + ")")

	query = queryBuff.String()

	dao := &Dao{}

	if err := dao.Save(query, data); err != nil {
		return nil, err
	}

	return nil, nil
}

// Login user
func (svc *Service) Login(fields []*entityField) ([]byte, error) {
	// for _, field := range fields {
	// 	log.Println(field.Key, field.Value)
	// }
	return nil, nil
}

// Logout user
func (svc *Service) Logout(fields []*entityField) ([]byte, error) {
	// for _, field := range fields {
	// 	log.Println(field.Key, field.Value)
	// }
	return nil, nil
}
