package api

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/x64integer/go-common/util"

	"github.com/x64integer/go-common/password"
	"github.com/x64integer/go-common/storage"
)

// Service is responsible to store data, for now
type Service struct {
	*storage.Container
}

// Register user account
func (svc *Service) Register(fields []*entityField) ([]byte, error) {
	var columns []string
	var queryParams []string
	var data []interface{}
	queryBuff := bytes.NewBufferString("")

	// construct query and data: "INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)"
	for i, field := range fields {

		// TODO: this is only prototype of what has to be done, will be changed
		var dataValue interface{}
		var param string

		switch field.AuthType {
		case "secret":
			pwd, err := password.Hash(fmt.Sprint(field.AuthValue))
			if err != nil {
				return nil, err
			}

			dataValue = pwd
		case "auto_gen":
			continue
		default:
			dataValue = field.AuthValue
		}

		switch util.Env("DB_DRIVER", "") {
		case strings.ToLower("mysql"):
			param = "?"
		default:
			param = "$" + fmt.Sprint(i)
		}

		columns = append(columns, field.AuthKey)
		queryParams = append(queryParams, param)
		data = append(data, dataValue)
	}

	queryBuff.WriteString(
		"INSERT INTO " + strings.ToLower(fields[0].AuthTable) + "s (" + strings.Join(columns, ", ") + ") VALUES (" + strings.Join(queryParams, ", ") + ")",
	)

	if _, err := svc.SQL.Exec(queryBuff.String(), data...); err != nil {
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
