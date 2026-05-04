package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Listing struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	UserID       int     `json:"user_id"`
	Status       string  `json:"status"`
	Price        int     `json:"price"`
	CreatedAt    *string `json:"created_at"`
	PublishedAt  *string `json:"published_at"`
}

func TestListingFlow(t *testing.T) {
	client := &http.Client{}
	state := struct {
		userData   LoginData
		adminData  LoginData
		categoryID int
		listingID  int
	}{}

	mustNotBePubliclyVisibleTest := func(t *testing.T) {
		resp, err := sendGet(t, client, fmt.Sprintf("%s/api/v1/listings/%d", gatewayAddr, state.listingID), "")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "listing must not be publicly visible")
	}

	t.Run("login as user", func(t *testing.T) {
		loginPayload := `{"password": "12345678", "email": "user@test.com"}`

		resp, err := send(t, client, "POST", gatewayAddr+"/api/v1/login", loginPayload, "")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		readJSON(t, resp.Body, &state.userData)

		assert.Equal(t, state.userData.UserID, 2)
		assert.NotEqual(t, state.userData.Token, "")
	})

	t.Run("login as admin", func(t *testing.T) {
		loginPayload := `{"password": "12345678", "email": "admin@test.com"}`

		resp, err := send(t, client, "POST", gatewayAddr+"/api/v1/login", loginPayload, "")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		readJSON(t, resp.Body, &state.adminData)

		assert.Equal(t, state.adminData.UserID, 1)
		assert.NotEqual(t, state.adminData.Token, "")
	})

	t.Run("get categories", func(t *testing.T) {
		resp, err := sendGet(t, client, gatewayAddr+"/api/v1/categories", "")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var categoriesResp struct {
			Categories []Category `json:"categories"`
		}

		readJSON(t, resp.Body, &categoriesResp)

		for _, c := range categoriesResp.Categories {
			if c.Name == "Electronics" {
				state.categoryID = c.ID
				break
			}
		}
		assert.NotZero(t, state.categoryID)
	})

	t.Run("create listing", func(t *testing.T) {
		listing := fmt.Sprintf(`{
    "title": "Продам iPhone 3G",
    "description": "Полный комплект: Зарядное устройство. Коробка. Наушники. Чехол в подарок.",
    "category_id": %d,
    "price": 1500
		}`, state.categoryID)

		resp, err := send(t, client, "POST", gatewayAddr+"/api/v1/listings", listing, state.userData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var listingResp struct {
			Listing Listing `json:"listing"`
		}

		readJSON(t, resp.Body, &listingResp)

		assert.NotZero(t, &listingResp.Listing.ID)
		assert.Equal(t, "Продам iPhone 3G", listingResp.Listing.Title)
		assert.Equal(t, "Полный комплект: Зарядное устройство. Коробка. Наушники. Чехол в подарок.", listingResp.Listing.Description)
		assert.Equal(t, state.categoryID, listingResp.Listing.CategoryID)
		assert.Equal(t, "Electronics", listingResp.Listing.CategoryName)
		assert.Equal(t, 2, listingResp.Listing.UserID)
		assert.Equal(t, "Draft", listingResp.Listing.Status)
		assert.Equal(t, 1500, listingResp.Listing.Price)
		assert.NotNil(t, listingResp.Listing.CreatedAt)
		assert.Nil(t, listingResp.Listing.PublishedAt)

		state.listingID = listingResp.Listing.ID
	})

	t.Run("get created listing: not visible for public", mustNotBePubliclyVisibleTest)

	t.Run("get created listing: visible for owner", func(t *testing.T) {
		resp, err := sendGet(t, client, fmt.Sprintf("%s/api/v1/users/%d/listings/%d", gatewayAddr, state.userData.UserID, state.listingID), state.userData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode, "listing must be accessible to user")

		var listingResp struct {
			Listing Listing `json:"listing"`
		}
		readJSON(t, resp.Body, &listingResp)

		assert.Equal(t, state.listingID, listingResp.Listing.ID)
		assert.Equal(t, "Продам iPhone 3G", listingResp.Listing.Title)
		assert.Equal(t, "Полный комплект: Зарядное устройство. Коробка. Наушники. Чехол в подарок.", listingResp.Listing.Description)
		assert.Equal(t, state.categoryID, listingResp.Listing.CategoryID)
		assert.Equal(t, "Electronics", listingResp.Listing.CategoryName)
		assert.Equal(t, 2, listingResp.Listing.UserID)
		assert.Equal(t, "Draft", listingResp.Listing.Status)
		assert.Equal(t, 1500, listingResp.Listing.Price)
		assert.NotNil(t, listingResp.Listing.CreatedAt)
		assert.Nil(t, listingResp.Listing.PublishedAt)
	})

	t.Run("update listing", func(t *testing.T) {
		listing := `{
    "title": "Продам iPhone 3GS (2011)",
    "description": "Полный комплект: Зарядное устройство. Коробка. Наушники. Чехол в подарок. UPD: Внешний вид с потёртостями",
    "price": 2500
		}`

		resp, err := send(t, client, "PATCH", fmt.Sprintf("%s/api/v1/listings/%d", gatewayAddr, state.listingID), listing, state.userData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var listingResp struct {
			Listing Listing `json:"listing"`
		}

		readJSON(t, resp.Body, &listingResp)

		assert.Equal(t, "Продам iPhone 3GS (2011)", listingResp.Listing.Title)
		assert.Equal(t, "Полный комплект: Зарядное устройство. Коробка. Наушники. Чехол в подарок. UPD: Внешний вид с потёртостями", listingResp.Listing.Description)
		assert.Equal(t, 2500, listingResp.Listing.Price)
	})

	t.Run("get updated listing: not visible for public", mustNotBePubliclyVisibleTest)

	t.Run("get updated listing: visible for owner", func(t *testing.T) {
		resp, err := sendGet(t, client, fmt.Sprintf("%s/api/v1/users/%d/listings/%d", gatewayAddr, state.userData.UserID, state.listingID), state.userData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode, "listing must be accessible to user")

		var listingResp struct {
			Listing Listing `json:"listing"`
		}
		readJSON(t, resp.Body, &listingResp)

		assert.Equal(t, "Продам iPhone 3GS (2011)", listingResp.Listing.Title)
		assert.Equal(t, "Полный комплект: Зарядное устройство. Коробка. Наушники. Чехол в подарок. UPD: Внешний вид с потёртостями", listingResp.Listing.Description)
		assert.Equal(t, 2500, listingResp.Listing.Price)
	})

	t.Run("send listing to moderation: success", func(t *testing.T) {
		resp, err := send(t, client, "POST", fmt.Sprintf("%s/api/v1/listings/%d/moderation", gatewayAddr, state.listingID), "", state.userData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "listing must be successfully send to moderation")
	})

	t.Run("send listing to moderation: error", func(t *testing.T) {
		resp, err := send(t, client, "POST", fmt.Sprintf("%s/api/v1/listings/%d/moderation", gatewayAddr, state.listingID), "", state.userData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status must be changed only once")
	})

	t.Run("get listing with moderation status", func(t *testing.T) {
		resp, err := sendGet(t, client, fmt.Sprintf("%s/api/v1/users/%d/listings/%d", gatewayAddr, state.userData.UserID, state.listingID), state.userData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var listingResp struct {
			Listing Listing `json:"listing"`
		}
		readJSON(t, resp.Body, &listingResp)

		assert.Equal(t, "Moderation", listingResp.Listing.Status)
	})

	t.Run("activate listing as moderation: success", func(t *testing.T) {
		resp, err := send(t, client, "POST", fmt.Sprintf("%s/api/v1/moderation/listings/%d/activate", gatewayAddr, state.listingID), "", state.adminData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "listing must be successfully activated")
	})

	t.Run("activate listing as moderation: error", func(t *testing.T) {
		resp, err := send(t, client, "POST", fmt.Sprintf("%s/api/v1/moderation/listings/%d/activate", gatewayAddr, state.listingID), "", state.adminData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status must be changed only once")
	})

	t.Run("get active public listing", func(t *testing.T) {
		resp, err := sendGet(t, client, fmt.Sprintf("%s/api/v1/listings/%d", gatewayAddr, state.listingID), "")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var listingResp struct {
			Listing Listing `json:"listing"`
		}
		readJSON(t, resp.Body, &listingResp)

		assert.Equal(t, "Active", listingResp.Listing.Status)

		assert.Equal(t, state.listingID, listingResp.Listing.ID)
		assert.Equal(t, "Продам iPhone 3GS (2011)", listingResp.Listing.Title)
		assert.Equal(t, "Полный комплект: Зарядное устройство. Коробка. Наушники. Чехол в подарок. UPD: Внешний вид с потёртостями", listingResp.Listing.Description)
		assert.Equal(t, state.categoryID, listingResp.Listing.CategoryID)
		assert.Equal(t, "Electronics", listingResp.Listing.CategoryName)
		assert.Equal(t, 2, listingResp.Listing.UserID)
		assert.Equal(t, 2500, listingResp.Listing.Price)
	})

	t.Run("deactivate listing as moderation: success", func(t *testing.T) {
		resp, err := send(t, client, "POST", fmt.Sprintf("%s/api/v1/moderation/listings/%d/deactivate", gatewayAddr, state.listingID), "", state.adminData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "listing must be successfully deactivated")
	})

	t.Run("deactivate listing as moderation: error", func(t *testing.T) {
		resp, err := send(t, client, "POST", fmt.Sprintf("%s/api/v1/moderation/listings/%d/deactivate", gatewayAddr, state.listingID), "", state.adminData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status must be changed only once")
	})

	t.Run("get deactivated listing: not visible for public", mustNotBePubliclyVisibleTest)

	t.Run("activate listing as user: success", func(t *testing.T) {
		resp, err := send(t, client, "POST", fmt.Sprintf("%s/api/v1/listings/%d/activate", gatewayAddr, state.listingID), "", state.userData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "listing must be successfully activated")
	})

	t.Run("activate listing as user: error", func(t *testing.T) {
		resp, err := send(t, client, "POST", fmt.Sprintf("%s/api/v1/listings/%d/activate", gatewayAddr, state.listingID), "", state.userData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status must be changed only once")
	})

	t.Run("get active listing (public api)", func(t *testing.T) {
		resp, err := sendGet(t, client, fmt.Sprintf("%s/api/v1/listings/%d", gatewayAddr, state.listingID), "")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var listingResp struct {
			Listing Listing `json:"listing"`
		}
		readJSON(t, resp.Body, &listingResp)

		assert.Equal(t, "Active", listingResp.Listing.Status)
	})

	t.Run("deactivate listing as user", func(t *testing.T) {
		resp, err := send(t, client, "POST", fmt.Sprintf("%s/api/v1/listings/%d/deactivate", gatewayAddr, state.listingID), "", state.userData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "listing must be successfully deactivated")
	})

	t.Run("deactivate listing as user", func(t *testing.T) {
		resp, err := send(t, client, "POST", fmt.Sprintf("%s/api/v1/listings/%d/deactivate", gatewayAddr, state.listingID), "", state.userData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status must be changed only once")
	})

	t.Run("get deactivated listing: not visible for public", mustNotBePubliclyVisibleTest)

	t.Run("delete listing", func(t *testing.T) {
		resp, err := send(t, client, "DELETE", fmt.Sprintf("%s/api/v1/listings/%d", gatewayAddr, state.listingID), "", state.userData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("get deleted listing: not accessible for public", mustNotBePubliclyVisibleTest)

	t.Run("get deleted listing: not accessible for owner", func(t *testing.T) {
		resp, err := sendGet(t, client, fmt.Sprintf("%s/api/v1/users/%d/listings/%d", gatewayAddr, state.userData.UserID, state.listingID), state.userData.Token)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "listing must not be accessible")
	})
}
