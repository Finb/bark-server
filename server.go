package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/lithammer/shortuuid"

	"github.com/go-zoo/bone"
	bolt "go.etcd.io/bbolt"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
	"github.com/sideshow/apns2/payload"
)

type BaseResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func responseString(code int, message string) string {
	t, _ := json.Marshal(BaseResponse{Code: code, Message: message})
	return string(t)
}
func responseData(code int, data map[string]interface{}, message string) string {
	t, _ := json.Marshal(BaseResponse{Code: code, Data: data, Message: message})
	return string(t)
}

func ping(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()
	_, err := fmt.Fprint(w, responseData(200, map[string]interface{}{
		"version": Version,
		"build":   BuildDate,
		"arch":    runtime.GOOS + "/" + runtime.GOARCH,
		"commit":  CommitID,
	}, "pong"))
	if err != nil {
		logrus.Error(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {

	key := bone.GetValue(r, "key")
	category := bone.GetValue(r, "category")
	title := bone.GetValue(r, "title")
	body := bone.GetValue(r, "body")

	deviceToken, err := getDeviceTokenByKey(key)
	if err != nil {
		logrus.Errorf("找不到 key 对应的 DeviceToken key: %s", key)
		_, err = fmt.Fprint(w, responseString(400, "找不到 Key 对应的 DeviceToken, 请确保 Key 正确! Key 可在 App 端注册获得。"))
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	err = r.ParseForm()
	if err != nil {
		logrus.Error(err)
	}

	if len(title) <= 0 && len(body) <= 0 {
		//url中不包含 title body，则从Form里取
		for key, value := range r.Form {
			if strings.ToLower(key) == "title" {
				title = value[0]
			} else if strings.ToLower(key) == "body" {
				body = value[0]
			}
		}

	}

	if body == "" {
		body = "无推送文字内容"
	}

	params := make(map[string]interface{})
	for key, value := range r.Form {
		params[strings.ToLower(key)] = value[0]
	}

	logrus.Println(" ========================== ")
	logrus.Println("key: ", key)
	logrus.Println("category: ", category)
	logrus.Println("title: ", title)
	logrus.Println("body: ", body)
	logrus.Println("params: ", params)
	logrus.Println(" ========================== ")

	err = postPush(category, title, body, deviceToken, params)
	if err != nil {
		_, err = fmt.Fprint(w, responseString(400, err.Error()))
		if err != nil {
			logrus.Error(err)
		}
	} else {
		_, err = fmt.Fprint(w, responseString(200, ""))
		if err != nil {
			logrus.Error(err)
		}
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()
	err := r.ParseForm()
	if err != nil {
		logrus.Error(err)
	}

	key := shortuuid.New()
	var deviceToken string
	for key, value := range r.Form {
		if strings.ToLower(key) == "devicetoken" {
			deviceToken = value[0]
			break
		}
	}

	if len(deviceToken) <= 0 {
		_, err = fmt.Fprint(w, responseString(400, "deviceToken 不能为空"))
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	oldKey := r.FormValue("key")
	err = boltDB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("device"))
		if err != nil {
			return err
		}

		if len(oldKey) > 0 {
			//如果已经注册，则更新DeviceToken的值
			val := bucket.Get([]byte(oldKey))
			if val != nil {
				key = oldKey
			}
		}

		return bucket.Put([]byte(key), []byte(deviceToken))
	})

	if err != nil {
		_, err = fmt.Fprint(w, responseString(400, "注册设备失败"))
		if err != nil {
			logrus.Error(err)
		}
		return
	}
	logrus.Info("注册设备成功")
	logrus.Info("key: ", key)
	logrus.Info("deviceToken: ", deviceToken)
	_, err = fmt.Fprint(w, responseData(200, map[string]interface{}{"key": key}, "注册成功"))
	if err != nil {
		logrus.Error(err)
	}
}

func getDeviceTokenByKey(key string) (string, error) {
	var deviceTokenBytes []byte
	err := boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("device"))
		deviceTokenBytes = bucket.Get([]byte(key))
		if deviceTokenBytes == nil {
			return errors.New("没找到 DeviceToken")
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return string(deviceTokenBytes), nil
}

func getb() []byte {
	return []byte{45, 45, 45, 45, 45, 66, 69, 71, 73, 78, 32, 80, 82, 73, 86, 65, 84, 69, 32, 75, 69, 89, 45, 45, 45, 45, 45, 10, 77, 73, 71, 84, 65, 103, 69, 65, 77, 66, 77, 71, 66, 121, 113, 71, 83, 77, 52, 57, 65, 103, 69, 71, 67, 67, 113, 71, 83, 77, 52, 57, 65, 119, 69, 72, 66, 72, 107, 119, 100, 119, 73, 66, 65, 81, 81, 103, 52, 118, 116, 67, 51, 103, 53, 76, 53, 72, 103, 75, 71, 74, 50, 43, 10, 84, 49, 101, 65, 48, 116, 79, 105, 118, 82, 69, 118, 69, 65, 89, 50, 103, 43, 106, 117, 82, 88, 74, 107, 89, 76, 50, 103, 67, 103, 89, 73, 75, 111, 90, 73, 122, 106, 48, 68, 65, 81, 101, 104, 82, 65, 78, 67, 65, 65, 83, 109, 79, 115, 51, 74, 107, 83, 121, 111, 71, 69, 87, 90, 10, 115, 85, 71, 120, 70, 115, 47, 52, 112, 119, 49, 114, 73, 108, 83, 86, 50, 73, 67, 49, 57, 77, 56, 117, 51, 71, 53, 107, 113, 51, 54, 117, 112, 79, 119, 121, 70, 87, 106, 57, 71, 105, 51, 69, 106, 99, 57, 100, 51, 115, 67, 55, 43, 83, 72, 82, 113, 88, 114, 69, 65, 74, 111, 119, 10, 56, 47, 55, 116, 82, 112, 86, 43, 10, 45, 45, 45, 45, 45, 69, 78, 68, 32, 80, 82, 73, 86, 65, 84, 69, 32, 75, 69, 89, 45, 45, 45, 45, 45}
}

func postPush(category string, title string, body string, deviceToken string, params map[string]interface{}) error {

	notification := &apns2.Notification{}
	notification.DeviceToken = deviceToken
	var sound string = "1107"
	if params["sound"] != nil {
		sound = params["sound"].(string) + ".caf"
	}
	pl := payload.NewPayload().Sound(sound).Category("myNotificationCategory")
	badge := params["badge"]
	if badge != nil {
		badgeStr, pass := badge.(string)
		if pass {
			badgeNum, err := strconv.Atoi(badgeStr)
			if err == nil {
				pl = pl.Badge(badgeNum)
			}
		}
	}

	for key, value := range params {
		pl = pl.Custom(key, value)
	}
	if len(title) > 0 {
		pl.AlertTitle(title)
	}
	if len(body) > 0 {
		pl.AlertBody(body)
	}
	notification.Payload = pl.MutableContent()
	notification.Topic = "me.fin.bark"
	res, err := apnsClient.Push(notification)

	if err != nil {
		logrus.Errorf("Error:", err)
		return fmt.Errorf("与苹果推送服务器传输数据失败: %w", err)
	}
	logrus.Infof("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	if res.StatusCode == 200 {
		return nil
	} else {
		return errors.New("推送发送失败 " + res.Reason)
	}

}

var boltDB *bolt.DB
var apnsClient *apns2.Client

func runBarkServer() {
	//f,_:= os.Open("./BarkPush.p12")
	//t,_ := ioutil.ReadAll(f)
	//
	//str := ""
	//for _,val := range t {
	//	str += ", "
	//	str += strconv.Itoa(int(val))
	//}
	//
	//fmt.Printf(string(t))

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err = os.Mkdir(dataDir, 0755); err != nil {
			logrus.Fatal(err)
		}
	} else if err != nil {
		logrus.Fatal(err)
	}

	db, err := bolt.Open(filepath.Join(dataDir, "bark.db"), 0600, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	boltDB = db

	err = boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("device"))
		return err
	})
	if err != nil {
		logrus.Fatal(err)
	}

	authKey, err := token.AuthKeyFromBytes(getb())
	if err != nil {
		logrus.Fatalf("token error:", err)
	}
	clientToken := &token.Token{
		AuthKey: authKey,
		KeyID:   "LH4T9V5U4R",
		TeamID:  "5U8LBRXG3A",
	}

	apnsClient = apns2.NewTokenClient(clientToken).Production()

	addr := fmt.Sprint(listenAddr, ":", listenPort)
	logrus.Info("Serving HTTP on " + addr)

	r := bone.New()
	r.Get("/ping", http.HandlerFunc(ping))
	r.Post("/ping", http.HandlerFunc(ping))

	r.Get("/register", http.HandlerFunc(register))
	r.Post("/register", http.HandlerFunc(register))

	r.Get("/:key/:body", http.HandlerFunc(index))
	r.Post("/:key/:body", http.HandlerFunc(index))

	r.Get("/:key/:title/:body", http.HandlerFunc(index))
	r.Post("/:key/:title/:body", http.HandlerFunc(index))

	r.Get("/:key/:category/:title/:body", http.HandlerFunc(index))
	r.Post("/:key/:category/:title/:body", http.HandlerFunc(index))

	err = http.ListenAndServe(addr, r)
	if err != nil {
		logrus.Fatal(err)
	}
}
