package httpserver

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/vlad19930514/webApp/internal/app/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vlad19930514/webApp/util"
)

type createUserRequest struct {
	FirstName string `json:"firstname" binding:"required,alpha"`
	LastName  string `json:"lastname" binding:"required,alpha"`
	Email     string `json:"email" binding:"required,email"`
	Age       uint8  `json:"age" binding:"required,min=1,max=130"`
}

func (server *HttpServer) createUser(ctx *gin.Context) {

	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {

		customErrors, errorsExist := util.GetValidationErrors(&err)

		if errorsExist {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": customErrors})
			return
		}

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg, err := toDomainUser(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	user, err := server.userService.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)

}

type getUserRequest struct {
	ID string `uri:"id"  binding:"required,uuid"`
}

func (server *HttpServer) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationErrors, errorsExist := util.GetValidationErrors(&err)
		fmt.Println(errorsExist)
		if errorsExist {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, _ := uuid.Parse(req.ID)

	user, err := server.userService.GetUser(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusAccepted, user)
}

type updateUserRequest struct {
	ID        uuid.UUID `json:"id" binding:"required,uuid"`
	FirstName string    `json:"firstname" binding:"required,alpha"`
	LastName  string    `json:"lastname" binding:"required,alpha"`
	Email     string    `json:"email" binding:"required,email"`
	Age       uint8     `json:"age" binding:"required,min=1,max=130"`
}

func (server *HttpServer) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {

		customErrors, errorsExist := util.GetValidationErrors(&err)

		if errorsExist {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": customErrors})
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	domainUser := domain.User{
		Id:        req.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Age:       req.Age,
	}

	user, err := server.userService.UpdateUser(ctx, domainUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}
func toDomainUser(req createUserRequest) (domain.User, error) {
	id := uuid.New()
	user := domain.User{
		Id:        id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Age:       req.Age,
	}
	return user, nil
}
