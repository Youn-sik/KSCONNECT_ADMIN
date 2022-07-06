package user_account

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/Youn-sik/KSCONNECT_ADMIN/database"
	"github.com/Youn-sik/KSCONNECT_ADMIN/setting"
	"github.com/bdwilliams/go-jsonify/jsonify"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type CreateReq struct {
	// Uid                    int
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

type MembershipCardCreateReq struct {
	Id             string
	Request_way    string
	Request_status string
	Request_reason string
	Address        string
}

type MembershipCardUpdateReq struct {
	Request_uid            int
	Membership_card_number string
	Request_way            string
	Request_status         string
	Request_reason         string
	Address                string
}

type MembershipCardDeleteReq struct {
	Request_uid int
}

type MembershipCardRequestReq struct {
	Request_uid     int
	Request_way     string
	Request_reason  string
	Request_address string
}

type MembershipCardRequestSubmitReq struct {
	Request_id      int    `json:"request_id"`
	Request_uid     int    `json:"request_uid"`
	Request_way     string `json:"request_way"`
	Request_reason  string `json:"request_reason"`
	Request_address string `json::"request_address"`
	Request_value   string `json:"request_value"` // permitted or reject
}

type InquiryBoardReplyReq struct {
	Inquiry_id int
	Reply      string
}

func UserList(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select uid, id, name, email, mobile, address, car_model, car_number, payment_card_company, payment_card_number, membership_card_number, point, rfid from user")
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		resultJson := jsonify.Jsonify(rows)

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

	_, err = conn1.Query("insert into user (id, password, name, email, mobile, address, car_model, car_number, payment_card_company, payment_card_number, membership_card_number, rfid) "+
		"value (?,?,?,?,?,?,?,?,?,?,?,?)",
		reqData.Id, reqData.Password, reqData.Name, reqData.Email, reqData.Mobile, reqData.Address, reqData.Car_model, reqData.Car_number, reqData.Payment_card_company, reqData.Payment_card_number, reqData.Membership_card_number, reqData.Rfid)
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

	if reqData.Password == "" {
		_, err = conn1.Query("update user set id = ?, name = ?, email = ?, mobile = ?, address = ?, car_model = ?, car_number = ?, "+
			"payment_card_company = ?, payment_card_number = ?, membership_card_number = ?, point = ?, rfid = ? where uid = ?",
			reqData.Id, reqData.Name, reqData.Email, reqData.Mobile, reqData.Address, reqData.Car_model, reqData.Car_number,
			reqData.Payment_card_company, reqData.Payment_card_number, reqData.Membership_card_number, reqData.Point, reqData.Rfid, reqData.Uid)
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
	} else {
		_, err = conn1.Query("update user set id = ?, password = ?, name = ?, email = ?, mobile = ?, address = ?, car_model = ?, car_number = ?, "+
			"payment_card_company = ?, payment_card_number = ?, membership_card_number = ?, point = ?, rfid = ? where uid = ?",
			reqData.Id, reqData.Password, reqData.Name, reqData.Email, reqData.Mobile, reqData.Address, reqData.Car_model, reqData.Car_number,
			reqData.Payment_card_company, reqData.Payment_card_number, reqData.Membership_card_number, reqData.Point, reqData.Rfid, reqData.Uid)
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

	_, err = conn1.Query(
		"delete from user where uid = ?", reqData.Uid)
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

func MembershipCardList(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select request_uid, user_membership_card.membership_card_number, request_way, request_status, user_membership_card.address, " +
		"request_time, request_reason , name, car_number from user_membership_card  inner join user on request_uid = uid")
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		resultJson := jsonify.Jsonify(rows)

		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "users": resultJson})
	}
}

func getUserUid(id string) int {
	send_uid := 0

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select uid from user where id = ?", id)
	if err != nil {
		log.Println(err)
		return -1
	}
	for rows.Next() {
		var uid int
		err = rows.Scan(&uid)
		if err != nil {
			log.Println(err)
			return -1
		}
		send_uid = uid
	}
	return send_uid
}

func getMembershipCardNumber(uid int) string {
	ntime := time.Now().Format(time.RFC3339)
	ntime = ntime[:19]

	random_value := strconv.Itoa(rand.Intn(9999))
	switch len(random_value) {
	case 1:
		random_value = "000" + random_value
	case 2:
		random_value = "00" + random_value
	case 3:
		random_value = "0" + random_value
	}

	uid_value := strconv.Itoa(uid)
	switch len(uid_value) {
	case 1:
		uid_value = "000" + uid_value
	case 2:
		uid_value = "00" + uid_value
	case 3:
		uid_value = "0" + uid_value
	}

	membership_card_number1 := ntime[:4]
	membership_card_number2 := ntime[5:7] + ntime[8:10]
	membership_card_number3 := random_value
	membership_card_number4 := uid_value
	membership_card_number := membership_card_number1 + "-" + membership_card_number2 + "-" + membership_card_number3 + "-" + membership_card_number4

	return membership_card_number
}

func MembershipCardCreate(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	reqData := MembershipCardCreateReq{}
	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	userUid := getUserUid(reqData.Id)
	switch userUid {
	case 0:
		send_data.result = "false"
		send_data.errStr = "존재하지 않는 아이디입니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	case -1:
		send_data.result = "false"
		send_data.errStr = "아이디를 가져오는 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}

	ntime := time.Now().Format(time.RFC3339)
	ntime = ntime[:19]

	// random_value := strconv.Itoa(rand.Intn(9999))
	// switch len(random_value) {
	// case 1:
	// 	random_value = "000" + random_value
	// case 2:
	// 	random_value = "00" + random_value
	// case 3:
	// 	random_value = "0" + random_value
	// }

	// uid_value := strconv.Itoa(userUid)
	// switch len(uid_value) {
	// case 1:
	// 	uid_value = "000" + uid_value
	// case 2:
	// 	uid_value = "00" + uid_value
	// case 3:
	// 	uid_value = "0" + uid_value
	// }

	// membership_card_number1 := ntime[:4]
	// membership_card_number2 := ntime[5:7] + ntime[8:10]
	// membership_card_number3 := random_value
	// membership_card_number4 := uid_value
	// membership_card_number := membership_card_number1 + "-" + membership_card_number2 + "-" + membership_card_number3 + "-" + membership_card_number4
	membership_card_number := getMembershipCardNumber(userUid)

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	_, err = conn1.Query("insert into user_membership_card (request_uid, membership_card_number, request_way, request_status, request_time, request_reason, address) value (?,?,?,?,?,?,?)",
		userUid, membership_card_number, reqData.Request_way, reqData.Request_status, ntime, reqData.Request_reason, reqData.Address)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}

	_, err = conn1.Query("update user set membership_card_number = ? where uid = ?", membership_card_number, userUid)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}

	send_data.result = "true"
	send_data.errStr = ""
	c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})

}

func MembershipCardUpdate(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	reqData := MembershipCardUpdateReq{}
	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	_, err = conn1.Query("update user_membership_card set request_way = ?, request_status = ?, request_reason = ?, address = ? where request_uid = ?",
		reqData.Request_way, reqData.Request_status, reqData.Request_reason, reqData.Address, reqData.Request_uid)
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

func MembershipCardDelete(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	reqData := MembershipCardDeleteReq{}
	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	_, err = conn1.Query("delete from user_membership_card where request_uid = ?", reqData.Request_uid)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}
	_, err = conn1.Query("update user set membership_card_number=NULL where uid = ?", reqData.Request_uid)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}

	send_data.result = "true"
	send_data.errStr = ""
	c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
}

func MembershipCardRequest(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	reqData := MembershipCardRequestReq{}
	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	ntime := time.Now().Format(time.RFC3339)
	ntime = ntime[:19]

	// MongoDB logging
	client := database.NewMongodbConnection()
	conn := client.Database("Admin_Service").Collection("request_membership_card")

	_, err = conn.InsertOne(context.TODO(), bson.D{
		{Key: "Request_uid", Value: reqData.Request_uid},
		{Key: "Request_way", Value: reqData.Request_way},
		{Key: "Request_reason", Value: reqData.Request_reason},
		{Key: "Request_address", Value: reqData.Request_address},
		{Key: "Request_value", Value: "wating"},
		{Key: "Timestamp", Value: ntime},
	})
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "MongoDB logging 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}
	// else {
	// 	log.Println(result)
	// }

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	_, err = conn1.Query("insert into request_user_membership_card (request_uid, request_time, request_way, request_reason, request_address) value (?,?,?,?,?)",
		reqData.Request_uid, ntime, reqData.Request_way, reqData.Request_reason, reqData.Request_address)
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

func MembershipCardRequestList(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select request_id, request_uid, request_time, request_way, request_reason, request_address, name, car_number, id " +
		"from request_user_membership_card inner join user on request_uid = uid")
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		resultJson := jsonify.Jsonify(rows)

		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "users": resultJson})
	}
}

func MembershipCardRequestSubmit(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	config, err := setting.LoadConfigSettingJSON()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	ntime := time.Now().Format(time.RFC3339)
	ntime = ntime[:19]

	reqData := MembershipCardRequestSubmitReq{}
	err = c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	//To user_service, it's permitted or reject
	var reqBody struct {
		ReqData              MembershipCardRequestSubmitReq `json:"reqData"`
		MembershipCardNumber string                         `json:"membershipCardNumber"`
		Timestamp            string                         `json:"timestamp"`
	}
	membershipCardNumber := getMembershipCardNumber(reqData.Request_uid)

	reqBody.ReqData = reqData
	reqBody.MembershipCardNumber = membershipCardNumber
	reqBody.Timestamp = ntime
	body, err := json.Marshal(reqBody)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다.."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	resp, err := http.Post("http://"+config.User_service.Host+":"+config.User_service.Port+"/user/membership_card_request_submit", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "User Service 로 Request 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		defer resp.Body.Close()
		// Response 체크
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			send_data.result = "false"
			send_data.errStr = "User Service 의 Response 처리중 문제가 발생하였습니다."
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		} else {
			// Response Body 체크
			respJson := make(map[string]string)
			respStr := string(respBody)
			err := json.Unmarshal([]byte(respStr), &respJson)
			if err != nil {
				log.Println(err)
				send_data.result = "false"
				send_data.errStr = "User Service 의 Response Body Parsing 처리중 문제가 발생하였습니다."
				c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
			}

			// log.Println(respStr)
			// log.Println(respJson)
			log.Println(respJson["result"])

			// 이후 resp["result"] 체크 후 정상 MongoDB 에 Logging / Mysql에 Delete
			if respJson["result"] == "false" {
				log.Println(respJson["errStr"])
				send_data.result = "false"
				send_data.errStr = "User Service 의 Response 가 올바르지 않습니다. => " + respJson["errStr"]
				c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
			} else {
				// mysql delete from request_charge_staion where request_id = Request_id
				conn1 := database.NewMysqlConnection()
				defer conn1.Close()

				_, err := conn1.Query("delete from request_user_membership_card where request_id = ?", reqData.Request_id)
				if err != nil {
					log.Println(err)
					send_data.result = "false"
					send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
					c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
				} else {
					// log.Println(rows)

					// MongoDB logging
					// "Request_value" : "waiting"/"permitted"/"reject"
					request_value := reqData.Request_value
					client := database.NewMongodbConnection()
					conn := client.Database("Admin_Service").Collection("request_membership_card")

					_, err := conn.InsertOne(context.TODO(), bson.D{
						{Key: "Request_uid", Value: reqData.Request_uid},
						{Key: "Request_way", Value: reqData.Request_way},
						{Key: "Request_reason", Value: reqData.Request_reason},
						{Key: "Request_address", Value: reqData.Request_address},
						{Key: "Request_value", Value: request_value},
						{Key: "Timestamp", Value: ntime},
					})
					if err != nil {
						log.Println(err)
						send_data.result = "false"
						send_data.errStr = "MongoDB logging 중 문제가 발생하였습니다."
						c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
					} else {
						// log.Println(result)

						send_data.result = "true"
						send_data.errStr = ""
						c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
					}
				}
			}
		}
	}
}

func MembershipCardRequestHistory(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	client := database.NewMongodbConnection()
	conn := client.Database("Admin_Service").Collection("request_membership_card")

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

func InquiryBoardList(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select * from inquiry_board")
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		resultJson := jsonify.Jsonify(rows)

		send_data.result = "true"
		send_data.errStr = ""
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "users": resultJson})
	}
}

func InquiryBoardReply(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	reqData := InquiryBoardReplyReq{}
	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	conn1 := database.NewMysqlConnection()
	_, err = conn1.Query("update inquiry_board set reply = ?, status = 'Y' where inquiry_id = ?", reqData.Reply, reqData.Inquiry_id)
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
