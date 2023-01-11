package soilmoisture

import (
	"bytes"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"oliujunk/input/api"
	"strconv"
	"strings"
	"time"
)

var (
	baseURL    = "http://swc.cau-iae.cn:8080"
	userID     string
	deviceList []DeviceInfo
)

type DeviceList struct {
	IsSuccess bool         `json:"IsSuccess"`
	DevList   []DeviceInfo `json:"devList"`
}

type DeviceInfo struct {
	AliasName   string  `json:"AliasName"`
	DevID       string  `json:"devID"`
	DevType     string  `json:"devType"`
	LAT         float32 `json:"LAT"`
	LNG         float32 `json:"LNG"`
	DevActiTime string  `json:"devActiTime"`
	DevTermTime string  `json:"devTermTime"`
}

type LastData struct {
	IsSuccess   bool          `json:"IsSuccess"`
	DevLastData []DevLastData `json:"devLastData"`
}

type DevLastData struct {
	POWER       int     `json:"POWER"`       // 电量
	TIME        string  `json:"TIME"`        // 更新时间
	CSQ         int     `json:"CSQ"`         // 信号强度
	DataTEMPStr string  `json:"DataTEMPStr"` // 不同深度温度
	DataSWCStr  string  `json:"DataSWCStr"`  // 不同深度湿度
	DataATM     int     `json:"DataATM"`     // 大气压强
	DataAT      float32 `json:"DataAT"`      // 空气温度
	DataATS     float32 `json:"DataATS"`     // 空气湿度
	DevID       string  `json:"devID"`       // 设备ID号
}

func Start() {
	log.Println("墒情API上报")

	userID = getUserID("武汉睿农科技有限公司", "123qwe")
	deviceList = getDeviceList(userID)

	job := cron.New(
		cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	_, _ = job.AddFunc("0 0 */1 * * *", updateData)
	//_, _ = job.AddFunc("0 */1 * * * *", updateData)

	job.Start()

}

func getUserID(name string, pwd string) string {
	urlValue := url.Values{}
	urlValue.Add("name", name)
	urlValue.Add("pwd", pwd)
	payload := strings.NewReader(urlValue.Encode())
	req, err := http.NewRequest("POST", baseURL+"/getID", payload)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	buf := bytes.NewBuffer(result)
	message, err := simplejson.NewFromReader(buf)
	if err != nil {
		log.Fatalln(err)
	}
	userID, _ := message.Get("UserID").String()
	return userID
}

func getDeviceList(userID string) []DeviceInfo {
	urlValue := url.Values{}
	urlValue.Add("userID", userID)
	payload := strings.NewReader(urlValue.Encode())
	req, err := http.NewRequest("POST", baseURL+"/getDevList", payload)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var deviceList DeviceList
	_ = json.Unmarshal(result, &deviceList)
	return deviceList.DevList
}

func getLastData(userID string, devID string) map[string]int {
	urlValue := url.Values{}
	urlValue.Add("userID", userID)
	urlValue.Add("devID", devID)
	payload := strings.NewReader(urlValue.Encode())
	req, err := http.NewRequest("POST", baseURL+"/getDevLastData", payload)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var lastData LastData
	err = json.Unmarshal(result, &lastData)
	if err != nil {
		log.Println(err)
		return nil
	}

	deviceID, _ := strconv.Atoi(devID[len(devID)-8:])
	dataArray := []int{0, 0, 0, 0, 0, 0}

	humidityStr := strings.Split(lastData.DevLastData[0].DataSWCStr, ",")
	temperatureStr := strings.Split(lastData.DevLastData[0].DataTEMPStr, ",")

	if len(humidityStr) <= 1 {
		return nil
	}

	tmpData, _ := strconv.ParseFloat(humidityStr[0], 32)
	dataArray[0] = int(tmpData * 100)
	tmpData, _ = strconv.ParseFloat(humidityStr[1], 32)
	dataArray[1] = int(tmpData * 100)
	tmpData, _ = strconv.ParseFloat(humidityStr[2], 32)
	dataArray[2] = int(tmpData * 100)

	tmpData, _ = strconv.ParseFloat(temperatureStr[0], 32)
	dataArray[3] = int(tmpData * 100)
	tmpData, _ = strconv.ParseFloat(temperatureStr[1], 32)
	dataArray[4] = int(tmpData * 100)
	tmpData, _ = strconv.ParseFloat(temperatureStr[2], 32)
	dataArray[5] = int(tmpData * 100)

	data := map[string]int{"facId": deviceID,
		"e1":  lastData.DevLastData[0].POWER,
		"e2":  dataArray[0],
		"e3":  dataArray[1],
		"e4":  dataArray[2],
		"e5":  dataArray[3],
		"e6":  dataArray[4],
		"e7":  dataArray[5],
		"e8":  lastData.DevLastData[0].CSQ,
		"e9":  int(lastData.DevLastData[0].DataAT * 10),
		"e10": int(lastData.DevLastData[0].DataATS * 10),
		"e11": lastData.DevLastData[0].DataATM,
		"e12": 0, "e13": 0, "e14": 0, "e15": 0, "e16": 0,
	}
	return data
}

func updateData() {

	for _, device := range deviceList {
		log.Println(device.DevID)
		data := getLastData(userID, device.DevID)
		if data == nil {
			continue
		}
		result := api.PostData(data)
		log.Println(result)

		time.Sleep(1 * time.Second)
	}
}
