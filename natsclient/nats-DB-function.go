package natsclient

import (
	"context"
	"database/sql"
	"log"

	"github.com/Youn-sik/KSCONNECT_ADMIN/database"
	v16 "github.com/aliml92/ocpp/v16"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MysqlClient *sql.DB

// MongoClient = database.NewMongodbConnection()
// conn := MongoClient.Database("Admin_Service").Collection("request_charge_device")
// defer conn.Close()
// MysqlClient = database.NewMysqlConnection()
// defer MysqlClient.Close()

func UpdateMeterValue() {
	log.Println("---UpdateMeterValueFunc---")
	MysqlClient = database.NewMysqlConnection()
	defer MysqlClient.Close()

	rows, err := MysqlClient.Query("select station_id, device_number from charge_device where status = 'I'")
	if err != nil {
		log.Fatal(err)
		panic(err)
	} else {
		MongoClient = database.NewMongodbConnection()
		conn := MongoClient.Database("Admin_Service").Collection("ocpp_MeterValues")
		for rows.Next() {
			var station_id string
			var device_number int
			err := rows.Scan(&station_id, &device_number)
			if err != nil {
				log.Println(err)
			} else {
				// Get MongoDB Data And MYSQL Update
				fileter := bson.M{"$and": []bson.M{{"MeterValues.chargepointid": station_id}, {"MeterValues.payload.connectorid": device_number}}}
				option := options.Find()
				option.SetSort(bson.M{"Timestamp": -1})
				option.SetLimit(1)
				cursor, err := conn.Find(context.TODO(), fileter, option)
				if err != nil {
					log.Println(err)
				} else {
					for cursor.Next(context.TODO()) {
						var elem bson.M
						err := cursor.Decode(&elem)
						if err != nil {
							log.Println(err)
						} else {
							log.Println(elem)
						}
					}
				}
			}
		}

		// update도 MongoDB logging ?
		return
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

// RDB Charge Device 충전량 값 update(기존 값 + 현재 값) => stop transaction 때 (meterStop - meterStart) 만큼 작업하기
