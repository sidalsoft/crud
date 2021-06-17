package middleware

import (
	"encoding/base64"
	"github.com/sidalsoft/crud/pkg/security"
	"net/http"
	"strings"
)

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
