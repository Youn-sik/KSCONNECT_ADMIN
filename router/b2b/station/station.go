package station

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
	// Station_id int
	Company_id int
	Name       string
	Status     string
	// Last_state
	Address        string
	Address_detail string
	Available      string
	Park_fee       int
	Pay_type       string
	Lat            string
	Longi          string
	Purpose        string
	Guide          string
	Reg_date       string
	// Up_date        string
}

type UpdateReq struct {
	Station_id int
	// Company_id int
	Name   string
	Status string
	// Last_state
	Address        string
	Address_detail string
	Available      string
	Park_fee       int
	Pay_type       string
	Lat            string
	Longi          string
	Purpose        string
	Guide          string
	// Reg_date       string
	Up_date string
}

type DeleteReq struct {
	Station_id int
}

type RequestReq struct {
	// Request_id    int
	Station_id     int
	Request_uid    int
	Company_id     int
	Name           string
	Status         string
	Address        string
	Address_detail string
	Available      string
	Park_fee       int
	Pay_type       string
	Lat            string
	Longi          string
	Purpose        string
	Guide          string
	Request_status string
	// Request_value  string // wating
}

type RequestSubmitReq struct {
	Request_id     int    `json:"request_id"`
	Request_uid    int    `json:"request_uid"`
	Company_id     int    `json:"company_id"`
	Name           string `json:"name"`
	Status         string `json:"status"`
	Address        string `json:"address"`
	Address_detail string `json:"address_detail"`
	Available      string `json:"available"`
	Park_fee       int    `json:"park_fee"`
	Pay_type       string `json:"pay_type"`
	Lat            string `json:"lat"`
	Longi          string `json:"longi"`
	Purpose        string `json:"purpose"`
	Guide          string `json:"guide"`
	Request_status string `json:"request_status"`
	Request_value  string `json:"request_value"` // permitted or reject
}

func getDeviceCount(station_id int) int {
	conn1 := database.NewMysqlConnection()
	defer conn1.Close()
	var device_cnt int

	rows, err := conn1.Query("select count(device_id) as device_cnt from charge_device where station_id = ?", station_id)
	if err != nil {
		log.Println(err)
		return -1
	}
	for rows.Next() {
		err := rows.Scan(&device_cnt)
		if err != nil {
			log.Println(err)
			return -1
		}
	}
	return device_cnt
}

func StationList(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select cs.station_id, cs.name, cs.address, c.name, c.company_id from charge_station as cs inner join company as c on c.company_id = cs.company_id")
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query ?????? ??? ????????? ?????????????????????."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		type Send_result struct {
			Station_id      int    `json:"station_id"`
			Station_name    string `json:"station_name"`
			Station_address string `json:"station_address"`
			Company_name    string `json:"company_name"`
			Company_id      int    `json:"company_id"`
			Device_cnt      int    `json:"device_cnt"`
		}
		var send_arr []Send_result
		var send_result Send_result

		for rows.Next() {
			// var station_id int
			// var station_name string
			// var station_address string
			// var company_name string
			// var company_id int
			// var device_cnt int

			err = rows.Scan(&send_result.Station_id, &send_result.Station_name, &send_result.Station_address, &send_result.Company_name, &send_result.Company_id)
			if err != nil {
				log.Println(err)
				send_data.result = "false"
				send_data.errStr = "Query Parsing ??? ????????? ?????????????????????."
				c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
				return
			} else {
				send_result.Device_cnt = getDeviceCount(send_result.Station_id)
				send_arr = append(send_arr, send_result)
			}
		}
		// log.Printf("%+v\n", send_arr)
		// resultJson := jsonify.Jsonify(rows)

		send_data.result = "true"
		send_data.errStr = ""
		// c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "charge_stations": resultJson})
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "charge_stations": send_arr})
	}
}

func StationCreate(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	reqData := CreateReq{}
	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing ????????? ?????????????????????."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query(
		"insert into charge_station (company_id, name, status, address, address_detail, available, park_fee, pay_type, lat, longi, purpose, guide, reg_date, up_date) "+
			"value (?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		reqData.Company_id, reqData.Name, reqData.Status, reqData.Address, reqData.Address_detail, reqData.Available, reqData.Park_fee, reqData.Pay_type,
		reqData.Lat, reqData.Longi, reqData.Purpose, reqData.Guide, reqData.Reg_date, reqData.Reg_date)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query ?????? ??? ????????? ?????????????????????."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		log.Println(rows)
		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}
}

func StationUpdate(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	reqData := UpdateReq{}
	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing ????????? ?????????????????????."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query(
		"update charge_station set name = ?, status = ?, address = ?, address_detail = ?, available = ?, park_fee = ?, pay_type = ?, "+
			"lat = ?, longi = ?, purpose = ?, guide = ?, up_date = ? where station_id = ?",
		reqData.Name, reqData.Status, reqData.Address, reqData.Address_detail, reqData.Available, reqData.Park_fee, reqData.Pay_type,
		reqData.Lat, reqData.Longi, reqData.Purpose, reqData.Guide, reqData.Up_date, reqData.Station_id)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query ?????? ??? ????????? ?????????????????????."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		log.Println(rows)
		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}
}

func StationDelete(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	reqData := DeleteReq{}
	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing ????????? ?????????????????????."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query(
		"delete from charge_station where station_id = ?", reqData.Station_id)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query ?????? ??? ????????? ?????????????????????."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		log.Println(rows)
		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}
}

func StationRequest(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	reqData := RequestReq{}
	err := c.Bind(&reqData)
	log.Println(reqData)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing ????????? ?????????????????????."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}
	if reqData.Station_id != 0 && reqData.Request_status == "C" {
		send_data.result = "false"
		send_data.errStr = "?????? ????????? station_id = 0 ?????????."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}

	// MongoDB logging
	ntime := time.Now().Format(time.RFC3339)
	ntime = ntime[:19]
	client := database.NewMongodbConnection()
	conn := client.Database("Admin_Service").Collection("request_charge_station")

	result, err := conn.InsertOne(context.TODO(), bson.D{
		{Key: "Station_id", Value: reqData.Station_id},
		{Key: "Request_uid", Value: reqData.Request_uid},
		{Key: "Company_id", Value: reqData.Company_id},
		{Key: "Name", Value: reqData.Name},
		{Key: "Status", Value: reqData.Status},
		{Key: "Address", Value: reqData.Address},
		{Key: "Address_detail", Value: reqData.Address_detail},
		{Key: "Available", Value: reqData.Available},
		{Key: "Park_fee", Value: reqData.Park_fee},
		{Key: "Pay_type", Value: reqData.Pay_type},
		{Key: "Lat", Value: reqData.Lat},
		{Key: "Longi", Value: reqData.Longi},
		{Key: "Purpose", Value: reqData.Purpose},
		{Key: "Guide", Value: reqData.Guide},
		{Key: "Request_value", Value: "wating"},
		{Key: "Request_status", Value: reqData.Request_status},
		{Key: "Timestamp", Value: ntime},
	})
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "MongoDB logging ??? ????????? ?????????????????????."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		log.Println(result)

		// MYSQL input
		conn1 := database.NewMysqlConnection()
		defer conn1.Close()

		rows, err := conn1.Query("insert into request_charge_station (station_id, request_uid, company_id, name, status, address, address_detail, "+
			"available, park_fee, pay_type, lat, longi, purpose, guide, request_time, request_status) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
			reqData.Station_id, reqData.Request_uid, reqData.Company_id, reqData.Name, reqData.Status, reqData.Address, reqData.Address_detail,
			reqData.Available, reqData.Park_fee, reqData.Pay_type, reqData.Lat, reqData.Longi, reqData.Purpose, reqData.Guide, ntime, reqData.Request_status)
		if err != nil {
			log.Println(err)
			send_data.result = "false"
			send_data.errStr = "DB Query ?????? ??? ????????? ?????????????????????."
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		} else {
			log.Println(rows)
			send_data.result = "true"
			send_data.errStr = ""
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		}
	}
}

func StationRequestList(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select * from request_charge_station")
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query ?????? ??? ????????? ?????????????????????."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		resultJson := jsonify.Jsonify(rows)

		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "request_charge_stations": resultJson})
	}
}

func StationRequestSubmit(c *gin.Context) {
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
		send_data.errStr = "Body parsing ????????? ?????????????????????."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	//To b2b_service, it's permitted or reject
	body, _ := json.Marshal(reqData)
	resp, err := http.Post("http://"+config.Btb_service.Host+":"+config.Btb_service.Port+"/charge_station/request_submit", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "B2B Service ??? Request ??? ????????? ?????????????????????."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		defer resp.Body.Close()
		// Response ??????
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			send_data.result = "false"
			send_data.errStr = "B2B Service ??? Response ????????? ????????? ?????????????????????."
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		} else {
			// Response Body ??????
			respJson := make(map[string]string)
			respStr := string(respBody)
			err := json.Unmarshal([]byte(respStr), &respJson)
			if err != nil {
				log.Println(err)
				send_data.result = "false"
				send_data.errStr = "B2B Service ??? Response Body Parsing ????????? ????????? ?????????????????????."
				c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
			}

			log.Println(respStr)
			log.Println(respJson)

			// ?????? resp["result"] ?????? ??? ?????? MongoDB ??? Logging / Mysql??? Delete
			if respJson["result"] != "false" {
				log.Println(respJson["errStr"])
				send_data.result = "false"
				send_data.errStr = "B2B Service ??? Response ??? ???????????? ????????????. => " + respJson["errStr"]
				c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
			} else {
				// mysql delete from request_charge_staion where request_id = Request_id
				conn1 := database.NewMysqlConnection()
				defer conn1.Close()

				rows, err := conn1.Query("delete from request_charge_station where request_id = ?", reqData.Request_id)
				if err != nil {
					log.Println(err)
					send_data.result = "false"
					send_data.errStr = "DB Query ?????? ??? ????????? ?????????????????????."
					c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
				} else {
					log.Println(rows)

					// MongoDB logging
					// "Request_value" : "waiting"/"permitted"/"reject"
					request_value := reqData.Request_value
					ntime := time.Now().Format(time.RFC3339)
					ntime = ntime[:19]
					client := database.NewMongodbConnection()
					conn := client.Database("Admin_Service").Collection("request_charge_station")

					result, err := conn.InsertOne(context.TODO(), bson.D{
						{Key: "Request_uid", Value: reqData.Request_uid},
						{Key: "Company_id", Value: reqData.Company_id},
						{Key: "Name", Value: reqData.Name},
						{Key: "Status", Value: reqData.Status},
						{Key: "Address", Value: reqData.Address},
						{Key: "Address_detail", Value: reqData.Address_detail},
						{Key: "Available", Value: reqData.Available},
						{Key: "Park_fee", Value: reqData.Park_fee},
						{Key: "Pay_type", Value: reqData.Pay_type},
						{Key: "Lat", Value: reqData.Lat},
						{Key: "Longi", Value: reqData.Longi},
						{Key: "Purpose", Value: reqData.Purpose},
						{Key: "Guide", Value: reqData.Guide},
						{Key: "Request_value", Value: request_value},
						{Key: "Request_status", Value: reqData.Request_status},
						{Key: "Timestamp", Value: ntime},
					})
					if err != nil {
						log.Println(err)
						send_data.result = "false"
						send_data.errStr = "MongoDB logging ??? ????????? ?????????????????????."
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

func StationRequestHistory(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	client := database.NewMongodbConnection()
	conn := client.Database("Admin_Service").Collection("request_charge_station")

	cursor, err := conn.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query ?????? ??? ????????? ?????????????????????."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		var result_arr []string
		for cursor.Next(context.TODO()) {
			var elem bson.M
			err := cursor.Decode(&elem)
			if err != nil {
				log.Println(err)
				send_data.result = "false"
				send_data.errStr = "DB Query Decode ??? ????????? ?????????????????????."
				c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
			}
			result, _ := json.Marshal(elem)
			result_arr = append(result_arr, string(result))
		}

		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "request_charge_stations_history": result_arr})
	}
}
