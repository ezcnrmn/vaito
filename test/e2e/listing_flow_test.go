package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type LoginData struct {
	UserID int    `json:"userID"`
	Token  string `json:"token"`
}

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

type state struct {
	userData   LoginData
	adminData  LoginData
	categoryID int
	listingID  int
}

func TestListingFlow(t *testing.T) {
	client := &http.Client{}
	state := struct {
		userData   LoginData
		adminData  LoginData
		categoryID int
		listingID  int
	}{}

	t.Run("login as user and admin", func(t *testing.T) {
		// login as user
		userPayload := `{"password": "12345678", "email": "user@test.com"}`

		userResp := sendPost(t, client, gatewayAddr+"/api/v1/login", userPayload, "")
		defer userResp.Body.Close()

		assert.Equal(t, http.StatusOK, userResp.StatusCode)
		readJSON(t, userResp.Body, &state.userData)

		assert.Equal(t, state.userData.UserID, 2)
		assert.NotEqual(t, state.userData.Token, "")

		// login as admin
		adminPayload := `{"password": "12345678", "email": "admin@test.com"}`

		adminResp := sendPost(t, client, gatewayAddr+"/api/v1/login", adminPayload, "")
		defer adminResp.Body.Close()

		assert.Equal(t, http.StatusOK, adminResp.StatusCode)
		readJSON(t, adminResp.Body, &state.adminData)

		assert.Equal(t, state.adminData.UserID, 1)
		assert.NotEqual(t, state.adminData.Token, "")
	})

	t.Run("get categories", func(t *testing.T) {
		resp := sendGet(t, client, gatewayAddr+"/api/v1/categories", "")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

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

		resp := sendPost(t, client, gatewayAddr+"/api/v1/listings", listing, state.userData.Token)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

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

	t.Run("get created listing", func(t *testing.T) {
		publicResp := sendGet(t, client, fmt.Sprintf("%s/api/v1/listings/%d", gatewayAddr, state.listingID), "")
		defer publicResp.Body.Close()

		assert.Equal(t, http.StatusNotFound, publicResp.StatusCode, "listing must not be publicly visible")

		userResp := sendGet(t, client, fmt.Sprintf("%s/api/v1/users/%d/listings/%d", gatewayAddr, state.userData.UserID, state.listingID), state.userData.Token)
		defer userResp.Body.Close()

		assert.Equal(t, http.StatusOK, userResp.StatusCode, "listing must be accessible to user")

		var listingResp struct {
			Listing Listing `json:"listing"`
		}
		readJSON(t, userResp.Body, &listingResp)

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

	// 4. update listing
	// 5. get this listing again (both ways)
	// 6. send this listing to moderation
	// 7. get this listing by user and check status
	// 8. activate this listing as moderation
	// 9 get this listing as public
	// 10. deactivate listing as moderation
	// 11. get listing as public
	// 12. activate listing as user
	// 13. get listing as public
	// 14. deactivate listing as user
	// 15. get listing as public
	// 16. delete listing
	// 17. get listing
}
