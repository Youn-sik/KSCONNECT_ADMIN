package user_account

import (
	"log"
	"net/http"

	"github.com/Youn-sik/KSCONNECT_ADMIN/database"
	"github.com/bdwilliams/go-jsonify/jsonify"
	"github.com/gin-gonic/gin"
)

type CreateReq struct {
	// Uid                    string
	Id                     string
	Password               string
	Name                   string
	Email                  string
	Mobile                 string
	Address                string
	Car_model              string
	Car_number             string
	Payment_card_company   string
	Payment_card_number    string
	Membership_card_number string
	// Point                  int
	Rfid string
}

type UpdateReq struct {
	Uid                    int
	Id                     string
	Password               string
	Name                   string
	Email                  string
	Mobile                 string
	Address                string
	Car_model              string
	Car_number             string
	Payment_card_company   string
	Payment_card_number    string
	Membership_card_number string
	Point                  int
	Rfid                   string
}

type DeleteReq struct {
	Uid int
}

func UserList(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select * from user")
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
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "users": resultJson})
	}
}

func UserCreate(c *gin.Context) {
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

	rows, err := conn1.Query("insert into user (id, password, name, email, mobile, address, car_model, car_number, payment_card_company, payment_card_number, membership_card_number, rfid) "+
		"value (?,?,?,?,?,?,?,?,?,?,?,?)",
		reqData.Id, reqData.Password, reqData.Name, reqData.Email, reqData.Mobile, reqData.Address, reqData.Car_model, reqData.Car_number, reqData.Payment_card_company, reqData.Payment_card_number, reqData.Membership_card_number, reqData.Rfid)
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

func UserUpdate(c *gin.Context) {
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

	rows, err := conn1.Query("update user set id = ?, password = ?, name = ?, email = ?, mobile = ?, address = ?, car_model = ?, car_number = ?, "+
		"payment_card_company = ?, payment_card_number = ?, membership_card_number = ?, point = ?, rfid = ? where uid = ?",
		reqData.Id, reqData.Password, reqData.Name, reqData.Email, reqData.Mobile, reqData.Address, reqData.Car_model, reqData.Car_number,
		reqData.Payment_card_company, reqData.Payment_card_number, reqData.Membership_card_number, reqData.Point, reqData.Rfid, reqData.Uid)
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

func UserDelete(c *gin.Context) {
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
		"delete from user where uid = ?", reqData.Uid)
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
