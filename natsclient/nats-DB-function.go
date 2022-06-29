package natsclient

import (
	"context"
	"database/sql"
	"log"

	"github.com/Youn-sik/KSCONNECT_ADMIN/database"
	v16 "github.com/aliml92/ocpp/v16"
	"github.com/bdwilliams/go-jsonify/jsonify"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client
var MysqlClient *sql.DB

// MongoClient = database.NewMongodbConnection()
// conn := MongoClient.Database("Admin_Service").Collection("request_charge_device")
// MysqlClient = database.NewMysqlConnection()
// defer MongoDB.Close()
// defer MYSQL.Close()

func UpdateMeterValue() {
	rows, err := MysqlClient.Query("select * from charge_device where status = 'I'")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	resultJson := jsonify.Jsonify(rows)
	if len(resultJson) == 0 {
		log.Println("From Mongo To Mysql Meter Value Update Block1111")
		return
	} else {
		// Get MongoDB And Update
		log.Println("From Mongo To Mysql Meter Value Update Block2222")
	}
}

func MeterValuesReq(GmeterValueReq v16.MeterValuesReq) {
	// meterValueReq := GmeterValueReq.MeterValue

	// MongoDB meter logging
	MongoClient = database.NewMongodbConnection()
	conn := MongoClient.Database("Admin_Service").Collection("OCPP_meter")

	result, err := conn.InsertOne(context.TODO(), bson.D{
		// {Key: "chargePointId", Value: GmeterValueReq.ChargePointId},
		// {Key: "connectorId", Value: meterValueReq.ConnectorId},
		// {Key: "transactionId", Value: meterValueReq.TransactionId},
		// {Key: "meterValue", Value: meterValueReq.MeterValue},
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Println(result)
	}

}

// RDB Charge Device 충전량 값 update => stop transaction 때 (meterStop - meterStart) 만큼 작업하기
