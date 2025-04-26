package storage

import (
	"fmt"
	"log/slog"
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
	SELECT *
	FROM users
	WHERE (name = $1 OR $1 IS NULL)
  	AND (surname = $2 OR $2 IS NULL)
  	AND (patronymic = $3 OR $3 IS NULL);
	`

	rows, err := s.db.Query(s.context, getUserQuery, u)
	if err != nil {
		slog.Info("")
	}

	var users []entity.User

	for rows.Next() {
		var row entity.User
		if err := rows.Scan(&row); err != nil {
			slog.Info("Failed to find the users: ", err)
			ch <- nil
		}
	}

	ch <- users
}

func (s Storage) DeleteUser(ID int, ch chan int) {
	deleteUserQuery := "DELETE FROM users WHERE id=$1"

	res, err := s.db.Exec(s.context, deleteUserQuery, ID)
	fmt.Println(res.String())
	if err != nil {
		slog.Info("Failed to remove the user: ", err)
		ch <- http.StatusInternalServerError
	}

	ch <- http.StatusOK
}

func (s Storage) PatchUser(u entity.User, ch chan int) {
	updateUserQuery := "UPDATE users SET name=$2, surname=$3, patronymic=$4 WHERE id=$1"

	result, err := s.db.Exec(s.context, updateUserQuery, u)
	fmt.Println(result.String())
	if err != nil {
		slog.Info("Failed to update the user data: ", err)
		ch <- http.StatusInternalServerError
	}

	ch <- http.StatusOK
}

func (s Storage) AddUser(u User, ch chan int) {
	createUserQuery := "INSERT INTO users(name, surname, patronymic) VALUES ($1, $2, $3)"

	result, err := s.db.Exec(s.context, createUserQuery, u)
	fmt.Println(result.String())
	if err != nil {
		slog.Info("failed to add the user: ", err)
		ch <- http.StatusInsufficientStorage
	}

	ch <- http.StatusCreated
}
