package v1

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"

	"auth-service/internal/entities"
	"auth-service/internal/usecase"
)

type authRoutes struct {
	uc usecase.Auth
}

func newAuthRoutes(handler *gin.RouterGroup, uc usecase.Auth) {
	r := &authRoutes{uc}

	h := handler.Group("/auth")
	{
		h.POST("/register", r.register)
		h.POST("/login", r.login)

		h.GET("/register", func(c *gin.Context) {
			c.HTML(http.StatusOK, "register.html", nil)
		})
		h.GET("/login", r.loginGet)
		h.GET("/user_profile", r.userProfile)
		h.GET("/logout", r.logout)
	}
}

func (r *authRoutes) register(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-store")
	type RequestUser struct {
		Login       string `form:"login" binding:"required"`
		Email       string `form:"email" binding:"required"`
		Password    string `form:"password" binding:"required"`
		PhoneNumber string `form:"phone_number" binding:"required"`
	}

	var rt RequestUser
	if err := c.ShouldBind(&rt); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.uc.Register(c.Request.Context(),
		&entities.User{Login: rt.Login,
			Email:       rt.Email,
			Password:    rt.Password,
			PhoneNumber: rt.PhoneNumber,
		}); err != nil {
		log.Println(err.Error())
		if err.Error() == entities.ErrNotValidLogin.Error() {
			errorResponse(c, http.StatusBadRequest, entities.ErrNotValidLogin.Error())
			return
		}
		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	location := url.URL{Path: "/v1/auth/login"}
	c.Redirect(http.StatusMovedPermanently, location.RequestURI())
}

func (r *authRoutes) login(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-store")

	type RequestUser struct {
		Login    string `form:"login" binding:"required"`
		Password string `form:"password" binding:"required"`
	}
	var rt RequestUser
	if err := c.ShouldBind(&rt); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}
	user := &entities.User{
		Login:    rt.Login,
		Password: rt.Password,
	}

	// чекаем логин-пароль в базе
	isValidCreds, err := r.uc.CheckCreds(c.Request.Context(), user)
	if err != nil {
		log.Println(err.Error())
		if err.Error() == entities.ErrNotValidLoginOrPass.Error() {
			errorResponse(c, http.StatusInternalServerError, entities.ErrNotValidLoginOrPass.Error())
			return
		}
		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}
	if !isValidCreds {
		errorResponse(c, http.StatusBadRequest, entities.ErrNotValidLoginOrPass.Error())
		return
	}

	// создаем сессию и добавляем ее в бд
	sessionID, err := r.uc.GenerateCookie(c.Request.Context(), user)
	if err != nil {
		log.Println(err.Error())
		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.SetCookie("session_id", sessionID, 60*60*12, "/v1/auth/", "localhost", false, true)

	location := url.URL{Path: "/v1/auth/user_profile"}
	c.Redirect(http.StatusMovedPermanently, location.RequestURI())
}

func (r *authRoutes) userProfile(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-store")

	cookie, err := c.Request.Cookie("session_id")
	if err != nil {
		log.Println(err.Error())
		location := url.URL{Path: "/v1/auth/login"}
		c.Redirect(http.StatusMovedPermanently, location.RequestURI())
		return
	}

	login, err := r.uc.GetLogin(c.Request.Context(), cookie.Value)
	if err != nil {
		log.Println(err.Error())
		location := url.URL{Path: "/v1/auth/login"}
		c.Redirect(http.StatusMovedPermanently, location.RequestURI())
		return
	}

	user, err := r.uc.GetUser(c.Request.Context(), login)
	if user == nil {
		log.Println(err.Error())
		location := url.URL{Path: "/v1/auth/register"}
		c.Redirect(http.StatusMovedPermanently, location.RequestURI())
		return
	}

	c.HTML(http.StatusOK, "user_profile.html", user)
}

func (r *authRoutes) logout(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-store")

	cookie, _ := c.Request.Cookie("session_id")
	login, err := r.uc.GetLogin(c.Request.Context(), cookie.Value)
	if err != nil {
		log.Println(err.Error())
		location := url.URL{Path: "/v1/auth/login"}
		c.Redirect(http.StatusMovedPermanently, location.RequestURI())
		return
	}

	err = r.uc.DeleteSession(c.Request.Context(), login)
	if err != nil {
		log.Println(err.Error())
		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	location := url.URL{Path: "/v1/auth/login"}
	c.Redirect(http.StatusMovedPermanently, location.RequestURI())
}

func (r *authRoutes) loginGet(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-store")

	cookie, err := c.Request.Cookie("session_id")
	if err == nil {
		expireDate, err := r.uc.GetExpire(c.Request.Context(), cookie.Value)
		if err != nil {
			// кука не найдена в базе
			log.Println(err.Error())
			c.HTML(http.StatusOK, "login.html", nil)
			return
		}
		//fmt.Printf("expireDate.After(time.Now()): %v\n", expireDate.After(time.Now()))
		if expireDate.After(time.Now()) {
			// кука валидна
			location := url.URL{Path: "/v1/auth/user_profile"}
			c.Redirect(http.StatusMovedPermanently, location.RequestURI())
			return
		}
	}
	c.HTML(http.StatusOK, "login.html", nil)
}
