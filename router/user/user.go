package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Youn-sik/KSCONNECT_ADMIN/database"

	"github.com/bdwilliams/go-jsonify/jsonify"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

type ReqUser struct {
	Id       string
	Password string
}

type User struct {
	Uid      int
	Id       string
	Password string
	Name     string
	Email    string
	Mobile   string
}

type AuthTokenClaims struct {
	ID                 string `json:"id"`     // 유저 ID
	Name               string `json:"name"`   // 유저 이름
	Email              string `json:"mail"`   // 유저 메일
	Mobile             string `json:"mobile"` // 유저 메일
	UID                int    `json:"uid"`    // 유저 UID
	jwt.StandardClaims        // 표준 토큰 Claims
}

// 토큰 발급
func TokenBuild(u User) string {
	at := AuthTokenClaims{
		ID:     u.Id,
		UID:    u.Uid,
		Name:   u.Name,
		Mobile: u.Mobile,
		Email:  u.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Minute * 360)),
		},
	}

	atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, &at)
	signedAuthToken, err := atoken.SignedString([]byte("cho"))

	if err != nil {
		log.Println(err)
		return "false"
	}

	// log.Println(signedAuthToken)
	return signedAuthToken
}

func GetUserInfo() ([]string, error) {
	conn1 := database.NewMysqlConnection()
	defer conn1.Close()

	rows, err := conn1.Query("select id, password, name from admin_user")

	if err != nil {
		log.Println(err)
	}

	results := (jsonify.Jsonify(rows))

	return results, err
}

func TokenCheck(authToken string) bool {
	key := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			ErrUnexpectedSigningMethod := errors.New("unexpected signing method")
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte("cho"), nil
	}

	user := AuthTokenClaims{}
	token, err := jwt.ParseWithClaims(authToken, &user, key)

	if err != nil {
		// token is expired by ...
		log.Println(err)
		return false
	}

	return token.Valid
}

func Login(c *gin.Context) {
	var send_data struct {
		result string
		errStr string
	}
	reqData := ReqUser{}
	err := c.Bind(&reqData)

	user := User{}

	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "body 데이터를 가져오는중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	// log.Println(reqData)

	conn := database.NewMysqlConnection()
	rows, err := conn.Query("select uid, id, password, name, email, mobile from admin_user where id=?", reqData.Id)

	if err != nil {
		log.Println(err)
		send_data.result = "false"
		send_data.errStr = "DB Query 실행 중 문제가 발생하였습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	}

	jsonRows := jsonify.Jsonify(rows)
	log.Println(jsonRows)
	json.Unmarshal([]byte(jsonRows[0]), &user)

	// 비밀번호가 숫자로만 이루어져 있을 시, 형태를 int 형으로 가져가게 된다.
	// log.Println(user)
	// log.Println(user.Password)
	// log.Println(reqData.Password)

	if user.Password != reqData.Password {
		send_data.result = "false"
		send_data.errStr = "비밀번호가 동일하지 않습니다."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
	} else {
		token := TokenBuild(user)

		if token == "false" {
			send_data.result = "false"
			send_data.errStr = "토큰 생성 중 문제가 발생하였습니다."
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr})
		} else {
			send_data.result = "true"
			send_data.errStr = ""
			c.JSON(http.StatusOK, gin.H{"result": send_data.result, "errStr": send_data.errStr, "token": token})
		}
	}
}

func UserList(res http.ResponseWriter, req *http.Request) {
	user_list, err := GetUserInfo()
	if err != nil {
		log.Fatal(err)
		// user_list 값 대체 필요
		fmt.Fprint(res, "No Result About User Inf ormation")
	} else {
		fmt.Fprint(res, user_list)
	}
}
