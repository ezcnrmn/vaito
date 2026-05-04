package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type User struct {
	ID    int
	Name  string
	Email string
}

type LoginData struct {
	UserID int    `json:"userID"`
	Token  string `json:"token"`
}

func TestUserFlow(t *testing.T) {
	client := &http.Client{}
	state := struct {
		UserID int
		Token  string
	}{}

	t.Run("create user", func(t *testing.T) {
		user := `{"name": "Test", "email": "test_user@test.com", "password": "12345678"}`

		resp, err := send(t, client, "POST", gatewayAddr+"/api/v1/users", user, "")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var userResp struct {
			User User `json:"user"`
		}
		readJSON(t, resp.Body, &userResp)

		assert.NotZero(t, userResp.User.ID)
		assert.Equal(t, "Test", userResp.User.Name)
		assert.Equal(t, "test_user@test.com", userResp.User.Email)

		state.UserID = userResp.User.ID
	})

	t.Run("create existing user", func(t *testing.T) {
		user := `{"name": "AnotherTest", "email": "test_user@test.com", "password": "another_password"}`

		resp, err := send(t, client, "POST", gatewayAddr+"/api/v1/users", user, "")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("login as created user", func(t *testing.T) {
		loginPayload := `{
		"email": "test_user@test.com",
		"password": "12345678"
		}`

		resp, err := send(t, client, "POST", gatewayAddr+"/api/v1/login", loginPayload, "")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var loginData LoginData
		readJSON(t, resp.Body, &loginData)

		assert.Equal(t, state.UserID, loginData.UserID)
		assert.Len(t, loginData.Token, 44)

		state.Token = loginData.Token
	})

	t.Run("login with wrong password", func(t *testing.T) {
		loginPayload := `{
		"email": "test_user@test.com",
		"password": "wrong_password"
		}`

		resp, err := send(t, client, "POST", gatewayAddr+"/api/v1/login", loginPayload, "")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("login with wrong email", func(t *testing.T) {
		loginPayload := `{
		"email": "wrong_email@test.com",
		"password": "12345678"
		}`

		resp, err := send(t, client, "POST", gatewayAddr+"/api/v1/login", loginPayload, "")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("get user data", func(t *testing.T) {
		resp, err := sendGet(t, client, fmt.Sprintf("%s/api/v1/users/%d", gatewayAddr, state.UserID), "")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var userResp struct {
			User User `json:"user"`
		}
		readJSON(t, resp.Body, &userResp)

		assert.Equal(t, state.UserID, userResp.User.ID)
		assert.Equal(t, "Test", userResp.User.Name)
		assert.Equal(t, "test_user@test.com", userResp.User.Email)
	})

	t.Run("update user name", func(t *testing.T) {
		user := `{"name": "UpdatedTest"}`

		resp, err := send(t, client, "PATCH", fmt.Sprintf("%s/api/v1/users/%d", gatewayAddr, state.UserID), user, state.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var userResp struct {
			User User `json:"user"`
		}
		readJSON(t, resp.Body, &userResp)

		assert.Equal(t, state.UserID, userResp.User.ID)
		assert.Equal(t, "UpdatedTest", userResp.User.Name)
		assert.Equal(t, "test_user@test.com", userResp.User.Email)
	})

	t.Run("update user email", func(t *testing.T) {
		user := `{"email": "updated_test_user@test.com"}`

		resp, err := send(t, client, "PATCH", fmt.Sprintf("%s/api/v1/users/%d", gatewayAddr, state.UserID), user, state.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var userResp struct {
			User User `json:"user"`
		}
		readJSON(t, resp.Body, &userResp)

		assert.Equal(t, state.UserID, userResp.User.ID)
		assert.Equal(t, "UpdatedTest", userResp.User.Name)
		assert.Equal(t, "updated_test_user@test.com", userResp.User.Email)
	})

	t.Run("update user password", func(t *testing.T) {
		user := `{"password": "updated_password"}`

		resp, err := send(t, client, "PUT", fmt.Sprintf("%s/api/v1/users/%d/update-password", gatewayAddr, state.UserID), user, state.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("get updated user data", func(t *testing.T) {
		resp, err := sendGet(t, client, fmt.Sprintf("%s/api/v1/users/%d", gatewayAddr, state.UserID), "")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var userResp struct {
			User User `json:"user"`
		}
		readJSON(t, resp.Body, &userResp)

		assert.Equal(t, state.UserID, userResp.User.ID)
		assert.Equal(t, "UpdatedTest", userResp.User.Name)
		assert.Equal(t, "updated_test_user@test.com", userResp.User.Email)
	})

	t.Run("login with old credentials", func(t *testing.T) {
		loginPayload := `{
		"email": "test_user@test.com",
		"password": "12345678"
		}`

		resp, err := send(t, client, "POST", gatewayAddr+"/api/v1/login", loginPayload, "")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("login with new credentials", func(t *testing.T) {
		loginPayload := `{
		"email": "updated_test_user@test.com",
		"password": "updated_password"
		}`

		resp, err := send(t, client, "POST", gatewayAddr+"/api/v1/login", loginPayload, "")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var loginData LoginData
		readJSON(t, resp.Body, &loginData)

		assert.Equal(t, state.UserID, loginData.UserID)
		assert.Len(t, loginData.Token, 44)
	})

	cleanUserTable(t, db)
}
