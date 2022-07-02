package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Youn-sik/KSCONNECT_ADMIN/database"
	"github.com/Youn-sik/KSCONNECT_ADMIN/setting"
	v16 "github.com/aliml92/ocpp/v16"
	"github.com/bdwilliams/go-jsonify/jsonify"
	"go.mongodb.org/mongo-driver/bson"

	// "github.com/Youn-sik/KSCONNECT_ADMIN/natsclient"
	n "github.com/Youn-sik/KSCONNECT_ADMIN/natsclient"
	"github.com/Youn-sik/KSCONNECT_ADMIN/router/admin/user"
	"github.com/Youn-sik/KSCONNECT_ADMIN/router/b2b/b2b_account"
	"github.com/Youn-sik/KSCONNECT_ADMIN/router/b2b/device"
	"github.com/Youn-sik/KSCONNECT_ADMIN/router/b2b/report"
	"github.com/Youn-sik/KSCONNECT_ADMIN/router/b2b/station"
	"github.com/Youn-sik/KSCONNECT_ADMIN/router/user/user_account"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

type IdTag struct {
	IdTag       string    `bson:"idTag" json:"idTag"`
	ExpiryDate  time.Time `bson:"expiryDate,omitempty" json:"expiryDate,omitempty"`
	ParentIdTag string    `bson:"parentIdTag,omitempty" json:"parentIdTag,omitempty"`
	Blocked     bool      `bson:"blocked,omitempty" json:"blocked,omitempty"`
}

type GenMeterValuesRes struct {
	ChargePointId string             `json:"id"`
	Payload       v16.MeterValuesReq `json:"payload,omitempty"`
}
type GenBootNotificationReq struct {
	ChargePointId string                  `json:"id"`
	Payload       v16.BootNotificationReq `json:"payload,omitempty"`
}
type GenStartTransactionReq struct {
	ChargePointId string                  `json:"id"`
	TransactionId int                     `json:"transactionId"`
	Payload       v16.StartTransactionReq `json:"payload,omitempty"`
}
type GenStopTransactionReq struct {
	ChargePointId string                 `json:"id"`
	Payload       v16.StopTransactionReq `json:"payload,omitempty"`
}
type GenStatusNotificationReq struct {
	ChargePointId string                    `json:"id"`
	Payload       v16.StatusNotificationReq `json:"payload,omitempty"`
}

type ResultTransaction struct {
	StartTimestamp string      `json:"StartTimestamp"`
	StopTimestamp  string      `json:"StopTimestamp"`
	Transaction    Transaction `json:"Transaction"`
	Id             int         `json:"_id"`
}
type Transaction struct {
	Chargepointid string  `json:"chargepointid"`
	Payload       Payload `json:"payload"`
	Transactionid int     `json:"transactionid"`
}
type Payload struct {
	Connectorid     int    `json:"connectorid"`
	Idtag           string `json:"idtag"`
	MeterStop       int    `json:"meterStop"`
	Meterstart      int    `json:"meterstart"`
	Reason          string `json:"reason"`
	Reservationid   int    `json:"reservationid"`
	Timestamp       string `json:"timestamp"`
	Transactiondata string `json:"transactiondata"`
}

var nc *n.NatsClient

func authenticateMiddleware(c *gin.Context) {
	authToken := c.Request.Header.Get("authorization")
	authToken = strings.Replace(authToken, "Bearer ", "", 1)

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

func randomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(CORSMiddleware())

	NAuth := router.Group("/NAuth")
	NAuth.POST("/user/login", func(c *gin.Context) {
		user.Login(c)
	})

	NAuth.POST("/user/auth", func(c *gin.Context) {
		user.Auth(c)
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

	btb_service.GET("/billing_list_company", func(c *gin.Context) {
		report.ReportList(c)
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

func SubscribeNats(subject string) {
	wg := sync.WaitGroup{}

	switch subject {
	case "ocpp/v16/MeterValues":
		{
			wg.Add(1)
			ch := make(chan GenMeterValuesRes)
			err := n.Subscribe[GenMeterValuesRes](nc, subject, ch)
			if err != nil {
				log.Println(err)
				wg.Done()
			}
			m := <-ch
			log.Println("===" + subject + "===")
			log.Println(m)

			// MongoDB Save
			ntime := time.Now().Format(time.RFC3339)
			ntime = ntime[:19]
			client := database.NewMongodbConnection()
			conn := client.Database("Admin_Service").Collection("ocpp_MeterValues")
			result, err := conn.InsertOne(context.TODO(), bson.D{
				{Key: "MeterValues", Value: m},
				{Key: "Timestamp", Value: ntime},
			})
			if err != nil {
				log.Println(err)
				wg.Done()
			} else {
				log.Println(result)
				wg.Done()
			}
		}
	case "ocpp/v16/BootNotification":
		{
			wg.Add(1)
			ch := make(chan GenBootNotificationReq)
			err := n.Subscribe[GenBootNotificationReq](nc, subject, ch)
			if err != nil {
				log.Println(err)
				wg.Done()
			}
			b := <-ch
			log.Println("===" + subject + "===")
			log.Println(b)
			// MongoDB Save
			ntime := time.Now().Format(time.RFC3339)
			ntime = ntime[:19]
			client := database.NewMongodbConnection()
			conn := client.Database("Admin_Service").Collection("ocpp_BootNotification")
			result, err := conn.InsertOne(context.TODO(), bson.D{
				{Key: "BootNotification", Value: b},
				{Key: "Timestamp", Value: ntime},
			})
			if err != nil {
				log.Println(err)
				wg.Done()
			} else {
				log.Println(result)
				wg.Done()
			}

			// MYSQL station status Y ?
			// 처리 어떻게 할지 생각 해 보기..
		}
	case "ocpp/v16/StartTransaction":
		{
			wg.Add(1)
			ch := make(chan GenStartTransactionReq)
			err := n.Subscribe[GenStartTransactionReq](nc, subject, ch)
			if err != nil {
				log.Println(err)
				wg.Done()
			}
			s := <-ch
			log.Println("===" + subject + "===")
			log.Println(s)

			// FCM PUSH
			go fcm_push("충전이 시작되었습니다.")

			// MongoDB Save
			ntime := time.Now().Format(time.RFC3339)
			ntime = ntime[:19]
			client := database.NewMongodbConnection()
			conn := client.Database("Admin_Service").Collection("ocpp_Transaction")
			_, err = conn.InsertOne(context.TODO(), bson.D{
				{Key: "Transaction", Value: s},
				{Key: "StartTimestamp", Value: ntime},
			})
			if err != nil {
				log.Println(err)
				wg.Done()
			}

			// MYSQL device status I
			// 문제 있음 => device_id 는 charge_device 의 PK 이고, connector id는 chargePoint 1개당 1부터 늘어나는 숫자임.
			conn1 := database.NewMysqlConnection()
			_, err = conn1.Query("update charge_device set status = 'I' where station_id = ? and device_number = ?", s.ChargePointId, s.Payload.ConnectorId)
			if err != nil {
				log.Println(err)
				wg.Done()
			}

			wg.Done()
		}
	case "ocpp/v16/StopTransaction":
		{
			wg.Add(1)
			ch := make(chan GenStopTransactionReq)
			err := n.Subscribe[GenStopTransactionReq](nc, subject, ch)
			if err != nil {
				log.Println(err)
				wg.Done()
			}
			s := <-ch
			log.Println("===" + subject + "===")
			log.Println(s)

			// FCM PUSH
			go fcm_push("충전이 완료되었습니다.")

			// MongoDB Save
			ntime := time.Now().Format(time.RFC3339)
			ntime = ntime[:19]

			client := database.NewMongodbConnection()
			conn := client.Database("Admin_Service").Collection("ocpp_Transaction")

			// transactionId := strconv.Itoa(*s.Payload.TransactionId)
			// log.Println(transactionId)
			filter := bson.M{"Transaction.transactionid": *s.Payload.TransactionId}
			value := bson.D{{"$set", bson.D{{"Transaction.payload.meterStop", s.Payload.MeterStop}, {"Transaction.payload.reason", s.Payload.Reason}, {"Transaction.payload.transactiondata", s.Payload.TransactionData}, {"StopTimestamp", ntime}}}}

			_, err = conn.UpdateOne(context.TODO(), filter, value)
			if err != nil {
				log.Println(err)
				wg.Done()
			}

			var transaction ResultTransaction
			filter = bson.M{"Transaction.transactionid": *s.Payload.TransactionId}
			cursor, err := conn.Find(context.TODO(), filter)
			if err != nil {
				log.Println(err)
				wg.Done()
			} else {
				for cursor.Next(context.TODO()) {
					var elem bson.M
					if err := cursor.Decode(&elem); err != nil {
						log.Println(err)
						wg.Done()
					}
					result, _ := json.Marshal(elem)
					json.Unmarshal(result, &transaction)
				}
			}

			// MYSQL device status Y
			// MYSQL device usage value add
			conn1 := database.NewMysqlConnection()
			defer conn1.Close()

			// 문제 있음 => device_id 는 charge_device 의 PK 이고, connector id는 chargePoint 1개당 1부터 늘어나는 숫자임.
			_, err = conn1.Query("update charge_device set status = 'Y' where station_id = ? and device_number = ?", transaction.Transaction.Chargepointid, transaction.Transaction.Payload.Connectorid)
			if err != nil {
				log.Println(err)
				wg.Done()
			}
			// 문제 있음 => device_id 는 charge_device 의 PK 이고, connector id는 chargePoint 1개당 1부터 늘어나는 숫자임.
			rows, err := conn1.Query("select charge_device.usage from charge_device where station_id = ? and device_number = ?", transaction.Transaction.Chargepointid, transaction.Transaction.Payload.Connectorid)
			if err != nil {
				log.Println(err)
				wg.Done()
			}
			for rows.Next() {
				var usage int
				err = rows.Scan(&usage)
				if err != nil {
					log.Println(err)
					wg.Done()
				} else {
					totalUsage := usage + (transaction.Transaction.Payload.MeterStop - transaction.Transaction.Payload.Meterstart)
					rows, err = conn1.Query("update charge_device set charge_device.usage = ? where station_id = ? and device_number = ?", totalUsage, transaction.Transaction.Chargepointid, transaction.Transaction.Payload.Connectorid)
					if err != nil {
						log.Println(err)
						wg.Done()
					}
				}
			}

			// User Service Payment Request(getBillingKey)
			config, err := setting.LoadConfigSettingJSON()
			if err != nil {
				log.Fatal(err)
				panic(err)
			}
			postBody, _ := json.Marshal(map[string]string{
				"cardNumber":          "5327501015763628",
				"cardExpirationYear":  "27",
				"cardExpirationMonth": "01",
				"cardPassword":        "05",
			})
			responseBody := bytes.NewBuffer(postBody)
			resp, err := http.Post("http://"+config.User_service.Host+":"+config.User_service.Port+"/payment/billingkey", "application/json", responseBody)
			if err != nil {
				log.Println(err)
				wg.Done()
			}
			var billingKey struct {
				BillingKey string `json:"billingKey"`
			}
			body, _ := ioutil.ReadAll(resp.Body)
			json.Unmarshal(body, &billingKey)

			// User Service Payment Request(pay)
			// getChargeFee
			postBody, _ = json.Marshal(map[string]string{
				"station_id": transaction.Transaction.Chargepointid,
			})
			responseBody = bytes.NewBuffer(postBody)
			resp, err = http.Post("http://"+config.User_service.Host+":"+config.User_service.Port+"/charge_station/charge_price", "application/json", responseBody)
			if err != nil {
				log.Println(err)
				wg.Done()
			}
			var price struct {
				Price float32 `json:"charge_price"`
			}
			body, _ = ioutil.ReadAll(resp.Body)
			json.Unmarshal(body, &price)

			// amount := (price.Price * float32(transaction.Transaction.Payload.MeterStop-transaction.Transaction.Payload.Meterstart))
			// amountStr := fmt.Sprintf("%f", amount)
			// 테스트에 0 값이여서 실데이터 넣으면 필수 파라메터 누락이라고 뜸
			// orderID 값 랜덤 값으로
			orderId := randomString(8)
			postBody, _ = json.Marshal(map[string]string{
				"billingKey": billingKey.BillingKey,
				// "amount":     amountStr,
				"amount":    "1000",
				"orderID":   string(orderId),
				"orderName": "전기차 충전 요금 정산",
			})
			responseBody = bytes.NewBuffer(postBody)
			resp, err = http.Post("http://"+config.User_service.Host+":"+config.User_service.Port+"/payment/pay", "application/json", responseBody)
			if err != nil {
				log.Println(err)
				wg.Done()
			}
			respBody, err := ioutil.ReadAll(resp.Body)
			// log.Println(string(respBody))

			// 결제 내역 (유저 정보, 빌링 키, 결제 요청 리턴에 대해 MongoDB에 저장)
			rows, err = conn1.Query("select * from user where uid = ?", transaction.Transaction.Payload.Idtag)
			if err != nil {
				log.Println(err)
				wg.Done()
			}
			userInfo := jsonify.Jsonify(rows)

			conn = client.Database("Admin_Service").Collection("service_payment")
			_, err = conn.InsertOne(context.TODO(), bson.D{
				{Key: "Payment", Value: string(respBody)},
				{Key: "User", Value: userInfo[0]},
			})
			if err != nil {
				log.Println(err)
				wg.Done()
			}

			wg.Done()
		}
	case "ocpp/v16/StatusNotification":
		{
			wg.Add(1)
			ch := make(chan GenStatusNotificationReq)
			err := n.Subscribe[GenStatusNotificationReq](nc, subject, ch)
			if err != nil {
				log.Println(err)
				wg.Done()
			}
			s := <-ch
			log.Println("===" + subject + "===")
			log.Println(s)
			// MongoDB Save
			ntime := time.Now().Format(time.RFC3339)
			ntime = ntime[:19]
			client := database.NewMongodbConnection()
			conn := client.Database("Admin_Service").Collection("ocpp_StatusNotification")
			result, err := conn.InsertOne(context.TODO(), bson.D{
				{Key: "StatusNotification", Value: s},
				{Key: "Timestamp", Value: ntime},
			})
			if err != nil {
				log.Println(err)
				wg.Done()
			} else {
				log.Println(result)
				// MYSQL device status set
				wg.Done()
			}
		}
	}
	wg.Wait()
}

func fcm_push(title string) {
	config, err := setting.LoadConfigSettingJSON()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	postBody, _ := json.Marshal(map[string]string{
		"title": title,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://"+config.User_service.Host+":"+config.User_service.Port+"/fcm/push", "application/json", responseBody)
	if err != nil {
		log.Println(err)
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes)
	// 코드 수정 필요
	if str != `{"result":"true"}` {
		log.Println("****FCM PUSH ERROR****")
	}
}

func main() {
	var port string = ":4001"

	nc = n.NewNatsClient()
	defer nc.Close()

	cr := cron.New()
	// 1분마다 단말 status(사용 중) 값 판단해서 충전중인 단말 MongoDB에서 MeterValue polling 해서 RDB charge_device의 usage update
	_ = cr.AddFunc("1 * * * * *", n.UpdateMeterValue)
	// _ = cr.AddFunc("*/10 * * * * *", n.UpdateMeterValue)
	// 매달 RDB의 charge_device usage 0으로 초기화 시 MongoDB 에 저장 필요.
	// cr.AddFunc("* * * * */1 *")
	cr.Start()

	go ReplyNats("ocpp/v16/chargepoints")
	go ReplyNats("ocpp/v16/idtags")

	go SubscribeNats("ocpp/v16/MeterValues")
	go SubscribeNats("ocpp/v16/BootNotification")
	go SubscribeNats("ocpp/v16/StartTransaction")
	go SubscribeNats("ocpp/v16/StopTransaction")
	go SubscribeNats("ocpp/v16/StatusNotification")

	router := setupRouter()
	log.Println("[SERVER] => Backend Admin application is listening on port " + port)
	router.Run(port)
}
