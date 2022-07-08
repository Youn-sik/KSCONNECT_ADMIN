package report

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/Youn-sik/KSCONNECT_ADMIN/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type CompanyList struct {
	Company_ids   []int    `json:"company_ids"`
	Company_names []string `json:"company_names"`
}
type StationList struct {
	Company_id  int   `json:"company_id"`
	StationList []int `json:"station_list"`
}
type DeviceList struct {
	Station_id int   `json:"station_id"`
	DeviceList []int `json:"device_id"`
}
type DeviceUsage struct {
	Device_id int    `json:"station_id"`
	Usage     int    `json:"usage"`
	Payment   string `json:"payment"` // 어떻게 측정 ?
}
type UsageList struct {
	Device_id     int `json:"device_id"`
	Device_number int `json:"device_number"`
	Usage         int `json:"usage"`
}
type Payment struct {
	Payment string `json:"Payment"`
}
type Report struct {
	Company_id     int    `json:"company_id`
	Company_name   string `json:"company_name"`
	Station_count  int    `json:"station_count"`
	Device_count   int    `json:"device_count"`
	Device_usage   int    `json:"device_usage"`
	Device_payment int    `json:"device_payment"`
}

func getCompanyList() CompanyList {
	companyList := CompanyList{}
	var comapany_ids []int
	var comapany_names []string

	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select company_id, name from company order by company_id ASC")
	if err != nil {
		log.Println(err)
	} else {
		for rows.Next() {
			var company_id int
			var company_name string
			err = rows.Scan(&company_id, &company_name)
			if err != nil {
				log.Println(err)
			} else {
				comapany_ids = append(comapany_ids, company_id)
				comapany_names = append(comapany_names, company_name)
			}
		}
		companyList.Company_ids = comapany_ids
		companyList.Company_names = comapany_names
	}
	return companyList
}

// func getCompanyNameList() []string {
// 	var comapany_names []string
// 	conn1 := database.NewMysqlConnection()
// 	defer conn1.Close()

// 	rows, err := conn1.Query("select name from company order by company_id ASC")
// 	if err != nil {
// 		log.Println(err)
// 	} else {
// 		for rows.Next() {
// 			var comapany_name string
// 			err = rows.Scan(&comapany_name)
// 			if err != nil {
// 				log.Println(err)
// 			} else {
// 				comapany_names = append(comapany_names, comapany_name)
// 			}
// 		}
// 	}
// 	return comapany_names
// }

func getStationList(comapany_ids []int) []StationList {
	resultStationList := []StationList{}
	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	for _, cid := range comapany_ids {
		stationList := StationList{}
		stationList.Company_id = cid
		rows, err := conn1.Query("select station_id from charge_station where company_id = ?", cid)
		if err != nil {
			log.Println(err)
			return resultStationList
		} else {
			var stationIdArr []int
			for rows.Next() {
				var station_id int
				err := rows.Scan(&station_id)
				if err != nil {
					log.Println(err)
					return resultStationList
				}
				stationIdArr = append(stationIdArr, station_id)
			}
			stationList.StationList = stationIdArr
			resultStationList = append(resultStationList, stationList)
		}
	}
	return resultStationList
}

func getDeviceList(station_ids []int) []DeviceList {
	resultDeviceList := []DeviceList{}
	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	for _, sid := range station_ids {
		DeviceList := DeviceList{}
		DeviceList.Station_id = sid
		rows, err := conn1.Query("select device_id from charge_device where station_id = ?", sid)
		if err != nil {
			log.Println(err)
			return resultDeviceList
		} else {
			var DeviceIdArr []int
			for rows.Next() {
				var Device_id int
				err := rows.Scan(&Device_id)
				if err != nil {
					log.Println(err)
					return resultDeviceList
				}
				DeviceIdArr = append(DeviceIdArr, Device_id)
			}
			DeviceList.DeviceList = DeviceIdArr
			resultDeviceList = append(resultDeviceList, DeviceList)
		}
	}
	return resultDeviceList
}

func getDeviceUsage(device_ids []int) []UsageList {
	resultUsageList := []UsageList{}
	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	for _, did := range device_ids {
		usageList := UsageList{}
		usageList.Device_id = did
		rows, err := conn1.Query("select device_number, charge_device.usage from charge_device where device_id = ?", did)
		if err != nil {
			log.Println(err)
			return resultUsageList
		} else {
			var usage int
			var device_number int
			for rows.Next() {
				// var usage int
				err := rows.Scan(&device_number, &usage)
				if err != nil {
					log.Println(err)
					return resultUsageList
				}
				// DeviceIdArr = append(DeviceIdArr, Device_id)
			}
			usageList.Usage = usage
			usageList.Device_number = device_number
			resultUsageList = append(resultUsageList, usageList)
		}
	}
	return resultUsageList
}

func getPayment(station_id, device_number int) int {
	totpayment := 0
	sid := strconv.Itoa(station_id)

	client := database.NewMongodbConnection()
	conn := client.Database("Admin_Service").Collection("service_payment")
	filter := bson.M{"$and": []bson.M{{"Station": sid}, {"Device": device_number}}}
	cursor, err := conn.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return totpayment
	}
	for cursor.Next(context.TODO()) {
		payment := Payment{}
		if err := cursor.Decode(&payment); err != nil {
			log.Println(err)
			return totpayment
		}
		pay, _ := strconv.Atoi(payment.Payment)
		totpayment += pay
	}
	return totpayment
}

func ReportList(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}

	// company list
	companyList := getCompanyList()
	if len(companyList.Company_ids) == 0 {
		send_data.result = "false"
		send_data.errStr = "등록 된 법인 고객이 존재하지 않습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}
	// log.Println(companyList)
	// log.Printf("\n\n")
	// log.Println(len(company_ids))

	// stationList {company_id int, []staionList}
	stationList := getStationList(companyList.Company_ids)
	if len(stationList) == 0 {
		send_data.result = "false"
		send_data.errStr = "등록 된 충전소가 존재하지 않습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}
	// log.Printf("%+v\n", stationList)
	// log.Printf("\n\n")
	// log.Println(len(stationList))

	// charge_device list(station_id int, []deivceList)
	deviceList := []DeviceList{}
	for _, sl := range stationList {
		device_list := getDeviceList(sl.StationList)
		deviceList = append(deviceList, device_list...)
	}
	if len(deviceList) == 0 {
		send_data.result = "false"
		send_data.errStr = "등록 된 충전기가 존재하지 않습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		return
	}
	// log.Printf("%+v\n", deviceList)
	// log.Printf("\n\n")
	// log.Println(len(deviceList))

	// charge_device usage, payment
	usageList := []UsageList{}
	for _, dl := range deviceList {
		usage_list := getDeviceUsage(dl.DeviceList)
		usageList = append(usageList, usage_list...)
	}
	// log.Printf("%+v\n", usageList)
	// log.Printf("\n\n")
	// log.Println(len(usageList))

	resultReport := []Report{}
	for i, cid := range companyList.Company_ids {
		report := Report{}
		report.Company_id = cid
		report.Company_name = companyList.Company_names[i]
		report.Station_count = len(stationList[i].StationList)
		report.Device_count = 0
		report.Device_usage = 0
		report.Device_payment = 0

		for _, dl := range deviceList {
			// if contains(sl.StationList, dl.Station_id) {
			if contains(stationList[i].StationList, dl.Station_id) {
				report.Device_count += len(dl.DeviceList)
				for _, usage := range usageList {
					if contains(dl.DeviceList, usage.Device_id) {
						report.Device_usage += usage.Usage
						report.Device_payment += getPayment(dl.Station_id, usage.Device_number)
					}
				}
			}
		}
		// log.Printf("%+v\n", report)
		resultReport = append(resultReport, report)
	}
	send_data.result = "true"
	send_data.errStr = ""
	c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "resultReport": resultReport})
	return
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
