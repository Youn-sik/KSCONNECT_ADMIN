package natsclient

import (
	"log"
	"sync"

	"github.com/Youn-sik/KSCONNECT_ADMIN/database"
	"github.com/bdwilliams/go-jsonify/jsonify"
)

func NatsReply(subject string) {

	if subject == "ocpp/v16/chargepoints" {
		conn1 := database.NewMysqlConnection()
		defer conn1.Close()

		rows, err := conn1.Query("select station_id from charge_station")

		if err != nil {
			log.Fatal(err)
		}

		results := (jsonify.Jsonify(rows))

		nc := NewNatsClient()
		defer nc.Close()

		// Subscribe
		wg := sync.WaitGroup{}
		wg.Add(1)

		log.Println(results)
		nc.Reply("ocpp/v16/chargepoints", results, &wg)
		wg.Wait()
	} else if subject == "ocpp/v16/idtags" {
		conn1 := database.NewMysqlConnection()
		defer conn1.Close()

		rows, err := conn1.Query("select rfid from user")

		if err != nil {
			log.Fatal(err)
		}

		results := (jsonify.Jsonify(rows))

		nc := NewNatsClient()
		defer nc.Close()

		// Subscribe
		wg := sync.WaitGroup{}
		wg.Add(1)

		log.Println(results)
		nc.Reply("ocpp/v16/idtags", results, &wg)
		wg.Wait()
	}

}
