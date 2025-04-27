package service

import (
	"fmt"
	"io"
	"sync"

	"log/slog"
	"net/http"
	"strconv"

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
	wg             = sync.WaitGroup{}
)

func getDataOnUrl(url string) string {
	response, err := http.Get(url)
	if err != nil {
		slog.Error("failed to get an answer: ", err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error("failed to get the body of the request: ", err)
	}
	slog.Info(string(data))

	return string(data)
}

func enrichment(u entity.User) GetUsersRespnse {
	stringAge := getDataOnUrl(fmt.Sprintf("https://api.agify.io/?name=%s", u.Name))
	gender := getDataOnUrl(fmt.Sprintf("https://api.genderize.io/?name=%s", u.Name))
	nationality := getDataOnUrl(fmt.Sprintf("https://api.nationalize.io/?name=%s", u.Name))

	age, err := strconv.Atoi(stringAge)
	if err != nil {
		slog.Error("failed to convert the string to the number", err)
	}

	return GetUsersRespnse{
		ID:          u.ID,
		Name:        u.Name,
		Surname:     u.Surname,
		Patronymic:  u.Patronymic,
		Age:         age,
		Gender:      gender,
		nationality: nationality,
	}

}

// @Summary      create user on database
// @Description  create user by name, surname and optional patronymic
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        input body storage.User true "user full name"
// @Success      200  {object}  entity.User
// @Failure      507  {string} error
// @Router       /v1/user/ [post]
func (s Service) AddUser(ctx *gin.Context) {
	var user storage.User
	ctx.BindJSON(&user)

	go s.Storage.AddUser(user, statusCodeChan)

	result := <-statusCodeChan

	if result > 299 {
		ctx.JSON(result, "failed to add the user")
		return
	}
	ctx.JSON(result, user)
}

// @Summary      update user on database
// @Description  update user by name, surname and optional patronymic
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        input body 	storage.User true "user full name"
// @Success      200  {object}  entity.User
// @Failure      507  {string} 	error
// @Router       /v1/user/ [patch]
func (s Service) PatchUser(ctx *gin.Context) {
	var user entity.User
	ctx.BindJSON(&user)

	go s.Storage.PatchUser(user, statusCodeChan)

	result := <-statusCodeChan

	if result > 299 {
		ctx.JSON(result, "failed to patch the user")
		return
	}
	ctx.JSON(result, user)
}

// @Summary      get user on database
// @Description  get user by name, surname and optional patronymic
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        name,surname,patronymic	query string	false "optional user full name"
// @Success      200  {object}  GetUsersRespnse
// @Failure      507  {string}	error
// @Router       /v1/user/ [get]
func (s Service) GetUsers(ctx *gin.Context) {
	user := storage.User{
		Name:       ctx.Query("name"),
		Surname:    ctx.Query("surnme"),
		Patronymic: ctx.Query("patronymic"),
	}

	go s.Storage.GetUsers(user, usersChan)

	userArray := <-usersChan
	if userArray == nil {
		ctx.JSON(http.StatusInternalServerError, "Failed to find the users")
		return
	}

	var result []GetUsersRespnse

	for i := range userArray {
		result = append(result, enrichment(userArray[i]))
	}

	ctx.JSON(http.StatusOK, result)
}

// @Summary      delete user on database
// @Description  delete user on database
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path     	int  true  "delete user"
// @Success      200  {int}  	"id"
// @Failure      507  {string} 	"error"
// @Router       /v1/user/{id} 	[delete]
func (s Service) DeleteUser(ctx *gin.Context) {
	var stringID string = ctx.Param("id")
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		slog.Error("failed to convert user id to number", err)
	}
	s.Storage.DeleteUser(ID, statusCodeChan)

	result := <-statusCodeChan
	if result > 299 {
		ctx.JSON(result, "user with such id does not exist")
		return
	}
	ctx.JSON(result, ID)
}
