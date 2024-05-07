package hof

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {
	var accessToken string
	// token from cookie
	if cookie, err := ctx.Request.Cookie("ACCESS_TOKEN"); err == nil {
		accessToken = cookie.Value
	}
	// token from header
	if authHeader := ctx.Request.Header.Get("Authorization"); authHeader != "" {
		header := strings.Split(authHeader, " ")
		accessToken = header[1]
	}
	// if token empty remove it if exist
	if accessToken == "" {
		http.SetCookie(ctx.Writer, &http.Cookie{
			Name:    "ACCESS_TOKEN",
			Value:   "",
			MaxAge:  -1,
			Domain:  "http://localhost:8000",
			Path:    "/",
			Expires: time.Now().Add(-time.Hour),
		})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "ACCESS_TOKEN_NOT_PROVIDE")
		return
	}
	// extract jwt
	claim, err := ExtractAndValidateJWT(TokenSecret, accessToken)
	if err != nil {
		http.SetCookie(ctx.Writer, &http.Cookie{
			Name:    "ACCESS_TOKEN",
			Value:   "",
			MaxAge:  -1,
			Domain:  "http://localhost:8000",
			Path:    "/",
			Expires: time.Now().Add(-time.Hour),
		})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}
	// set item
	ctx.Set("uid", claim.Payload["id"])
	ctx.Set("uname", claim.Payload["username"])
	ctx.Next()
}
