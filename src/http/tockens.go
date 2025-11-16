package http

import (
	"avito_pr_service/src/conf"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type claim struct {
	UserRole string `json:"user_role"`
	jwt.RegisteredClaims
}

const roleAdmin, roleUser = "admin", "user"

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func generateTokens() (adminTocken, userTocken string, err error) {
	generateToken := func(role string) (tokenString string, err error) {
		claim := &claim{UserRole: role}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		tokenString, err = token.SignedString(jwtKey)
		return
	}

	if adminTocken, err = generateToken(roleAdmin); err != nil {
		return
	}
	if userTocken, err = generateToken(roleUser); err != nil {
		return
	}
	return
}

func authenticationRequest(gctx *gin.Context) (role string) {
	const (
		authHeader, prefix = "Authorization", "Bearer "
		prefixLength       = 7
	)

	tokenString := gctx.GetHeader(authHeader)
	if tokenString == "" {
		conf.Logger.Error(fmt.Sprintf("%s: error on authentication: token not found", conf.LogHeaders.HTTPServer))
		return
	}
	if len(tokenString) > prefixLength && tokenString[:prefixLength] == prefix {
		tokenString = tokenString[prefixLength:]
	}
	claim := &claim{}
	if token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	}); err != nil || !token.Valid {
		conf.Logger.Error(fmt.Sprintf("%s: authorization falid (error value): %v", conf.LogHeaders.HTTPServer, err))
		return
	}
	role = claim.UserRole
	return
}
