package device

import (
	"log"
	"net/http"

	"github.com/Youn-sik/KSCONNECT_ADMIN/database"
	"github.com/bdwilliams/go-jsonify/jsonify"
	"github.com/gin-gonic/gin"
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

func DeviceList(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select * from charge_device")
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

func DeviceCreate(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	reqData := CreateReq{}
	err := c.Bind(&reqData)
	log.Println(reqData)

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
