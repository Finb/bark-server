package main

import (
	"net/http"
	"fmt"
	"log"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"

	"github.com/boltdb/bolt"

	"errors"
	"encoding/json"
	"github.com/go-zoo/bone"
	"github.com/renstrom/shortuuid"
	"strconv"
	"flag"
	"strings"
)

type BaseResponse struct {
	Code    int         `json:"code"`
	Data interface{} `json:"data"`
	Message string      `json:"message"`
}

func responseString(code int, message string)string {
	t, _ := json.Marshal(BaseResponse{Code:code,Message:message})
	return string(t)
}
func responseData(code int, data map[string]interface{} , message string)string {
	t, _ := json.Marshal(BaseResponse{Code:code,Data:data, Message:message})
	return string(t)
}

func ping(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Fprint(w, responseData(200, map[string]interface{}{"version": "1.0.0"},"pong"))
}

func Index(w http.ResponseWriter, r *http.Request) {
	key := bone.GetValue(r, "key")
	category := bone.GetValue(r, "category")
	title := bone.GetValue(r, "title")
	body := bone.GetValue(r, "body")

	defer r.Body.Close()

	deviceToken ,err := getDeviceTokenByKey(key)
	if err != nil {
		log.Println("找不到key对应的DeviceToken key: " + key)
		fmt.Fprint(w, responseString(400,"找不到key对应的DeviceToken, 请确保Key正确! Key可在App端注册获得。"))
		return
	}

	r.ParseForm()

	if len(title) <= 0 && len(body) <= 0 {
		//url中不包含 title body，则从Form里取
		for key,value := range r.Form{
			if strings.ToLower(key) == "title" {
				title = value[0]
			} else if strings.ToLower(key) == "body"{
				body = value[0]
			}
		}


	}

	if len(body) <= 0 {
		body = "无推送文字内容"
	}

	params := make(map[string]interface{})
	for key,value := range r.Form {
		params[strings.ToLower(key)] = value[0]
	}

	log.Println(" ========================== ")
	log.Println("key: ", key)
	log.Println("category: ", category)
	log.Println("title: ", title)
	log.Println("body: ", body)
	log.Println("params: ", params)
	log.Println(" ========================== ")

	err = postPush(category,title,body,deviceToken,params)
	if err != nil {
		fmt.Fprint(w, responseString(400, err.Error()))
	} else{
		fmt.Fprint(w, responseString(200, ""))
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.ParseForm()
	key := shortuuid.New()
	var deviceToken string
	for key,value := range r.Form {
		if strings.ToLower(key) == "devicetoken" {
			deviceToken = value[0]
			break
		}
	}

	if len(deviceToken) <= 0 {
		fmt.Fprint(w, responseString(400, "deviceToken 不能为空"))
		return
	}

	oldKey := r.FormValue("key")
	boltDB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("device"))
		if err != nil {
			return  err
		}

		if len(oldKey) >0 {
			//如果已经注册，则更新DeviceToken的值
			val := bucket.Get([]byte(oldKey))
			if val != nil {
				key = oldKey
			}
		}

		bucket.Put([]byte(key), []byte(deviceToken))
		return nil
	})
	log.Println("注册设备成功")
	log.Println("key: ", key)
	log.Println("deviceToken: ", deviceToken)
	fmt.Fprint(w, responseData(200, map[string]interface{}{"key":key}, "注册成功"))
}


func getDeviceTokenByKey(key string) (string,error){
	var deviceTokenBytes []byte
	err := boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("device"))
		deviceTokenBytes = bucket.Get([]byte(key))
		if deviceTokenBytes == nil {
			return errors.New("没找到DeviceToken")
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return string(deviceTokenBytes), nil
}

func getb() []byte {
	//测试证书
	if IsDev{
		return []byte{}
	}
	//线上证书
	return []byte{}
}

func postPush(category string, title string, body string, deviceToken string,  params map[string]interface{}) error{

	notification := &apns2.Notification{}
	notification.DeviceToken = deviceToken

	payload := payload.NewPayload().Sound("1107").Category("myNotificationCategory")
	badge := params["badge"]
	if badge != nil {
		badgeStr, pass := badge.(string)
		if pass {
			badgeNum, err := strconv.Atoi(badgeStr)
			if err == nil {
				payload = payload.Badge(badgeNum)
			}
		}
	}

	for key, value := range params {
		payload = payload.Custom(key, value)
	}
	if len(title) > 0 {
		payload.AlertTitle(title)
	}
	if len(body) > 0 {
		payload.AlertBody(body)
	}
	notification.Payload = payload
	notification.Topic = "me.fin.bark"
	res, err := apnsClient.Push(notification)

	if err != nil {
		log.Println("Error:", err)
		return errors.New("与苹果推送服务器传输数据失败")
	}
	log.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	if res.StatusCode == 200 {
		return nil
	}else{
		return errors.New("推送发送失败 " + res.Reason)
	}


}

var IsDev bool = false
var boltDB *bolt.DB
var apnsClient *apns2.Client
func main()  {
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
	ip := flag.String("ip",  "0.0.0.0", "http listen ip")
	port := flag.Int("port",  8080, "http listen port")
	dev := flag.Bool("dev", false, "develop推送，请忽略此参数，设置此参数为True会导致推送失败")
	flag.Parse()

	IsDev = *dev

	db, err := bolt.Open("bark.db", 0600, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer  db.Close()
	boltDB = db

	boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("device"))
		if err != nil {
			log.Fatalln(err)
		}
		return err
	})

	cert, err := certificate.FromP12Bytes(getb(),"bp")
	if err != nil {
		log.Fatalln("cer error")
	}
	apnsClient = apns2.NewClient(cert).Production()



	addr := *ip + ":" + strconv.Itoa(*port)
	log.Println("Serving HTTP on " + addr)

	r := bone.New()
	r.Get("/ping", http.HandlerFunc(ping))
	r.Post("/ping", http.HandlerFunc(ping))

	r.Get("/register", http.HandlerFunc(register))
	r.Post("/register", http.HandlerFunc(register))

	r.Get("/:key/:body", http.HandlerFunc(Index))
	r.Post("/:key/:body", http.HandlerFunc(Index))

	r.Get("/:key/:title/:body", http.HandlerFunc(Index))
	r.Post("/:key/:title/:body", http.HandlerFunc(Index))

	r.Get("/:key/:category/:title/:body", http.HandlerFunc(Index))
	r.Post("/:key/:category/:title/:body", http.HandlerFunc(Index))


	log.Fatal(http.ListenAndServe(addr, r))
}

