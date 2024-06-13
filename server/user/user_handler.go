package user

import (
	"encoding/json"
	"net/http"
	"sync"
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
	var username string
	if uname, ok := ctx.MustGet("uname").(string); ok {
		username = uname
	}
	data, err := h.service.EventType(ctx, uid, username)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity,
			gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

func (h handler) event(ctx *gin.Context) {
	var googleEvents []*calendar.Event
	var username, googleAuthURL, microsoftAuthURL,
		googleEmail, googleName, microsoftName, microsoftEmail string
	var wg sync.WaitGroup
	var mu sync.Mutex
	googleAuthToken, microsoftAuthToken := &oauth2.Token{}, &oauth2.Token{}
	// get user internal data
	if uname, ok := ctx.MustGet("uname").(string); ok {
		username = uname
	}
	data, err := h.service.Profile(ctx, username, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": err.Error()})
		return
	}
	wg.Add(2)
	go func() {
		defer wg.Done()
		if data != nil && data.GoogleToken.Valid {
			if err = json.Unmarshal([]byte(data.GoogleToken.String), googleAuthToken); err != nil {
				ctx.JSON(http.StatusUnprocessableEntity,
					gin.H{"error": err.Error()})
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		if data != nil && data.MicrosoftToken.Valid {
			if err = json.Unmarshal([]byte(data.MicrosoftToken.String), microsoftAuthToken); err != nil {
				ctx.JSON(http.StatusUnprocessableEntity,
					gin.H{"error": err.Error()})
				return
			}
		}
	}()
	wg.Wait()
	// get user google data
	if googleAuthToken != nil {
		cfg := hof.GetGoogleOAuthConfig()
		if data != nil && !data.GoogleToken.Valid {
			googleAuthURL, googleAuthToken = hof.GetGoogleOAuthTokenFromWeb(cfg)
		}
		if googleAuthURL == "" {
			wg.Add(2)
			// get user calendar data
			go func() {
				defer wg.Done()
				calendarService := hof.GetGoogleCalendarService(ctx, googleAuthToken, cfg)
				event, err := hof.GetGoogleCalendarData(calendarService)
				if err != nil {
					mu.Lock()
					defer mu.Unlock()
					ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
					return
				}
				mu.Lock()
				googleEvents = event
				mu.Unlock()
			}()
			// get user profile from oauth
			go func() {
				defer wg.Done()
				displayName, emailAddress, err := hof.GetGoogleUserData(ctx, googleAuthToken, cfg)
				if err != nil {
					mu.Lock()
					defer mu.Unlock()
					ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
					return
				}
				mu.Lock()
				googleName = displayName
				googleEmail = emailAddress
				mu.Unlock()
			}()
			wg.Wait()
		}
	}
	// get user microsoft data
	if microsoftAuthToken != nil {
		cfg := hof.GetMicrosoftOAuthConfig()
		if data != nil && !data.MicrosoftToken.Valid {
			microsoftAuthURL, microsoftAuthToken = hof.GetMicrosoftOAuthTokenFromWeb(cfg)
		}
		if microsoftAuthURL == "" {
			wg.Add(1)
			// get user profile
			go func() {
				defer wg.Done()
				userData, err := hof.GetMicrosoftUserProfile(microsoftAuthToken.AccessToken)
				if err != nil {
					mu.Lock()
					defer mu.Unlock()
					ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
					return
				}
				mu.Lock()
				if userData != nil {
					microsoftName = userData["name"]
					microsoftEmail = userData["email"]
				}
				mu.Unlock()
			}()
			wg.Wait()
		}
	}
	//return data
	ctx.JSON(http.StatusOK, gin.H{
		"google_name":        googleName,
		"google_email":       googleEmail,
		"google_scheduled":   googleEvents,
		"google_auth_url":    googleAuthURL,
		"microsoft_auth_url": microsoftAuthURL,
		"microsoft_name":     microsoftName,
		"microsoft_email":    microsoftEmail,
	})
}

func (h handler) googleExchange(ctx *gin.Context) {
	var username string
	if uname, ok := ctx.MustGet("uname").(string); ok {
		username = uname
	}
	cfg := hof.GetGoogleOAuthConfig()
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

func (h handler) microsoftExchange(ctx *gin.Context) {
	var username string
	if uname, ok := ctx.MustGet("uname").(string); ok {
		username = uname
	}
	cfg := hof.GetMicrosoftOAuthConfig()
	data, err := cfg.Exchange(ctx, ctx.Query("code"))
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity,
			gin.H{"error": err.Error()})
		return
	}
	if err := h.service.SaveMicrosoftToken(
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
	router.GET("/profile/google/exchange", hof.Auth, h.googleExchange)
	router.GET("/profile/microsoft/exchange", hof.Auth, h.microsoftExchange)
}
