package middleware

import (
	"fmt"
	"net/http"
	"rest-api-example/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(user models.Login) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = user.Username
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something went wrong : %s", err.Error())
		return "", err
	}
	return tokenString, nil

}

var mySigningKey = []byte("SECRET_KET_!123")

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
			if token.Valid {
				endpoint(w, r)
			}

		} else {

			fmt.Fprintf(w, "Not Authorized", http.StatusForbidden)
		}

	})
}
