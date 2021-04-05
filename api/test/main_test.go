package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidwarshaw/golang-user-crud/api/models"
	"github.com/davidwarshaw/golang-user-crud/api/server"
	"github.com/stretchr/testify/assert"
)

type UserAccounts struct {
	Data       []models.UserAccount `json:"data"`
	Pagination models.Pagination    `json:"pagination"`
}

func retrieveAllUsers(ts *httptest.Server, t *testing.T, pageSizeString string) UserAccounts {
	response, _ := http.Get(fmt.Sprintf("%s/users%s", ts.URL, pageSizeString))
	defer response.Body.Close()
	assert.Equal(t, response.StatusCode, 200, "Response should be OK")

	var userAccounts UserAccounts
	json.NewDecoder(response.Body).Decode(&userAccounts)

	return userAccounts
}

func retrieveUser(ts *httptest.Server, t *testing.T, id uint) models.UserOutgoing {
	response, _ := http.Get(fmt.Sprintf(fmt.Sprintf("%s/users/%d", ts.URL, id)))
	defer response.Body.Close()
	assert.Equal(t, response.StatusCode, 200, "Response should be OK")

	var userAccount models.UserOutgoing
	json.NewDecoder(response.Body).Decode(&userAccount)

	return userAccount
}

func createUser(ts *httptest.Server, t *testing.T, userJson []byte, expectedStatus int, expectedResponse string) models.UserOutgoing {
	response, _ := http.Post(fmt.Sprintf("%s/users", ts.URL), "application/json", bytes.NewReader(userJson))
	defer response.Body.Close()
	assert.Equal(t, response.StatusCode, expectedStatus, expectedResponse)

	var createdUser models.UserOutgoing
	json.NewDecoder(response.Body).Decode(&createdUser)

	return createdUser
}

func updateUser(ts *httptest.Server, t *testing.T, id uint, userJson []byte) models.UserOutgoing {
	request, _ := http.NewRequest("PUT", fmt.Sprintf("%s/users/%d", ts.URL, id), bytes.NewReader(userJson))
	client := &http.Client{}
	response, _ := client.Do(request)
	defer response.Body.Close()
	assert.Equal(t, response.StatusCode, 200, "Response should be OK")

	var createdUser models.UserOutgoing
	json.NewDecoder(response.Body).Decode(&createdUser)

	return createdUser
}

func deleteUser(ts *httptest.Server, t *testing.T, id uint) {
	request, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%d", ts.URL, id), nil)
	client := &http.Client{}
	response, _ := client.Do(request)
	assert.Equal(t, response.StatusCode, 204, "Response should be NO_CONTENT")
}

func TestUserRoute(t *testing.T) {
	// Create server
	ts := httptest.NewServer(server.Setup())
	defer ts.Close()

	// Read fixtures
	goodUser1Json, err := ioutil.ReadFile("fixtures/goodUser1.json")
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	goodUser2Json, err := ioutil.ReadFile("fixtures/goodUser2.json")
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	badUser3Json, err := ioutil.ReadFile("fixtures/badUser3.json")
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	// No users to start with
	userAccounts := retrieveAllUsers(ts, t, "")
	assert.Equal(t, len(userAccounts.Data), 0, "There should be no initial users")

	// Add some users
	newUser1 := createUser(ts, t, goodUser1Json, 201, "Response should be CREATED")
	assert.Equal(t, newUser1.UserName, "user1", "User Name should match")
	assert.Equal(t, newUser1.PrimaryPhoneNumber, "(555) 555-1234", "Primary Phone Number should be formatted")
	assert.Greater(t, newUser1.Id, uint(0), "Id should be set by DB (greater than 0)")

	newUser2 := createUser(ts, t, goodUser2Json, 201, "Response should be CREATED")
	assert.Equal(t, newUser2.UserName, "user2", "User Name should match")
	assert.Greater(t, newUser2.Id, uint(0), "Id should be set by DB (greater than 0)")

	// Users with bad data
	var badUser models.UserIncoming
	var jsonData []byte
	json.Unmarshal(badUser3Json, &badUser)
	badUser.UserName = "user1" // Duplicate username
	jsonData, _ = json.Marshal(badUser)
	createUser(ts, t, jsonData, 400, "Response should be BAD_REQUEST")

	json.Unmarshal(badUser3Json, &badUser)
	badUser.Password = "a" // Password too short
	jsonData, _ = json.Marshal(badUser)
	createUser(ts, t, jsonData, 400, "Response should be BAD_REQUEST")

	json.Unmarshal(badUser3Json, &badUser)
	badUser.Email = "abc" // Bad email
	jsonData, _ = json.Marshal(badUser)
	createUser(ts, t, jsonData, 400, "Response should be BAD_REQUEST")

	json.Unmarshal(badUser3Json, &badUser)
	badUser.PrimaryPhoneNumber = "abc" // Bad phone number
	jsonData, _ = json.Marshal(badUser)
	createUser(ts, t, jsonData, 400, "Response should be BAD_REQUEST")

	// Pagination
	userAccounts = retrieveAllUsers(ts, t, "?page=1&page_size=1")
	assert.Equal(t, len(userAccounts.Data), 1, "There should one user per page")
	firstPageId := userAccounts.Data[0].Id
	userAccounts = retrieveAllUsers(ts, t, "?page=2&page_size=1")
	assert.Equal(t, len(userAccounts.Data), 1, "There should one user per page")
	secondPageId := userAccounts.Data[0].Id
	// Ids on first and second page should be different
	assert.NotEqual(t, firstPageId, secondPageId)

	// Retrieve and update a user
	firstPageUser := retrieveUser(ts, t, firstPageId)
	firstPageUserUpdate := &models.UserIncoming{
		UserBase: firstPageUser.UserBase,
		Password: "anewpassword",
	}
	firstPageUserUpdate.FirstName = "a new name"
	jsonData, _ = json.Marshal(firstPageUserUpdate)
	firstUserUpdated := updateUser(ts, t, firstPageId, jsonData)
	assert.Equal(t, firstUserUpdated.FirstName, "a new name")

	// Delete all users
	userAccounts = retrieveAllUsers(ts, t, "")
	for _, userAccount := range userAccounts.Data {
		deleteUser(ts, t, userAccount.Id)
	}

	// All the users were deleted, so there should be no more
	userAccounts = retrieveAllUsers(ts, t, "")
	assert.Equal(t, len(userAccounts.Data), 0, "There should be no users after deleting them all")
}
