package handler

import (
	"log"
	"net/http"
	authdto "restaurant/dto/auth"
	dto "restaurant/dto/result"
	"restaurant/models"
	"restaurant/pkg/bcrypt"
	jwtToken "restaurant/pkg/jwt"
	"restaurant/repositories"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type handlerAuth struct {
	AuthRepository repositories.AuthRepository
}

func HandlerAuth(AuthRepository repositories.AuthRepository) *handlerAuth {
	return &handlerAuth{AuthRepository}
}

func (h *handlerAuth) Register(c echo.Context) error {
	request := new(authdto.RegisterRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	password, err := bcrypt.HashingPassword(request.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: password,
	}

	data, err := h.AuthRepository.Register(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Message: "Your registration is success", Data: data})
}

func (h *handlerAuth) Login(c echo.Context) error {
	request := new(authdto.LoginRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	data := models.User{
		Email:    request.Email,
		Password: request.Password,
	}

	userLogin, err := h.AuthRepository.Login(data.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: "email not registered"})
	}

	isValid := bcrypt.CheckPasswordHash(request.Password, userLogin.Password)
	if !isValid {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: "Incorrect email or password"})
	}

	claims := jwt.MapClaims{}
	claims["id"] = userLogin.ID
	// claims["listAs"] = userLogin.ListAs
	// claims["role"] = userLogin.Role
	claims["exp"] = time.Now().Add(time.Hour * 4).Unix() // 4 hours expired

	token, generateTokenErr := jwtToken.GenerateToken(&claims)
	if generateTokenErr != nil {
		log.Println(generateTokenErr)
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	loginResponse := authdto.LoginResponse{
		ID:    userLogin.ID,
		Email: userLogin.Email,
		Token: token,
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: loginResponse})
}

func (h *handlerAuth) CheckAuth(c echo.Context) error {
	userLogin := c.Get("userLogin")
	userId := userLogin.(jwt.MapClaims)["id"].(float64)

	user, _ := h.AuthRepository.CheckAuth(int(userId))

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: responseCheckAuth(user)})
}

func responseCheckAuth(u models.User) authdto.CheckAuthResponse {
	return authdto.CheckAuthResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
