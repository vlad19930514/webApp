package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/vlad19930514/webApp/db/sqlc"
	"github.com/vlad19930514/webApp/util"
)

type createUserRequest struct {
	Firstname string `json:"firstname" binding:"required,alpha"` // a person can have a home and cottage...
	Lastname  string `json:"lastname" binding:"required,alpha"`
	Email     string `json:"email" binding:"required,email"`
	Age       int16  `json:"age" binding:"required,min=1,max=130"`
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

	// Преобразуем ошибки в JSON-формат

	generatedUUID, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("Ошибка при создании UUID: %v", err)
	}
	dbUUID := pgtype.UUID{}
	copy(dbUUID.Bytes[:], generatedUUID[:]) // Устанавливаем bytes field
	dbUUID.Valid = true                     // Устанавливаем valid field

	id, _ := dbUUID.UUIDValue()
	arg := db.CreateUserParams{
		ID:        id,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Email:     req.Email,
		Age:       req.Age,
		Created:   util.PgtypeCurrentTime(),
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusAccepted, user)

}

type getUserRequest struct {
	ID string `uri:"id" binding:"uuid"`
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
	dbUUID := pgtype.UUID{}
	dbUUID.Scan(req.ID)

	user, err := server.store.GetUser(ctx, dbUUID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusAccepted, user)
}

type updateUserRequest struct {
	ID        string `json:"id" binding:"required"`
	Firstname string `json:"firstname" binding:"required,alpha"`
	Lastname  string `json:"lastname" binding:"required,alpha"`
	Email     string `json:"email" binding:"required,email"`
	Age       int16  `json:"age"  binding:"required,min=1,max=130"`
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
	dbUUID := pgtype.UUID{}
	dbUUID.Scan(req.ID)
	arg := db.UpdateUserParams{
		ID:        dbUUID,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Email:     req.Email,
		Age:       req.Age,
	}
	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}
