package auth

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"in-backend/internal/config"
	"in-backend/internal/handlers"
	"in-backend/internal/user"
	"in-backend/pkg/logging"
	"net/http"
	"regexp"
	"time"
)

const (
	registerUrl = "/register"
)

type handler struct {
	logger     *logging.Logger
	repository user.Repository
	cfg        *config.Config
}

func NewHandler(logger *logging.Logger, repository user.Repository, cfg *config.Config) handlers.Handler {
	return &handler{
		logger:     logger,
		repository: repository,
		cfg:        cfg,
	}
}

func (h *handler) Register(router *gin.Engine) {
	router.POST(registerUrl, h.CreateUser)
}

type register struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Familia  string `form:"familia" json:"familia" binding:"required"`
	Role     string `form:"role" json:"role" binding:"required"` //TODO role check user or admin.
	Name     string `form:"name" json:"name" binding:"required"`
}

func (h *handler) CreateUser(ctx *gin.Context) {
	var registerVals register
	if err := ctx.ShouldBind(&registerVals); err != nil {
		ctx.String(http.StatusUnauthorized, "missing vals")
		return
	}
	usernameConvention := h.cfg.Pattern.User
	re, _ := regexp.Compile(usernameConvention)
	if !(len(registerVals.Username) > 4 && re.MatchString(registerVals.Username)) {
		ctx.String(http.StatusUnauthorized, "wrong username")
		return
	}
	passwordConvention := h.cfg.Pattern.Password
	re, _ = regexp.Compile(passwordConvention)
	if !(len(registerVals.Password) <= 12 && re.MatchString(registerVals.Password)) {
		ctx.String(http.StatusUnauthorized, "wrong password")
		return
	}
	var usr = &user.User{
		Username: registerVals.Username,
		Password: registerVals.Password,
		Role:     registerVals.Role,
		Familia:  registerVals.Familia,
		Name:     registerVals.Name,
	}
	password, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Errorf("error to bcrypt password %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
		return
	}
	usr.PasswordHash = string(password)
	if err := h.repository.Create(context.TODO(), usr); err != nil {
		h.logger.Errorf("Error create user %t", err)
		ctx.String(http.StatusInternalServerError, "error user is duplicate")
		return
	}
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = usr.ID
	expire := time.Now().Add(time.Hour)
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = time.Now().Unix()
	tokenString, err := h.signedString(token)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	})
}

func (h *handler) signedString(token *jwt.Token) (string, error) {
	var tokenString string
	var err error
	tokenString, err = token.SignedString([]byte(h.cfg.Secret))
	return tokenString, err
}
