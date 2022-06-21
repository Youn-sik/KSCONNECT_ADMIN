package station

import (
	"log"
	"net/http"

	"github.com/Youn-sik/KSCONNECT_ADMIN/database"
	"github.com/bdwilliams/go-jsonify/jsonify"
	"github.com/gin-gonic/gin"
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

func StationList(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select * from charge_station")
	if err != nil {
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	resultJson := jsonify.Jsonify(rows)
	log.Println(resultJson)

	send_data.result = "true"
	send_data.errStr = ""
	c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "charge_stations": resultJson})
}

func StationCreate(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	reqData := CreateReq{}
	err := c.Bind(&reqData)

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
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
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
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
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

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query(
		"delete from charge_station where station_id = ?", reqData.Station_id)
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
