package service

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Flikest/testovoe-effective-mobile/internal/entity"
	"github.com/Flikest/testovoe-effective-mobile/internal/storage"
	"github.com/gin-gonic/gin"
)

type GetUsersRespnse struct {
	ID          int
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	nationality string
}

var (
	statusCodeChan = make(chan int)
	usersChan      = make(chan []entity.User)
	enrichmentChan = make(chan []GetUsersRespnse)
)

func getDataOnUrl(url string) string {
	response, err := http.Get(url)
	if err != nil {
		slog.Info("failed to get an answer: ", err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Info("failed to get the body of the request: ", err)
	}
	slog.Info(string(data))

	return string(data)
}

func enrichment(u entity.User) GetUsersRespnse {
	age := getDataOnUrl(fmt.Sprintf("https://api.agify.io/?name=%s", u.Name))

	resp_age.Body.Read([]byte())

	resp_gender, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", u.Name))
	if err != nil {
		slog.Info("gender not found: ", err)
	}

	resp_nationality, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", u.Name))
	if err != nil {
		slog.Info("nationality not found: ", err)
	}
	if err != nil {
		slog.Info("age not found: ", err)
	}

	return GetUsersRespnse{
		ID:          u.ID,
		Name:        u.Name,
		Surname:     u.Surname,
		Patronymic:  u.Patronymic,
		Age:         int(age),
		Gender:      gender,
		nationality: nationality,
	}

}

func (s Service) AddUser(ctx *gin.Context) {
	var user storage.User
	ctx.BindJSON(&user)

	go s.Storage.AddUser(user, statusCodeChan)

	result := <-statusCodeChan

	if result > 299 {
		ctx.JSON(http.StatusInternalServerError, "failed to add the user")
	}
	ctx.JSON(result, user)
}

func (s Service) PatchUser(ctx *gin.Context) {
	var user entity.User
	ctx.BindJSON(&user)

	go s.Storage.PatchUser(user, statusCodeChan)

	result := <-statusCodeChan

	if result > 299 {
		ctx.JSON(http.StatusInternalServerError, "failed to patch the user")
	}
	ctx.JSON(result, user)
}

func (s Service) GetUsers(ctx *gin.Context) {
	user := storage.User{
		Name:       ctx.Query("name"),
		Surname:    ctx.Query("surnme"),
		Patronymic: ctx.Query("patronymic"),
	}

	go s.Storage.GetUsers(user, usersChan)

	result := <-usersChan
	if result == nil {
		ctx.JSON(http.StatusInternalServerError, "Failed to find the users")
	}
	ctx.JSON(http.StatusOK, result)
}
