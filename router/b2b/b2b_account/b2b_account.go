package b2b_account

import (
	"log"
	"net/http"

	"github.com/Youn-sik/KSCONNECT_ADMIN/database"
	"github.com/bdwilliams/go-jsonify/jsonify"
	"github.com/gin-gonic/gin"
)

type CreateReq struct {
	// Uid                  int
	Id                   string
	Password             string
	Name                 string
	Email                string
	Mobile               string
	Zip_no               string
	Address              string
	Car_model            string
	Car_number           string
	Payment_card_company string
	Payment_card_number  string
	Pay_type             string
	Registration         string
	Certificate          string
	Warrant              string
	Copyid               string
	Cal_name             string
	Cal_email            string
	Cal_mobile           string
	Cal_tel              string
	Cal_fax              string
	Sms_status           string
	Email_status         string
}

func UserList(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select * from company_manager")
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

	_, err = conn1.Query("insert into company_manager () " +
		"value ()",
	)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		// log.Println(rows)
		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}
}
