package net

//  Generate RSA signing files via shell (adjust as needed):
//
//  $ openssl genrsa -out app.rsa 1024
//  $ openssl rsa -in app.rsa -pubout > app.rsa.pub
//
// Code borrowed and modified from the following sources:
// https://www.youtube.com/watch?v=dgJFeqeXVKw
// https://goo.gl/ofVjK4
// https://github.com/dgrijalva/jwt-go
//

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"crypto/rsa"
)

type UserCredentials struct {
	Username string
	Password string
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

type JwtController struct {
	privKeyPath string //app.rsa
	pubKeyPath  string //app.rsa.pub

	authenticate func(UserCredentials) error

	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
}

func NewJwtController(privKeyPath, pubKeyPath string,
	authenticator func(UserCredentials) error) *JwtController {
	return &JwtController{privKeyPath: privKeyPath, pubKeyPath: pubKeyPath, authenticate: authenticator}
}

func (o *JwtController) Setup() (err error) {
	var keyBytes []byte
	if keyBytes, err = ioutil.ReadFile(o.privKeyPath); err == nil {
		o.signKey, err = jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	}

	if err != nil {
		return
	}

	if keyBytes, err = ioutil.ReadFile(o.pubKeyPath); err == nil {
		o.verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	}
	return
}

func (o *JwtController) LoginHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var user UserCredentials
		if err := Decode(&user, r); err != nil {
			ResponseResultErr(err, "Can't retrieve credentials", http.StatusForbidden, w)
			return
		}

		if err := o.authenticate(user); err != nil {
			ResponseResultErr(err, "Wrong credentials", http.StatusForbidden, w)
			return
		}

		token := jwt.New(jwt.SigningMethodRS256)
		claims := make(jwt.MapClaims)
		claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
		claims["iat"] = time.Now().Unix()
		token.Claims = claims

		if tokenString, err := token.SignedString(o.signKey); err != nil {
			ResponseResultErr(err, "Error while signing the token", http.StatusInternalServerError, w)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			ResponseJson(Token{tokenString}, w)
		}
	})
}

func (o *JwtController) ValidateTokenHandler(protected http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o.ValidateToken(w, r, protected)
	})
}

func (o *JwtController) ValidateToken(w http.ResponseWriter, r *http.Request, next http.Handler) {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return o.verifyKey, nil
		})

	if err == nil {
		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource")
	}
}
