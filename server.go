package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Youn-sik/KSCONNECT_ADMIN/database"
	"github.com/Youn-sik/KSCONNECT_ADMIN/natsclient"
	n "github.com/Youn-sik/KSCONNECT_ADMIN/natsclient"
	"github.com/Youn-sik/KSCONNECT_ADMIN/router/admin/user"
	"github.com/Youn-sik/KSCONNECT_ADMIN/router/b2b/b2b_account"
	"github.com/Youn-sik/KSCONNECT_ADMIN/router/b2b/device"
	"github.com/Youn-sik/KSCONNECT_ADMIN/router/b2b/station"
	"github.com/Youn-sik/KSCONNECT_ADMIN/router/user/user_account"

	v16 "github.com/aliml92/ocpp/v16"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

type IdTag struct {
	IdTag       string    `bson:"idTag" json:"idTag"`
	ExpiryDate  time.Time `bson:"expiryDate,omitempty" json:"expiryDate,omitempty"`
	ParentIdTag string    `bson:"parentIdTag,omitempty" json:"parentIdTag,omitempty"`
	Blocked     bool      `bson:"blocked,omitempty" json:"blocked,omitempty"`
}

var nc *n.NatsClient

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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, DELETE, POST")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(CORSMiddleware())

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
	btb_service.POST("/charge_device_request", func(c *gin.Context) {
		device.DeviceRequest(c)
	})
	btb_service.GET("/charge_device_request_list", func(c *gin.Context) {
		device.DeviceRequestList(c)
	})
	btb_service.POST("/charge_device_request_submit", func(c *gin.Context) {
		device.DeviceRequestSubmit(c)
	})
	btb_service.GET("/charge_device_request_history", func(c *gin.Context) {
		device.DeviceRequestHistory(c)
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
	user_service.GET("/membership_card_list", func(c *gin.Context) {
		user_account.MembershipCardList(c)
	})
	user_service.POST("/membership_card_create", func(c *gin.Context) {
		user_account.MembershipCardCreate(c)
	})
	user_service.POST("/membership_card_update", func(c *gin.Context) {
		user_account.MembershipCardUpdate(c)
	})
	user_service.POST("/membership_card_delete", func(c *gin.Context) {
		user_account.MembershipCardDelete(c)
	})
	user_service.POST("/membership_card_request", func(c *gin.Context) {
		user_account.MembershipCardRequest(c)
	})
	user_service.GET("/membership_card_request_list", func(c *gin.Context) {
		user_account.MembershipCardRequestList(c)
	})
	user_service.POST("/membership_card_request_submit", func(c *gin.Context) {
		user_account.MembershipCardRequestSubmit(c)
	})
	user_service.GET("/membership_card_request_history", func(c *gin.Context) {
		user_account.MembershipCardRequestHistory(c)
	})
	user_service.GET("/inquiry_board_list", func(c *gin.Context) {
		user_account.InquiryBoardList(c)
	})
	user_service.POST("/inquiry_board_reply", func(c *gin.Context) {
		user_account.InquiryBoardReply(c)
	})

	return router
}

func ReplyNats(subject string) {
	// Subscribe
	wg := sync.WaitGroup{}
	wg.Add(1)

	conn1 := database.NewMysqlConnection()

	switch subject {
	case "ocpp/v16/chargepoints":
		{
			var send_data []string
			rows, err := conn1.Query("select station_id from charge_station")
			if err != nil {
				log.Println(err)
				wg.Done()
			} else {
				for rows.Next() {
					var chargePoints string
					err = rows.Scan(&chargePoints)
					if err != nil {
						log.Println(err)
						wg.Done()
					} else {
						send_data = append(send_data, chargePoints)
					}
				}
			}
			nc.Reply(subject, send_data, &wg)
		}
	case "ocpp/v16/idtags":
		{
			var ttts []IdTag
			var ttt IdTag
			rows, err := conn1.Query("select uid as idTag from user")
			if err != nil {
				log.Println(err)
				wg.Done()
			} else {
				// resultJson := jsonify.Jsonify(rows)
				// log.Println(reflect.TypeOf(resultJson[0]))
				// send_data = resultJson

				for rows.Next() {
					// var idtags int
					err = rows.Scan(&ttt.IdTag)
					if err != nil {
						log.Println(err)
						wg.Done()
					} else {
						ttts = append(ttts, ttt)
					}
				}

				// log.Println(ttts)
				nc.Reply(subject, ttts, &wg)
			}
		}
	}

	wg.Wait()
}

func SubscribeNats() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	ch1 := make(chan v16.GMeterValuesReq)
	err := n.Subscribe[v16.GMeterValuesReq](nc, "ocpp/v16/MeterValues", ch1)
	if err != nil {
		wg.Done()
		log.Println(err)
	}
	for m := range ch1 {
		natsclient.MeterValuesReq(m)
	}

	wg.Add(1)
	ch2 := make(chan v16.GBootNotificationReq)
	err = n.Subscribe[v16.GBootNotificationReq](nc, "ocpp/v16/BootNotification", ch2)
	if err != nil {
		wg.Done()
		log.Println(err)
	}
	for m := range ch2 {
		log.Println(m)
	}

	wg.Add(1)
	ch3 := make(chan v16.GStartTransactionReq)
	err = n.Subscribe[v16.GStartTransactionReq](nc, "ocpp/v16/StartTranscation", ch3)
	if err != nil {
		wg.Done()
		log.Println(err)
	}
	for m := range ch3 {
		log.Println(m)
	}

	wg.Add(1)
	ch4 := make(chan v16.GStopTransactionReq)
	err = n.Subscribe[v16.GStopTransactionReq](nc, "ocpp/v16/StartTranscation", ch4)
	if err != nil {
		wg.Done()
		log.Println(err)
	}
	for m := range ch4 {
		log.Println(m)
	}

	wg.Wait()
}

func main() {
	var port string = ":4001"

	nc = n.NewNatsClient()
	defer nc.Close()

	cr := cron.New()
	// 1분마다 단말 status(사용 중) 값 판단해서 충전중인 단말 MongoDB에서 MeterValue polling 해서 RDB charge_device의 usage update
	cr.AddFunc("* */1 * * * *", n.UpdateMeterValue)
	// 매달 RDB의 charge_device usage 0으로 초기화 시 MongoDB 에 저장 필요.
	// cr.AddFunc("* * * * */1 *")

	log.Println(nc)

	go ReplyNats("ocpp/v16/chargepoints")
	go ReplyNats("ocpp/v16/idtags")

	router := setupRouter()
	log.Println("[SERVER] => Backend Admin application is listening on port " + port)
	router.Run(port)
}
