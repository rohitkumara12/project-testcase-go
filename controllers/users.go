package controllers

import (
	"log"
	"net/http"
	"strconv"
	"test-case/models"
	"test-case/views"

	"github.com/gin-gonic/gin"
)

func DetailsUsers(c *gin.Context) {
	var responsemodel views.Baseresponse
	var responsedata models.DataNewUser

	userid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Error = append(responsemodel.Error, err.Error())
		c.JSON(http.StatusBadRequest, responsedata)
		return
	}
	result := models.DB.First(&responsedata, userid)
	if result.Error != nil {
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Error = append(responsemodel.Error, result.Error.Error())
		c.JSON(http.StatusBadRequest, responsedata)
		return
	}
	responsemodel.Status = http.StatusOK
	responsemodel.Data = responsedata

	c.JSON(http.StatusOK, responsedata)
}

func GetUsers(c *gin.Context) {

	var responsemodel views.Baseresponse
	var user models.DataNewUser
	var Addresses models.DataNewUser

	iduser, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Error = append(responsemodel.Error, "invalid user id")
		c.JSON(http.StatusBadRequest, responsemodel)
		return
	}
	result := models.DB.First(&user, iduser)
	if result.Error != nil {
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Error = append(responsemodel.Error, "user not found")
		c.JSON(http.StatusBadRequest, responsemodel)
		return
	}
	addressesresult := models.DB.Where("iduser = ?", iduser).Find(&Addresses)
	if addressesresult.Error != nil {
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Error = append(responsemodel.Error, "failed to retrive address user")
		c.JSON(http.StatusBadRequest, responsemodel)
		return
	}
	responsedata := map[string]interface{}{
		"id":     user.ID,
		"nama":   user.Name,
		"email":  user.Email,
		"alamat": Addresses,
	}
	responsemodel.Status = http.StatusOK
	responsemodel.Data = responsedata
	c.JSON(http.StatusOK, responsedata)

	log.Println("iduser:", iduser)
}

func Index(c *gin.Context) {
	var responsemodel views.Baseresponse
	var responsdata views.Userslist

	result := models.DB.Find(&responsdata.Userlist)
	if result.Error != nil {
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Error = append(responsemodel.Error, result.Error.Error())
		c.JSON(http.StatusBadRequest, responsemodel)
		return
	}
	responsemodel.Status = http.StatusOK
	responsemodel.Data = responsdata

	c.JSON(http.StatusOK, responsdata)

}
