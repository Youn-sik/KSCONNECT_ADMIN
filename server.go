package main

import (
	// "context"
	// "encoding/json"
	// "fmt"
	"log"
	"net/http"

	// "net/http"

	// "KSCONNECT_ADMIN/database"
	// "KSCONNECT_ADMIN/setting"
	"github.com/Youn-sik/KSCONNECT_ADMIN/router/user"

	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
)

// func GetTestOne() {
// 	config, err := setting.LoadConfigSettingJSON()

// 	if err != nil {
// 		log.Fatal(err)
// 		panic(err)
// 	}

// 	conn, err := database.NewMongodbConnection()

// 	if err != nil {
// 		log.Fatal(err)
// 		panic(err)
// 	}

// 	mongoColl := conn.Database(config.Mongodb.Database).Collection("test")

// 	var result bson.M
// 	err = mongoColl.FindOne(context.TODO(), bson.D{{"item", "card"}}).Decode(&result)

// 	if err == mongo.ErrNoDocuments {
// 		log.Println("No document was found with the item: 'card'")
// 		return
// 	}
// 	if err != nil {
// 		log.Fatal(err)
// 		panic(err)
// 	}

// 	log.Println(result)
// }

// func GetTest() []string {
// 	var results []string
// 	config, err := setting.LoadConfigSettingJSON()

// 	if err != nil {
// 		log.Fatal(err)
// 		panic(err)
// 	}
// 	conn, err := database.NewMongodbConnection()

// 	if err != nil {
// 		log.Fatal(err)
// 		panic(err)
// 	}

// 	mongoColl := conn.Database(config.Mongodb.Database).Collection("test")

// 	filter := bson.D{{"item", "card"}}
// 	cursor, err := mongoColl.Find(context.TODO(), filter)

// 	if err == mongo.ErrNoDocuments {
// 		log.Println("No document was found with the item: 'card'")
// 		return results // empty array needed
// 	}
// 	if err != nil {
// 		log.Fatal(err)
// 		panic(err)
// 	}

// 	for cursor.Next(context.TODO()) {
// 		var number bson.M
// 		err := cursor.Decode(&number)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		b, _ := json.Marshal(number)
// 		// log.Println(string(b))
// 		results = append(results, string(b))
// 		// results = append(results, string(b)+",")
// 	}

// 	// log.Println(results)
// 	return results
// }

// func UserList(res http.ResponseWriter, req *http.Request) {
// 	user_list, err := GetUserInfo()
// 	if err != nil {
// 		log.Fatal(err)
// 		// user_list 값 대체 필요
// 		fmt.Fprint(res, "No Result About User Inf ormation")
// 	} else {
// 		fmt.Fprint(res, user_list)
// 	}
// }

// func TestList(res http.ResponseWriter, req *http.Request) {
// 	test_list := GetTest()
// 	fmt.Fprint(res, test_list)
// }

func authenticateMiddleware(c *gin.Context) {
	authToken := c.Request.Header.Get("authorization")

	if authToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"result": "false", "errStr": "No Token"})
		c.Abort()
		return
	}

	// token 검증 로직
	isValid := user.TokenCheck(authToken)
	if isValid == true {
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"result": "false", "errStr": "Expired Token"})
		c.Abort()
		return
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	// 그룹화 하여 관리 => 토큰 미들웨어 사용에 따라 그룹하기
	NAuth := router.Group("/NAuth")
	NAuth.POST("/user/login", func(c *gin.Context) {
		user.Login(c)
	})

	router.Use(authenticateMiddleware)
	router.GET("/", func(c *gin.Context) {
		c.String(200, "pongggg")
	})

	return router
}

func main() {
	var port string = ":4001"

	router := setupRouter()
	router.Run(port)
	log.Println("[SERVER] => Backend Admin application is listening on port " + port)
}
