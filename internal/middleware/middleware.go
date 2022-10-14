package middleware

import (
	"fmt"
	"githib.com/dkischenko/company-api/pkg/logger"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"strings"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l, err := logger.GetLogger()
		if err != nil {
			panic(fmt.Sprintf("cannot init logger: %w", err))
		}
		start := time.Now()
		next.ServeHTTP(w, r)
		l.Entry.Logger.Infof("Method: %s | Reqest: %s | Latency: %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func PanicAndRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			l, err := logger.GetLogger()
			if err != nil {
				panic(fmt.Sprintf("cannot init logger: %w", err))
			}
			if err := recover(); err != nil {
				l.Entry.Logger.Errorf("panic: %+v", err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l, err := logger.GetLogger()
		if err != nil {
			panic(fmt.Sprintf("cannot init logger: %w", err))
		}

		if r.Method != http.MethodGet && strings.Contains(r.URL.Path, "companies") {
			l.Entry.Logger.Infof("just some info for logger")
			tokenString := r.Header.Get("Authorization")
			if len(tokenString) == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Missing Authorization Header"))
				l.Entry.Logger.Warning("Missing Authorization Header")
				return
			}

			tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
			_, err := verifyToken(tokenString)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(fmt.Sprintf("Error verifying JWT token: %w", err)))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func verifyToken(tokenString string) (user_id string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SIGNINKEY")), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid || claims["user_id"] == nil {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["user_id"].(string), nil
}
