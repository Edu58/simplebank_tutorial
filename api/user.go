package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
)

type UserResponse struct {
	Username         string    `json:"username"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	PasswordChangeAt time.Time `json:"password_change_at"`
	CreatedAt        time.Time `json:"created_at"`
}

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var user createUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(user.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not hash password",
		})

		return
	}

	args := db.CreateUserParams{
		Username:       user.Username,
		HashedPassword: hashedPassword,
		FullName:       user.FullName,
		Email:          user.Email,
	}

	userRecord, err := server.store.CreateUser(ctx, args)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
	}

	response := UserResponse{
		Username: userRecord.Username,
		FullName: userRecord.Username,
		Email: userRecord.Email,
		CreatedAt: userRecord.CreatedAt,
		PasswordChangeAt: userRecord.PasswordChangeAt,
	}

	ctx.JSON(http.StatusOK, response)
}

type getUserRequest struct {
	Username string `uri:"username" binding:"required,alphanum"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var user getUserRequest

	if err := ctx.ShouldBindUri(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userRecord, err := server.store.GetUser(ctx, user.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response := UserResponse{
		Username: userRecord.Username,
		FullName: userRecord.Username,
		Email: userRecord.Email,
		CreatedAt: userRecord.CreatedAt,
		PasswordChangeAt: userRecord.PasswordChangeAt,
	}

	ctx.JSON(http.StatusOK, response)
}
