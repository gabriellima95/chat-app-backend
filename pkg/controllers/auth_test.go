package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestAuthMiddleware(t *testing.T) {
	t.Run("case=auth-middleware-must-return-success-when-receiving-valid-token", func(t *testing.T) {
		var b bytes.Buffer
		req := httptest.NewRequest(http.MethodPost, "/test", &b)
		token, err := generateToken(uuid.NewString())
		if err != nil {
			t.Errorf("Error generating token: %s", err.Error())
		}
		req.Header.Set("Authorization", token)
		w := httptest.NewRecorder()

		AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {})(w, req)
		if w.Code != 200 {
			t.Errorf("Error authenticating token: response status code is not 200")
		}
	})

	t.Run("case=auth-middleware-must-return-unauthorized-when-receiving-invalid-token", func(t *testing.T) {
		var b bytes.Buffer
		req := httptest.NewRequest(http.MethodPost, "/test", &b)
		token := "abc"
		req.Header.Set("Authorization", token)
		w := httptest.NewRecorder()

		AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {})(w, req)
		if w.Code != 401 {
			t.Errorf("Error authenticating token: response status code is not 401")
		}
	})
}
