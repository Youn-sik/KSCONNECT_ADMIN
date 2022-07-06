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
	Company_name                 string `json:"company_name"`
	Id                           string `json:"id"`
	Password                     string `json:"password"`
	Company_president            string `json:"company_president"`
	Company_certification_number string `json:"company_certification_number"`
	Company_certification_date   string `json:"company_certification_date"`
	Company_type                 string `json:"company_type"`
	Company_job_type             string `json:"company_job_type"`
	Name                         string `json:"name"`
	Company_certification        string `json:"company_certification"`
	Pay_type                     string `json:"pay_type"`
	Pay_company                  string `json:"pay_company"`
	Pay_card_number              string `json:"pay_card_number"`
}

type UpdateReq struct {
	Uid                          int    `json:"uid"`
	Company_id                   int    `json:"company_id"`
	Company_name                 string `json:"company_name"`
	Company_president            string `json:"company_president"`
	Company_certification_number string `json:"company_certification_number"`
	Company_certification_date   string `json:"company_certification_date"`
	Company_type                 string `json:"company_type"`
	Company_job_type             string `json:"company_job_type"`
	Id                           string `json:"id"`
	Password                     string `json:"password"`
	Name                         string `json:"name"`
	Company_certification        string `json:"company_certification"`
	Pay_type                     string `json:"pay_type"`
	Pay_company                  string `json:"pay_company"`
	Pay_card_number              string `json:"pay_card_number"`
}

type DeleteReq struct {
	Uid        int `json:"uid"`
	Company_id int `json:"company_id"`
}

//CRUD 시 유의 사항 -> C: 회사 등록도 함께, U: 회사 정보도 함께, D: 회사 정보도 함께

func UserList(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select cm.uid, cm.company_id, c.name as company_name, cm.id, cm.company_president, cm.company_certification_number, cm.company_certification_date, " +
		"cm.company_type, cm.company_job_type, cm.name, cm.pay_type, cm.pay_company, cm.pay_card_number from company_manager as cm inner join company as c on cm.company_id = c.company_id")
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		resultJson := jsonify.Jsonify(rows)
		// log.Println(resultJson)

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

	// log.Printf("%+v\n", reqData)
	_, err = conn1.Query("insert into company (name) value(?)", reqData.Company_name)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다.."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}
	// rows에서 company_id 받아서 company_manager 에 company FK 하나 만들어서 저장 같이 하기
	rows, err := conn1.Query("select company_id from company where name = ?", reqData.Company_name)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}
	var cid int
	for rows.Next() {
		err := rows.Scan(&cid)
		if err != nil {
			log.Println(err)
			send_data.result = "false"
			send_data.errStr = "Query Parsing 중 문제가 발생하였습니다."
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
			return
		}
	}

	if reqData.Company_certification == "" {
		_, err = conn1.Query("insert into company_manager (id, company_id, password, company_president, company_certification_number, company_certification_date, company_type, company_job_type, "+
			"name, pay_type, pay_company, pay_card_number) value (?,?,?,?,?,?,?,?,?,?,?,?)", reqData.Id, cid, reqData.Password, reqData.Company_president, reqData.Company_certification_number,
			reqData.Company_certification_date, reqData.Company_type, reqData.Company_job_type, reqData.Name, reqData.Pay_type, reqData.Pay_company, reqData.Pay_card_number)
		if err != nil {
			log.Println(err)
			send_data.result = "false"
			send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
			return
		}
	} else {
		_, err = conn1.Query("insert into company_manager (id, password, company_president, company_certification_number, company_certification_date, company_type, company_job_type, "+
			"name, company_certification, pay_type, pay_company, pay_card_number) value (?,?,?,?,?,?,?,?,?,?,?,?)", reqData.Id, reqData.Password, reqData.Company_president,
			reqData.Company_certification_number, reqData.Company_certification_date, reqData.Company_type, reqData.Company_job_type, reqData.Name, reqData.Company_certification,
			reqData.Pay_type, reqData.Pay_company, reqData.Pay_card_number)
		if err != nil {
			log.Println(err)
			send_data.result = "false"
			send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
			return
		}
	}

	send_data.result = "true"
	send_data.errStr = ""
	c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
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

	_, err = conn1.Query("update company set name = ? where company_id = ?", reqData.Company_name, reqData.Company_id)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다.."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}
	if reqData.Company_certification != "" {
		_, err = conn1.Query("update company_manager set id = ?, password = ?, company_president = ?, company_certification_number = ?, company_certification_date = ?, "+
			"company_type = ?, company_job_type = ? , name = ? , company_certification = ?, pay_type = ?, pay_company = ?, pay_card_number = ? where uid = ?",
			reqData.Id, reqData.Password, reqData.Company_president, reqData.Company_certification_number, reqData.Company_certification_date,
			reqData.Company_type, reqData.Company_job_type, reqData.Name, reqData.Company_certification, reqData.Pay_type, reqData.Pay_company, reqData.Pay_card_number, reqData.Uid)
		if err != nil {
			log.Println(err)
			send_data.result = "false"
			send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다.."
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
			return
		}
	} else {
		_, err = conn1.Query("update company_manager set id = ?, password = ?, company_president = ?, company_certification_number = ?, company_certification_date = ?, "+
			"company_type = ?, company_job_type = ? , name = ? , pay_type = ?, pay_company = ?, pay_card_number = ? where uid = ?",
			reqData.Id, reqData.Password, reqData.Company_president, reqData.Company_certification_number, reqData.Company_certification_date,
			reqData.Company_type, reqData.Company_job_type, reqData.Name, reqData.Pay_type, reqData.Pay_company, reqData.Pay_card_number, reqData.Uid)
		if err != nil {
			log.Println(err)
			send_data.result = "false"
			send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다.."
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
			return
		}
	}

	send_data.result = "true"
	send_data.errStr = ""
	c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
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

	_, err = conn1.Query("delete from company_manager where uid = ?", reqData.Uid)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}
	_, err = conn1.Query("delete from company where company_id = ?", reqData.Company_id)
	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "Body parsing 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	send_data.result = "true"
	send_data.errStr = ""
	c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
}
