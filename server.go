package main

import (
	"log"
	"net/http"

	"github.com/Youn-sik/KSCONNECT_ADMIN/router/b2b_account"
	"github.com/Youn-sik/KSCONNECT_ADMIN/router/device"
	"github.com/Youn-sik/KSCONNECT_ADMIN/router/station"
	"github.com/Youn-sik/KSCONNECT_ADMIN/router/user"
	"github.com/Youn-sik/KSCONNECT_ADMIN/router/user_account"

	"github.com/gin-gonic/gin"
)

func authenticateMiddleware(c *gin.Context) {
	authToken := c.Request.Header.Get("authorization")

	if authToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"result": "false", "errStr": "No Token"})
		c.Abort()
		return
	}

	// token 검증 로직
	isValid := user.TokenCheck(authToken)
	if isValid {
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"result": "false", "errStr": "Expired Token"})
		c.Abort()
		return
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	NAuth := router.Group("/NAuth")
	NAuth.POST("/user/login", func(c *gin.Context) {
		user.Login(c)
	})

	btb_service := router.Group("/btb_service")
	btb_service.Use(authenticateMiddleware)
	btb_service.GET("/charge_station_list", func(c *gin.Context) {
		station.StationList(c)
	})
	btb_service.POST("/charge_station_create", func(c *gin.Context) {
		station.StationCreate(c)
	})
	btb_service.POST("/charge_station_update", func(c *gin.Context) {
		station.StationUpdate(c)
	})
	btb_service.POST("/charge_station_delete", func(c *gin.Context) {
		station.StationDelete(c)
	})
	btb_service.POST("/charge_station_request", func(c *gin.Context) {
		station.StationRequest(c)
	})
	btb_service.GET("/charge_station_request_list", func(c *gin.Context) {
		station.StationRequestList(c)
	})
	btb_service.POST("/charge_station_request_submit", func(c *gin.Context) {
		station.StationRequestSubmit(c)
	})
	btb_service.GET("/charge_station_request_history", func(c *gin.Context) {
		station.StationRequestHistory(c)
	})

	btb_service.GET("/charge_device_list", func(c *gin.Context) {
		device.DeviceList(c)
	})
	btb_service.POST("/charge_device_create", func(c *gin.Context) {
		device.DeviceCreate(c)
	})
	btb_service.POST("/charge_device_update", func(c *gin.Context) {
		device.DeviceUpdate(c)
	})
	btb_service.POST("/charge_device_delete", func(c *gin.Context) {
		device.DeviceDelete(c)
	})

	btb_service.GET("/user_list", func(c *gin.Context) {
		b2b_account.UserList(c)
	})

	user_service := router.Group("/user_service")
	user_service.Use(authenticateMiddleware)
	user_service.GET("/user_list", func(c *gin.Context) {
		user_account.UserList(c)
	})
	user_service.POST("/user_create", func(c *gin.Context) {
		user_account.UserCreate(c)
	})
	user_service.POST("/user_update", func(c *gin.Context) {
		user_account.UserUpdate(c)
	})
	user_service.POST("/user_delete", func(c *gin.Context) {
		user_account.UserDelete(c)
	})

	return router
}

func main() {
	var port string = ":4001"

	router := setupRouter()
	router.Run(port)
	log.Println("[SERVER] => Backend Admin application is listening on port " + port)
}
