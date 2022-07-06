package user_service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Youn-sik/KSCONNECT_ADMIN/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Alert struct {
	Id        string `json:"_id"`
	Title     string `json:"Title"`
	Timestamp string `json:"Timestamp"`
	Uid       string `json:"Uid"`
	Station   string `json:"Station"`
	Device    int    `json:"Device"`
	Usage     int    `json:"Usage,omitempty"`
	Payment   string `json:"Payment,omitempty"`
}

type ResultData struct {
	AlertInfo struct {
		Title     string
		Context   string
		Timestamp string
	}
	ETCInfo struct {
		StationName  string
		DeviceNumber int
		Usage        int
		Payment      string
	}
}

type ResultTransaction struct {
	StartTimestamp string      `json:"StartTimestamp,omitempty"`
	StopTimestamp  string      `json:"StopTimestamp,omitempty"`
	Transaction    Transaction `json:"Transaction,omitempty"`
	Id             int         `json:"_id,omitempty"`
}
type Transaction struct {
	Chargepointid string             `json:"chargepointid"`
	Payload       TransactionPayload `json:"payload"`
	Transactionid int                `json:"transactionid"`
}
type TransactionPayload struct {
	Connectorid     int    `json:"connectorid"`
	Idtag           string `json:"idtag"`
	MeterStop       int    `json:"meterStop"`
	Meterstart      int    `json:"meterstart"`
	Reason          string `json:"reason"`
	Reservationid   int    `json:"reservationid"`
	Timestamp       string `json:"timestamp"`
	Transactiondata string `json:"transactiondata"`
}

// type ResultMeterValues struct {
// 	Id          int        `json:"_id,omitempty"`
// 	MeterValues MeterValue `json:"MeterValues,omitempty"`
// 	Timestamp   string     `json:"Timestamp"`
// }
// type MeterValue struct {
// 	ChargePointId string            `json:"chargepointid"`
// 	Payload       MeterValuePayload `json:"payload"`
// }
// type MeterValuePayload struct {
// 	Connectorid     int    `json:"connectorid"`
// 	Idtag           string `json:"idtag"`
// 	MeterStop       int    `json:"meterStop"`
// 	Meterstart      int    `json:"meterstart"`
// 	Reason          string `json:"reason"`
// 	Reservationid   int    `json:"reservationid"`
// 	Timestamp       string `json:"timestamp"`
// 	Transactiondata string `json:"transactiondata"`
// }

func AlertList(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	var Uid struct {
		Uid string `json:"uid"`
	}
	err := c.Bind(&Uid)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	conn1 := database.NewMysqlConnection()

	client := database.NewMongodbConnection()
	conn := client.Database("Admin_Service").Collection("service_alert")
	filter := bson.M{"Uid": Uid.Uid}
	options := options.Find()
	options.SetSort(bson.M{"Timestamp": -1})
	options.SetLimit(20)
	cursor, err := conn.Find(context.TODO(), filter, options)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Query 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}
	var alertListArr []ResultData
	for cursor.Next(context.TODO()) {
		alert := Alert{}
		var elem bson.M
		if err := cursor.Decode(&elem); err != nil {
			log.Println(err)
			send_data.result = "false"
			send_data.errStr = "Query parsing 문제가 발생하였습니다."
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
			return
		}
		result, _ := json.Marshal(elem)
		json.Unmarshal(result, &alert)
		// log.Printf("%+v\n", alert)

		var StationName string
		var DeviceNumber int
		rows, err := conn1.Query("select name from charge_station where station_id = ?", alert.Station)
		if err != nil {
			log.Println(err)
			send_data.result = "false"
			send_data.errStr = "DB Query 중 문제가 발생하였습니다."
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
			return
		}
		for rows.Next() {
			err := rows.Scan(&StationName)
			if err != nil {
				log.Println(err)
				send_data.result = "false"
				send_data.errStr = "Query Parsing 중 문제가 발생하였습니다."
				c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
				return
			}
		}
		rows, err = conn1.Query("select device_number from charge_device where station_id = ? and device_number = ?", alert.Station, alert.Device)
		if err != nil {
			log.Println(err)
			send_data.result = "false"
			send_data.errStr = "DB Query 중 문제가 발생하였습니다."
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
			return
		}
		for rows.Next() {
			err := rows.Scan(&DeviceNumber)
			if err != nil {
				log.Println(err)
				send_data.result = "false"
				send_data.errStr = "Query Parsing 중 문제가 발생하였습니다."
				c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
				return
			}
		}

		resultData := ResultData{}
		resultData.AlertInfo.Title = alert.Title
		resultData.AlertInfo.Timestamp = alert.Timestamp
		resultData.ETCInfo.StationName = StationName
		resultData.ETCInfo.DeviceNumber = DeviceNumber
		resultData.ETCInfo.Usage = alert.Usage
		resultData.ETCInfo.Payment = alert.Payment

		DeviceNumberStr := strconv.Itoa(DeviceNumber)
		UsageStr := strconv.Itoa(alert.Usage)
		PaymentStr := alert.Payment

		switch alert.Title {
		case "충전 시작":
			resultData.AlertInfo.Context = StationName + "의 " + DeviceNumberStr + "번 충전기에서 충전 시작되었습니다."
		case "충전 완료":
			resultData.AlertInfo.Context = StationName + "의 " + DeviceNumberStr + "번 충전기에서 " + UsageStr + "kW 충전 완료되었습니다."
		case "결제 완료":
			resultData.AlertInfo.Context = StationName + "의 " + DeviceNumberStr + "번 충전기에서 " + PaymentStr + "원 결제 완료되었습니다."
		}
		// log.Printf("%+v\n", resultData)
		alertListArr = append(alertListArr, resultData)
	}

	send_data.result = "true"
	send_data.errStr = ""
	c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "alert_list": alertListArr})
}

func CheckChargeStatus(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	var Uid struct {
		Uid string `json:"uid"`
	}
	err := c.Bind(&Uid)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}

	// MongoDB의 Transaction 에서 확인
	transaction := ResultTransaction{}

	ntime := time.Now().Format(time.RFC3339)
	ntime = ntime[:10]
	log.Println(ntime)
	// ytime := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
	// ytime = ytime[:10]
	// log.Println(ytime)

	client := database.NewMongodbConnection()
	conn := client.Database("Admin_Service").Collection("ocpp_Transaction")
	filter := bson.M{"$and": []bson.M{{"StopTimestamp": nil}, {"Transaction.payload.idtag": Uid.Uid}}}
	options := options.Find()
	options.SetSort(bson.M{"StartTimestamp": -1})
	options.SetLimit(1)
	cursor, err := conn.Find(context.TODO(), filter, options)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Query parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}
	if cursor.Next(context.TODO()) {
		if err := cursor.Decode(&transaction); err != nil {
			log.Println(err)
		}
		// log.Printf("%+v\n", transaction)

		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "transaction": transaction})
		return
	} else {
		send_data.result = "false"
		send_data.errStr = "현재 충전중이 아닙니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}
}

func GetMeterValue(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	var Tid struct {
		Tid int `json:"tid"`
	}
	err := c.Bind(&Tid)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}

	client := database.NewMongodbConnection()
	conn := client.Database("Admin_Service").Collection("ocpp_MeterValues")
	filter := bson.M{"MeterValues.payload.transactionid": Tid.Tid}
	options := options.Find()
	options.SetSort(bson.M{"Timestamp": -1})
	options.SetLimit(1)
	cursor, err := conn.Find(context.TODO(), filter, options)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Query parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}
	if cursor.Next(context.TODO()) {
		var elem bson.M
		if err := cursor.Decode(&elem); err != nil {
			log.Println(err)
		}
		// log.Printf("%+v\n", elem)

		// if err := cursor.Decode(&transaction); err != nil {
		// 	log.Println(err)
		// }
		// // log.Printf("%+v\n", transaction)

		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "metervalue": elem})
		return
	} else {
		send_data.result = "false"
		send_data.errStr = "현재 충전중이 아닙니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}
}
