package report

import (
	"log"
	"net/http"

	"github.com/Youn-sik/KSCONNECT_ADMIN/database"
	"github.com/gin-gonic/gin"
)

// var TotalSendData []SendData

// // type SendData struct {
// // 	Company Company `json:"company"`
// // }
// // type Company struct {
// // 	Company_id   string `json:"company_id"`
// // 	Company_name string `json:"company_name"`
// // }

func getCompanyList() []int {
	var comapany_ids []int
	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select company_id from company")
	if err != nil {
		log.Println(err)
	} else {
		for rows.Next() {
			var company_id int
			err = rows.Scan(&company_id)
			if err != nil {
				log.Println(err)
			} else {
				comapany_ids = append(comapany_ids, company_id)
			}
		}
	}
	return comapany_ids
}

func getStationList(comapany_ids []int) [][]int {
	var Totalstation_ids [][]int
	var station_ids []int
	i := 0
	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	for _, cid := range comapany_ids {
		log.Print(cid)
		rows, err := conn1.Query("select station_id from charge_station where company_id = ?", cid)
		if err != nil {
			log.Println(err)
		} else {
			for rows.Next() {
				var station_id int
				err := rows.Scan(&station_id)
				if err != nil {
					log.Println(err)
				} else {
					station_ids = append(station_ids, station_id)
				}
			}
		}
		// 조금 이상하게 나옴
		Totalstation_ids = append(Totalstation_ids, station_ids)
		i++
	}
	return Totalstation_ids
}

func ReportList(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	// company list
	company_ids := getCompanyList()
	if len(company_ids) == 0 {
		send_data.result = "false"
		send_data.errStr = "등록 된 법인 고객이 존재하지 않습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}

	// charge_station list(count, ids)
	log.Println(getStationList(company_ids))

	// charge_device list(count)

	// charge_device usage, payment

	// rows, err := conn1.Query("select company.company_id, company.name, count(charge_station.station_id) from company inner join charge_station on company.company_id = charge_station.company_id where company.company_id = 1")
	// if err != nil {
	// 	log.Println(err)
	// 	send_data.result = "false"
	// 	send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
	// 	c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	// } else {
	// 	log.Println(rows)
	// }
}
