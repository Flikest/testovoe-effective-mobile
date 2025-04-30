package storage

import (
	"net/http"

	"github.com/Flikest/testovoe-effective-mobile/internal/entity"
)

type User struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

func (s Storage) GetUsers(u User, ch chan []entity.User) {
	getUserQuery := `
		SELECT * FROM users 
		WHERE ($1::text IS NULL OR name LIKE $1) 
		AND ($2::text IS NULL OR surname LIKE $2) 
		AND ($3::text IS NULL OR patronymic LIKE $3)
	`

	name := "%" + u.Name + "%"
	surname := "%" + u.Surname + "%"
	patronymic := "%" + u.Patronymic + "%"

	s.Log.Info("Executing query with parameters: ", name, surname, patronymic)

	rows, err := s.DB.Query(s.Context, getUserQuery, name, surname, patronymic)
	if err != nil {
		s.Log.Error("SQL error, failed to find the users: ", err)
		ch <- nil
		return
	}
	defer rows.Close()

	var users []entity.User

	for rows.Next() {
		var row entity.User
		if err := rows.Scan(&row.ID, &row.Name, &row.Surname, &row.Patronymic); err != nil {
			s.Log.Error("Failed to scan user: ", err)
			ch <- nil
			return
		}
		users = append(users, row)
	}

	if err := rows.Err(); err != nil {
		s.Log.Error("Error occurred during row iteration: ", err)
		ch <- nil
		return
	}

	ch <- users
}

func (s Storage) DeleteUser(ID int, ch chan int) {
	deleteUserQuery := "DELETE FROM users WHERE id=$1"

	_, err := s.DB.Exec(s.Context, deleteUserQuery, ID)
	if err != nil {
		s.Log.Info("Failed to remove the user: ", err)
		ch <- http.StatusInternalServerError
		return
	}

	ch <- http.StatusOK
}

func (s Storage) PatchUser(u entity.User, ch chan int) {
	updateUserQuery := "UPDATE users SET name=$2, surname=$3, patronymic=$4 WHERE id=$1"

	_, err := s.DB.Exec(s.Context, updateUserQuery, u.ID, u.Name, u.Surname, u.Patronymic)
	if err != nil {
		s.Log.Info("Failed to update the user data: ", err)
		ch <- http.StatusInternalServerError
		return
	}

	ch <- http.StatusOK
}

func (s Storage) AddUser(u User, ch chan int) {
	createUserQuery := "INSERT INTO users(name, surname, patronymic) VALUES ($1, $2, $3)"

	_, err := s.DB.Exec(s.Context, createUserQuery, u.Name, u.Surname, u.Patronymic)
	if err != nil {
		s.Log.Info("failed to add the user: ", err)
		ch <- http.StatusInternalServerError
		return
	}

	ch <- http.StatusCreated
}
