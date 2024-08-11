package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vlad19930514/webApp/util"
)

type createUserRequest struct {
	FirstName string json:"firstname" binding:"required,alpha"
	LastName  string json:"lastname" binding:"required,alpha"
	Email     string json:"email" binding:"required,email"
	Age       int16  json:"age" binding:"required,min=1,max=130"
}

func (server *Server) createUser(ctx *gin.Context) {

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
	user, err := server.services.userService.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)

}

type getUserRequest struct {
	ID string uri:"id"  binding:"required,uuid"
}

func (server *Server) getUser(ctx *gin.Context) {
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

	user, err := server.services.userService.GetUser(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusAccepted, user)
}

type updateUserRequest struct {
	ID        string json:"id" binding:"required,uuid"
	FirstName string json:"firstname" binding:"required,alpha"
	LastName  string json:"lastname" binding:"required,alpha"
	Email     string json:"email" binding:"required,email"
	Age       int16  json:"age" binding:"required,min=1,max=130"
}

func (server *Server) updateUser(ctx *gin.Context) {
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
	domainUser, err := updateUserToDomain(req)
	user, err := server.services.userService.UpdateUser(ctx, domainUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}