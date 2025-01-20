package restapi

import (
	"database/sql"
	"fmt"
	"protchain/internal/dal/crudder"
	"protchain/internal/dal/model"
	"protchain/internal/schema"
	"protchain/pkg/function"

	"github.com/lucsky/cuid"
	"github.com/pkg/errors"
)

func (a *API) RegisterUser(req schema.RegisterReq) (schema.RegisterRes, error) {
	if req.Email == "" || req.Password == "" {
		return schema.RegisterRes{}, errors.New("email and password cannot be empty")
	}

	uc := crudder.DefaultCrudder(&model.User{}, a.Deps.DAL.SqlDB)
	uc.Filter.Exact = map[string]any{
		"email": req.Email,
	}
	exists, err := uc.Exists()
	if err != nil {
		return schema.RegisterRes{}, errors.Wrap(err, "[DAL]: Failed to check if user exists")
	}
	if exists {
		return schema.RegisterRes{}, fmt.Errorf("you already have an account. Please login")
	}
	hashedPass, err := function.HashPassword(req.Password)
	if err != nil {
		return schema.RegisterRes{}, fmt.Errorf("you already have an account. Please login")
	}
	u := model.User{
		Id:           cuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPass,
	}
	uc.DataModel = &u
	if _, err := uc.Insert(); err != nil {
		return schema.RegisterRes{}, errors.Wrap(err, "[DAL]: Failed to insert user")
	}

	return schema.RegisterRes{UserID: u.Id}, nil
}

func (a *API) LoginUser(req schema.LoginReq) (schema.LoginRes, error) {
	if req.Email == "" || req.Password == "" {
		return schema.LoginRes{}, errors.New("email and password cannot be empty")
	}

	var u model.User
	uc := crudder.DefaultCrudder(&u, a.Deps.DAL.SqlDB)
	uc.Filter.Exact = map[string]any{
		"email": req.Email,
	}
	if err := uc.Fetch(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schema.LoginRes{}, fmt.Errorf("Invalid authentication credentials")
		}
		return schema.LoginRes{}, errors.Wrap(err, "[DAL]: Failed to retrieve user information")
	}

	if !function.CheckPasswordHash(req.Password, u.PasswordHash) {
		return schema.LoginRes{}, fmt.Errorf("Invalid authentication credentials")
	}

	token := function.GenerateJWT(req.Email, a.Config.JwtKey)
	return schema.LoginRes{Token: token}, nil
}

func (a *API) GoogleOAuth(req schema.GoogleOAuthReq) (schema.GoogleOAuthRes, error) {
	if req.GoogleToken == "" {
		return schema.GoogleOAuthRes{}, errors.New("google token cannot be empty")
	}

	// Simulate Google token validation (replace with actual validation)
	email := "user@example.com" // Mock email retrieval from Google token
	token := function.GenerateJWT(email, a.Config.JwtKey)
	return schema.GoogleOAuthRes{Token: token}, nil
}
