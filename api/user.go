package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/kaushikkampli/neobank/db/sqlc"
	"github.com/kaushikkampli/neobank/utils"
	"github.com/lib/pq"
)

type CreateUserRequest struct {
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	Username        string    `json:"username"`
	FullName        string    `json:"full_name"`
	EmailID         string    `json:"email_id"`
	PasswdUpdatedAt time.Time `json:"passwd_updated_at"`
	CreatedAt       time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashPasswd, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
	}

	arg := db.CreateUserParams{
		Username:     req.UserName,
		HashedPasswd: hashPasswd,
		FullName:     req.FullName,
		EmailID:      req.Email,
	}

	User, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
	}

	resp := DBtoUserResp(User)

	ctx.JSON(http.StatusOK, resp)
}

func DBtoUserResp(User db.User) CreateUserResponse {
	return CreateUserResponse{
		Username:        User.Username,
		FullName:        User.FullName,
		EmailID:         User.EmailID,
		CreatedAt:       User.CreatedAt,
		PasswdUpdatedAt: User.PasswdUpdatedAt,
	}
}

type LoginUserRequest struct {
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	Token string             `json:"token"`
	User  CreateUserResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req LoginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.UserName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	err = utils.ComparePassword(req.Password, user.HashedPasswd)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	token, err := server.token.CreateToken(req.UserName, server.config.TokenExpirationTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := LoginUserResponse{
		Token: token,
		User:  DBtoUserResp(user),
	}

	ctx.JSON(http.StatusOK, resp)
}
