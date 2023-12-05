package middleware

import (
	"net/http"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/dgrijalva/jwt-go"
)

var key, _ = beego.AppConfig.String("JWT_SEC_KEY")
var jwtKey = []byte(key)

func JWTMiddleware(ctx *context.Context) {
	tokenString := ctx.Input.Header("Authorization")
	if tokenString == "" {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Unauthorized"}, true, false)
		return
	}
	bearer := ContainsBearer(tokenString)
	if bearer {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Invalid token"}, true, false)
		return
	}
	ctx.Input.SetData("user", token.Claims.(jwt.MapClaims))
}

func ContainsBearer(token string) bool {
	// Convert the token to lowercase to make the comparison case-insensitive
	lowerToken := strings.ToLower(token)

	// Check if the token starts with "bearer "
	return strings.HasPrefix(lowerToken, "bearer ")
}