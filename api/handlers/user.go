package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/davidwarshaw/golang-user-crud/api/database"
	"github.com/davidwarshaw/golang-user-crud/api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/nyaruka/phonenumbers"
	"golang.org/x/crypto/bcrypt"
)

func normalizeIncomingUserAccount(userIncoming models.UserIncoming) (*models.UserAccount, error) {
	// Parse the phone number
	primaryPhoneNumber, err := phonenumbers.Parse(userIncoming.PrimaryPhoneNumber, "US")
	if err != nil {
		return &models.UserAccount{}, errors.New("primary_phone_number must be a valid US telephone number")
	}
	primaryPhoneNumberString := phonenumbers.Format(primaryPhoneNumber, phonenumbers.NATIONAL)

	// Hash the password
	password := []byte(userIncoming.Password)
	passwordHash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return &models.UserAccount{}, errors.New("error hasing password")
	}

	// Create the DB model from the API model
	userAccount := &models.UserAccount{
		UserBase:     userIncoming.UserBase,
		PasswordHash: string(passwordHash),
	}

	// Use the reformatted phone number
	userAccount.PrimaryPhoneNumber = primaryPhoneNumberString

	return userAccount, nil
}

// @Summary Retrieve all users
// @Accept  json
// @Produce  json
// @Param   page      	query	int	false  "default: 1"
// @Param   page_size   query	int	false  "default: 20"
// @Success 200 {array} models.UserOutgoing	"The user entities"
// @Router /users [get]
func RetrieveAllUsers(c *gin.Context) {
	db := c.MustGet("DB").(*pg.DB)

	// Get pagination
	var paginationIncoming models.Pagination
	if err := c.ShouldBindQuery(&paginationIncoming); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Set default pagination and offset
	if paginationIncoming.Page == 0 {
		paginationIncoming.Page = 1
	}
	if paginationIncoming.PageSize == 0 {
		paginationIncoming.PageSize = 20
	}
	offset := (paginationIncoming.Page - 1) * paginationIncoming.PageSize

	// Retrieve all the user accounts
	var userAccounts []models.UserAccount
	db.Model(&userAccounts).Limit(paginationIncoming.PageSize).Offset(offset).Select()

	// Transform models
	var usersOutgoing []models.UserOutgoing
	for _, userAccount := range userAccounts {
		userOutgoing := &models.UserOutgoing{
			UserID:   userAccount.UserID,
			UserBase: userAccount.UserBase,
		}
		usersOutgoing = append(usersOutgoing, *userOutgoing)
	}

	c.JSON(http.StatusOK, gin.H{"data": usersOutgoing, "pagination": paginationIncoming})
}

// @Summary Create a user
// @Accept  json
// @Produce  json
// @Param   user      	body	models.UserIncoming	true "The user data to be created"
// @Success 201 {body} models.UserOutgoing
// @Router /users [post]
func CreateUser(c *gin.Context) {
	db := c.MustGet("DB").(*pg.DB)

	// Get the request body
	var userIncoming models.UserIncoming
	if err := c.BindJSON(&userIncoming); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userAccount, err := normalizeIncomingUserAccount(userIncoming)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Save to the DB
	if _, err = db.Model(userAccount).Insert(); err != nil {
		c.Error(err)
		if strings.Contains(err.Error(), database.PK_ERROR_CODE) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "user_name already exists"})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		return
	}

	userOutgoing := &models.UserOutgoing{
		UserID:   userAccount.UserID,
		UserBase: userAccount.UserBase,
	}

	c.JSON(http.StatusCreated, userOutgoing)
}

// @Summary Retrieve a user by id
// @Produce  json
// @Param   id path int true "The id of the user to be retrieved"
// @Success 200 {object} models.UserOutgoing "The user entity for that id"
// @Router /users/:id [get]
func RetrieveUser(c *gin.Context) {
	db := c.MustGet("DB").(*pg.DB)

	// Get URL param
	var userId models.UserID
	if err := c.ShouldBindUri(&userId); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Retrieve all the user accounts
	var userAccount models.UserAccount
	userAccount.UserID = userId

	if err := db.Model(&userAccount).WherePK().Select(); err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"message": "User Account not found"})
		return
	}

	userOutgoing := &models.UserOutgoing{
		UserID:   userAccount.UserID,
		UserBase: userAccount.UserBase,
	}

	c.JSON(http.StatusOK, userOutgoing)
}

// @Summary Update a user by id
// @Accept  json
// @Produce  json
// @Param   id path int true "The id of the user to be updated"
// @Param   user      	body	models.UserIncoming	true "The user data to be updated"
// @Success 200 {object} models.UserOutgoing "The updated user entity for that id"
// @Router /users/:id [put]
func UpdateUser(c *gin.Context) {
	db := c.MustGet("DB").(*pg.DB)

	// Get URL param
	var userId models.UserID
	if err := c.ShouldBindUri(&userId); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Get the request body
	var userIncoming models.UserIncoming
	if err := c.BindJSON(&userIncoming); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userAccount, err := normalizeIncomingUserAccount(userIncoming)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// The URL ID overrides any model ID
	userAccount.UserID = userId

	if _, err := db.Model(userAccount).WherePK().Update(); err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"message": "User Account not found"})
		return
	}

	userOutgoing := &models.UserOutgoing{
		UserID:   userAccount.UserID,
		UserBase: userAccount.UserBase,
	}

	c.JSON(http.StatusOK, userOutgoing)
}

// @Summary Delete a user by id
// @Produce  json
// @Param   id path int true "The id of the user to be deleted"
// @Success 204 {string} nil
// @Router /users/:id [delete]
func DeleteUser(c *gin.Context) {
	db := c.MustGet("DB").(*pg.DB)

	// Get URL param
	var userId models.UserID
	if err := c.ShouldBindUri(&userId); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Retrieve all the user accounts
	var userAccount models.UserAccount
	userAccount.UserID = userId

	if _, err := db.Model(&userAccount).WherePK().Delete(); err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"message": "User Account not found"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
