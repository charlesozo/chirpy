package auth

import (
	"fmt"
	"net/http"
	"time"
    "errors"
	"strings"
	"strconv"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)
type TokenType string
var ErrNoAuthHeaderIncluded = errors.New("not auth header included in request")
const (
	TokenTypeAccess TokenType= "chirpy-access"
	TokenTypeRefresh TokenType = "chirpy-refresh"
)
type MyCustomsClaims struct{
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err!=nil{
		return "", err
	}
	return string(dat), err
}
func ComparePassword(password, hash string) error{
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID int, tokenSecret string, expiresIn time.Duration, tokenType TokenType) (string, error) {
	signingKey := []byte(tokenSecret)
	claims := MyCustomsClaims{
		jwt.RegisteredClaims{
			Issuer: string(tokenType),
			IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}
func RefreshToken(tokenString, tokenSecret string) (string, error){
	token, err := jwt.ParseWithClaims(
		tokenString,
		&MyCustomsClaims{},
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil{
		return "", err
	}
	userIDString, err := token.Claims.GetSubject()
	if err != nil{
		return "", err
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil{
		return "", err
	}
	if issuer != string(TokenTypeRefresh){
		return "", errors.New("invalid issuer")
	}
	userid, err := strconv.Atoi(userIDString)
	if err != nil{
       return "", err
	}
	newToken, err := MakeJWT(userid, tokenSecret, time.Hour, TokenTypeAccess)
	if err !=nil{
		return "", err
	}
	return newToken, nil
}
func ValidateJWT(tokenString, tokenSecret string) (string, error){
   token, err := jwt.ParseWithClaims(
	tokenString,
	&MyCustomsClaims{},
    func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
)
 if err !=nil{
	return "", err
 }
 userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil{
		return "", err
	}
	if issuer != string(TokenTypeAccess) {
		return "", errors.New("invalid issuer")
	}

	return userIDString, nil
}

func GetBearerToken(headers http.Header) (string, error){
	authHeader := headers.Get("Authorization")
	if authHeader=="" {
		return "",  ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}
	return splitAuth[1], nil
}
func GetAPIKey(headers http.Header) (string, error){
	authHeader := headers.Get("Authorization")
	if authHeader=="" {
		return "",  ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}
	return splitAuth[1], nil
}