package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/0xForked/goca/server/hof"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

type handler struct {
	service IUserService
	cal     *calendar.Service
}

func (h handler) login(ctx *gin.Context) {
	var body LoginForm
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity,
			gin.H{"error": err.Error()})
		return
	}
	if err := body.Validate(); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity,
			gin.H{"error": err})
		return
	}
	// Call the service
	data, err := h.service.Login(ctx, &body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": err.Error()})
		return
	}
	ctx.SetCookie("ACCESS_TOKEN", data["token"].(string),
		1800, "/", "", false, true)
	ctx.JSON(http.StatusOK,
		gin.H{"data": data})
}

func (h handler) logout(ctx *gin.Context) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "ACCESS_TOKEN",
		Value:   "",
		MaxAge:  -1,
		Domain:  "http://localhost:8000",
		Path:    "/",
		Expires: time.Now().Add(-time.Hour),
	})
	ctx.JSON(http.StatusUnauthorized, nil)
}

func (h handler) profile(ctx *gin.Context) {
	var username string
	if uname, ok := ctx.MustGet("uname").(string); ok {
		username = uname
	}
	data, err := h.service.Profile(ctx, username, false)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity,
			gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK,
		gin.H{"data": data})
}

func (h handler) availability(ctx *gin.Context) {
	var uid int
	if id, ok := ctx.MustGet("uid").(float64); ok {
		uid = int(id)
	}
	data, err := h.service.Availability(ctx, uid)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity,
			gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

func (h handler) eventType(ctx *gin.Context) {
	var uid int
	if id, ok := ctx.MustGet("uid").(float64); ok {
		uid = int(id)
	}
	data, err := h.service.EventType(ctx, uid)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity,
			gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

func (h handler) event(ctx *gin.Context) {
	cfg := hof.GetOAuthConfig()
	var url string
	tok := &oauth2.Token{}
	var username string
	if uname, ok := ctx.MustGet("uname").(string); ok {
		username = uname
	}
	// get user data
	data, err := h.service.Profile(ctx, username,
		false)
	if err != nil || !data.Token.Valid {
		url, tok = hof.GetOAuthTokenFromWeb(cfg)
		if url != "" {
			ctx.JSON(http.StatusOK, gin.H{"auth_url": url})
			return
		}
	}
	if data != nil {
		if err = json.Unmarshal([]byte(data.Token.String), tok); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity,
				gin.H{"error": err.Error()})
			return
		}
	}
	// get user calendar data
	calendarService := hof.GetCalendarService(ctx, tok, cfg)
	event, err := hof.GetCalendarData(calendarService)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity,
			gin.H{"error": err.Error()})
		return
	}
	// get user profile from oauth
	peopleService := hof.GetPeopleService(ctx, tok, cfg)
	profile, err := hof.GetProfileData(peopleService)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity,
			gin.H{"error": err.Error()})
		return
	}
	//return data
	ctx.JSON(http.StatusOK, gin.H{
		"name":   profile,
		"events": event,
	})
}

func (h handler) exchange(ctx *gin.Context) {
	cfg := hof.GetOAuthConfig()
	var username string
	if uname, ok := ctx.MustGet("uname").(string); ok {
		username = uname
	}
	data, err := cfg.Exchange(ctx, ctx.Query("code"))
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity,
			gin.H{"error": err.Error()})
		return
	}
	if err := h.service.SaveGoogleToken(
		ctx, username, data,
	); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity,
			gin.H{"error": err.Error()})
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, "/fe/")
}

func newUserHandler(
	service IUserService,
	router *gin.RouterGroup,
) {
	h := &handler{service: service}
	router.POST("/login", h.login)
	router.POST("/logout", h.logout)
	router.GET("/profile", hof.Auth, h.profile)
	router.GET("/profile/availabilities", hof.Auth, h.availability)
	router.GET("/profile/event-types", hof.Auth, h.eventType)
	router.GET("/profile/events", hof.Auth, h.event)
	router.GET("/profile/:provider/exchange", hof.Auth, h.exchange)
}
