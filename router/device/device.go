package device

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Youn-sik/KSCONNECT_ADMIN/database"
	"github.com/Youn-sik/KSCONNECT_ADMIN/setting"
	"github.com/bdwilliams/go-jsonify/jsonify"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type CreateReq struct {
	// Device_id int
	Station_id  int
	Name        string
	Sirial      string
	Charge_type string
	Charge_way  string
	Available   string
	Status      string
	// Last_state string
	Device_number int
	Reg_date      string
	// Up_date       string
}

type UpdateReq struct {
	Device_id int
	// Station_id  int
	Name        string
	Sirial      string
	Charge_type string
	Charge_way  string
	Available   string
	Status      string
	// Last_state string
	Device_number int
	// Reg_date      string
	Up_date string
}

type DeleteReq struct {
	Device_id int
}

type RequestReq struct {
	// Request_id int
	Request_uid    int
	Station_id     int
	Name           string
	Sirial         string
	Charge_type    string
	Charge_way     string
	Available      string
	Status         string
	Device_number  int
	Request_status string
	// Request_value  string // wating
}

type RequestSubmitReq struct {
	Request_id     int    `json:"request_id"`
	Request_uid    int    `json:"request_uid"`
	Station_id     int    `json:"station_id"`
	Name           string `json:"name"`
	Sirial         string `json:"sirial"`
	Charge_type    string `json:"charge_type"`
	Charge_way     string `json:"charge_way"`
	Available      string `json:"available"`
	Status         string `json:"status"`
	Device_number  int    `json:"device_number"`
	Request_status string `json:"request_status"`
	Request_value  string `json:"request_value"` // permitted or reject
}

func DeviceList(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select * from charge_device")
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		resultJson := jsonify.Jsonify(rows)
		log.Println(resultJson)

		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "charge_devices": resultJson})
	}
}

func DeviceCreate(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	reqData := CreateReq{}
	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("insert into charge_device (station_id, name, sirial, charge_type, charge_way, available, status, device_number, reg_date, up_date) value (?,?,?,?,?,?,?,?,?,?)",
		reqData.Station_id, reqData.Name, reqData.Sirial, reqData.Charge_type, reqData.Charge_way, reqData.Available, reqData.Status, reqData.Device_number, reqData.Reg_date, reqData.Reg_date)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		log.Println(rows)
		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}
}

func DeviceUpdate(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	reqData := UpdateReq{}
	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("update charge_device set name = ?, sirial = ?,  charge_type = ?, charge_way = ?, available = ?, status = ?, device_number = ?, up_date = ? where device_id = ?",
		reqData.Name, reqData.Sirial, reqData.Charge_type, reqData.Charge_way, reqData.Available, reqData.Status, reqData.Device_number, reqData.Up_date, reqData.Device_id)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		log.Println(rows)
		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}
}

func DeviceDelete(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	reqData := DeleteReq{}
	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query(
		"delete from charge_device where device_id = ?", reqData.Device_id)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		log.Println(rows)
		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}
}

func DeviceRequest(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	reqData := RequestReq{}
	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	// MongoDB logging
	ntime := time.Now().Format(time.RFC3339)
	ntime = ntime[:19]
	client := database.NewMongodbConnection()
	conn := client.Database("Admin_Service").Collection("request_charge_device")

	result, err := conn.InsertOne(context.TODO(), bson.D{
		{Key: "Request_uid", Value: reqData.Request_uid},
		{Key: "Station_id", Value: reqData.Station_id},
		{Key: "Name", Value: reqData.Name},
		{Key: "Sirial", Value: reqData.Sirial},
		{Key: "Charge_type", Value: reqData.Charge_type},
		{Key: "Charge_way", Value: reqData.Charge_way},
		{Key: "Available", Value: reqData.Available},
		{Key: "Status", Value: reqData.Status},
		{Key: "Device_number", Value: reqData.Device_number},
		{Key: "Request_value", Value: "wating"},
		{Key: "Request_status", Value: reqData.Request_status},
		{Key: "Timestamp", Value: ntime},
	})
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "MongoDB logging 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		log.Println(result)

		// MYSQL input
		conn1 := database.NewMysqlConnection()
		defer conn1.Close()

		rows, err := conn1.Query("insert into request_charge_device (request_uid, station_id, name, sirial, charge_type, charge_way, available, status, "+
			"device_number, request_time, request_status) value (?,?,?,?,?,?,?,?,?,?,?)",
			reqData.Request_uid, reqData.Station_id, reqData.Name, reqData.Sirial, reqData.Charge_type, reqData.Charge_way, reqData.Available, reqData.Status,
			reqData.Device_number, ntime, reqData.Request_status)
		if err != nil {
			log.Println(err)
			send_data.result = "false"
			send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		} else {
			log.Println(rows)
			send_data.result = "true"
			send_data.errStr = ""
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		}
	}
}

func DeviceRequestList(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select * from request_charge_device")
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		resultJson := jsonify.Jsonify(rows)

		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "request_charge_devices": resultJson})
	}
}

func DeviceRequestSubmit(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	config, err := setting.LoadConfigSettingJSON()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	reqData := RequestSubmitReq{}
	err = c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	//To b2b_service, it's permitted or reject
	body, _ := json.Marshal(reqData)
	resp, err := http.Post("http://"+config.Btb_service.Host+":"+config.Btb_service.Port+"/charge_device/request_submit", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "B2B Service 로 Request 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		defer resp.Body.Close()
		// Response 체크
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			send_data.result = "false"
			send_data.errStr = "B2B Service 의 Response 처리중 문제가 발생하였습니다."
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		} else {
			// Response Body 체크
			respJson := make(map[string]string)
			respStr := string(respBody)
			err := json.Unmarshal([]byte(respStr), &respJson)
			if err != nil {
				log.Println(err)
				send_data.result = "false"
				send_data.errStr = "B2B Service 의 Response Body Parsing 처리중 문제가 발생하였습니다."
				c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
			}

			log.Println(respStr)
			log.Println(respJson)

			// 이후 resp["result"] 체크 후 정상 MongoDB 에 Logging / Mysql에 Delete
			if respJson["result"] != "false" {
				log.Println(respJson["errStr"])
				send_data.result = "false"
				send_data.errStr = "B2B Service 의 Response 가 올바르지 않습니다. => " + respJson["errStr"]
				c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
			} else {
				// mysql delete from request_charge_staion where request_id = Request_id
				conn1 := database.NewMysqlConnection()
				defer conn1.Close()

				rows, err := conn1.Query("delete from request_charge_device where request_id = ?", reqData.Request_id)
				if err != nil {
					log.Println(err)
					send_data.result = "false"
					send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
					c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
				} else {
					log.Println(rows)

					// MongoDB logging
					// "Request_value" : "waiting"/"permitted"/"reject"
					request_value := reqData.Request_value
					ntime := time.Now().Format(time.RFC3339)
					ntime = ntime[:19]
					client := database.NewMongodbConnection()
					conn := client.Database("Admin_Service").Collection("request_charge_device")

					result, err := conn.InsertOne(context.TODO(), bson.D{
						{Key: "Request_uid", Value: reqData.Request_uid},
						{Key: "Station_id", Value: reqData.Station_id},
						{Key: "Name", Value: reqData.Name},
						{Key: "Sirial", Value: reqData.Sirial},
						{Key: "Charge_type", Value: reqData.Charge_type},
						{Key: "Charge_way", Value: reqData.Charge_way},
						{Key: "Available", Value: reqData.Available},
						{Key: "Status", Value: reqData.Status},
						{Key: "Device_number", Value: reqData.Device_number},
						{Key: "Request_value", Value: request_value},
						{Key: "Request_status", Value: reqData.Request_status},
						{Key: "Timestamp", Value: ntime},
					})
					if err != nil {
						log.Println(err)
						send_data.result = "false"
						send_data.errStr = "MongoDB logging 중 문제가 발생하였습니다."
						c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
					} else {
						log.Println(result)

						send_data.result = "true"
						send_data.errStr = ""
						c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
					}
				}
			}
		}
	}
}

func DeviceRequestHistory(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	client := database.NewMongodbConnection()
	conn := client.Database("Admin_Service").Collection("request_charge_device")

	cursor, err := conn.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		var result_arr []string
		for cursor.Next(context.TODO()) {
			var elem bson.M
			err := cursor.Decode(&elem)
			if err != nil {
				log.Println(err)
				send_data.result = "false"
				send_data.errStr = "DB Query Decode 중 문제가 발생하였습니다."
				c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
			}
			result, _ := json.Marshal(elem)
			result_arr = append(result_arr, string(result))
		}

		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "request_charge_devices_history": result_arr})
	}
}
