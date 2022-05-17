package auth

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"in-backend/internal/config"
	"in-backend/internal/handlers"
	"in-backend/internal/user"
	"in-backend/pkg/logging"
	"net/http"
	"regexp"
)

const (
	registerUrl = "/register"
)

type handler struct {
	logger     *logging.Logger
	repository user.Repository
	cfg        *config.Config
}

func (h *handler) RegisterAdmin(router *gin.RouterGroup) {
	router.POST(registerUrl, h.CreateUser)
}

func (h *handler) RegisterAuth(router *gin.RouterGroup) {
	//router.POST(registerUrl, h.CreateUser)
}

func NewHandler(logger *logging.Logger, repository user.Repository, cfg *config.Config) handlers.HandlerAuth {
	return &handler{
		logger:     logger,
		repository: repository,
		cfg:        cfg,
	}
}

type register struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Familia  string `form:"familia" json:"familia" binding:"required"`
	Name     string `form:"name" json:"name" binding:"required"`
}

func (h *handler) CreateUser(ctx *gin.Context) {
	var registerVals register
	if err := ctx.ShouldBind(&registerVals); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    500,
			"message": "missing vals",
		})
		return
	}
	usernameConvention := h.cfg.Pattern.User
	re, _ := regexp.Compile(usernameConvention)
	if !re.MatchString(registerVals.Username) || len(registerVals.Username) < 4 {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    500,
			"message": "wrong username",
		})
		return
	}
	passwordConvention := h.cfg.Pattern.Password
	re, _ = regexp.Compile(passwordConvention)
	if !re.MatchString(registerVals.Password) || len(registerVals.Password) <= 5 {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    500,
			"message": "wrong password",
		})
		return
	}
	var usr = &user.User{
		Username:     registerVals.Username,
		PasswordHash: registerVals.Password,
		Role:         "user",
		Familia:      registerVals.Familia,
		Name:         registerVals.Name,
	}
	password, err := bcrypt.GenerateFromPassword([]byte(usr.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Errorf("error to bcrypt password %t", err)
		ctx.String(http.StatusInternalServerError, "error %t", err)
		return
	}
	usr.PasswordHash = string(password)
	if err := h.repository.Create(context.TODO(), usr); err != nil {
		h.logger.Errorf("Error create user %t", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "error user is duplicate",
		})
		return
	}
	//token := jwt.New(jwt.GetSigningMethod("HS256"))
	//claims := token.Claims.(jwt.MapClaims)
	//claims["id"] = usr.ID
	//expire := time.Now().Add(time.Hour)
	//claims["exp"] = expire.Unix()
	//claims["orig_iat"] = time.Now().Unix()
	//tokenString, err := h.signedString(token)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	ctx.JSON(http.StatusOK, usr)
}

func (h *handler) signedString(token *jwt.Token) (string, error) {
	var tokenString string
	var err error
	tokenString, err = token.SignedString([]byte(h.cfg.Secret))
	return tokenString, err
}
