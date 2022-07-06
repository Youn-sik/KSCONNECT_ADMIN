package natsclient

import (
	"context"
	"database/sql"
	"log"

	"github.com/Youn-sik/KSCONNECT_ADMIN/database"
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
		log.Println(err)
		return
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
		return
	}
}

func ResetDeviceUsage() {
	log.Println("---ResetDeviceUsageFunc---")
	MysqlClient = database.NewMysqlConnection()
	defer MysqlClient.Close()

	rows, err := MysqlClient.Query("select device_id, station_id, device_number, charge_device.usage from charge_device")
	if err != nil {
		log.Println(err)
		return
	}
	MongoClient = database.NewMongodbConnection()
	conn := MongoClient.Database("Admin_Service").Collection("service_device_usage_reset")
	conn1 := database.NewMysqlConnection()
	defer conn1.Close()
	for rows.Next() {
		var save_data struct {
			Device_id     int
			Station_id    int
			Device_number int
			Usage         int
		}
		err := rows.Scan(&save_data.Device_id, &save_data.Station_id, &save_data.Device_number, &save_data.Usage)
		if err != nil {
			log.Println(err)
			return
		}

		// Insert MongoDB Data
		_, err = conn.InsertOne(context.TODO(), bson.D{
			{Key: "Device_id", Value: save_data.Device_id},
			{Key: "Station_id", Value: save_data.Station_id},
			{Key: "Device_number", Value: save_data.Device_number},
			{Key: "Usage", Value: save_data.Usage},
		})
		if err != nil {
			log.Println(err)
			return
		}

		//Update MYSQL Data
		log.Printf("%+v\n", save_data)
		_, err = conn1.Query("update charge_device set charge_device.usage = 0 where device_id = ?", save_data.Device_id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	log.Println("Month Reset Complete")
}
