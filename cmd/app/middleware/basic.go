package middleware

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/sidalsoft/crud/pkg/security"
	"net/http"
	"strings"
)

var ErrNoAuthentication = errors.New("no authentication")

type IDFunc func(ctx context.Context, token string) (int64, error)

type HasAnyRoleFunc func(ctx context.Context, role string) bool

func CheckRole(hasAnyRoleFunc HasAnyRoleFunc, role string) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if !hasAnyRoleFunc(request.Context(), role) {
				http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			handler.ServeHTTP(writer, request)
		})
	}
}

func Authentication(ctx context.Context) (int64, error) {
	if value, ok := ctx.Value("idUser").(int64); ok {
		return value, nil
	}
	return 0, ErrNoAuthentication
}

func Authenticate(idFunc IDFunc) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			token := request.Header.Get("Authorization")

			id, err := idFunc(request.Context(), token)
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(request.Context(), "idUser", id)
			request = request.WithContext(ctx)
			handler.ServeHTTP(writer, request)
		})
	}
}

func Basic(authSvc *security.AuthService) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			authorization := request.Header.Get("Authorization")

			arr1 := strings.Split(string(authorization), " ")
			if len(arr1) < 2 {
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			data, err := base64.StdEncoding.DecodeString(arr1[1])
			if err != nil || len(data) == 0 {
				http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			arr := strings.Split(string(data), ":")
			login := arr[0]
			pass := arr[1]
			if !authSvc.Auth(login, pass) {
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			handler.ServeHTTP(writer, request)
		})
	}
}
