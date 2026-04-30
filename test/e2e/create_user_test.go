package e2e

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	ID    int64
	Name  string
	Email string
}

func TestCreateUser(t *testing.T) {
	t.Run("creating user", func(t *testing.T) {
		payload := `{"name": "Test", "email": "test_user@test.com", "password": "12345678"}`
		reader := strings.NewReader(payload)

		resp, err := http.Post(gatewayAddr+"/api/v1/users", "application/json", reader)
		assert.NoError(t, err)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("creating existing user", func(t *testing.T) {
		payload := `{"name": "AnotherTest", "email": "test_user@test.com", "password": "another_password"}`
		reader := strings.NewReader(payload)

		resp, err := http.Post(gatewayAddr+"/api/v1/users", "application/json", reader)
		assert.NoError(t, err)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	cleanUserTable(t, db)
}
