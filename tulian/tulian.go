package tulian

import (
	"encoding/json"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"log"
	"net/http"
	"oliujunk/input/api"
	"strconv"
	"time"
)

var (
	token     string
	deviceIDs = [...]int{
		16081799, 16081798, 16081797,
	}
)

func updateToken() {
	token = api.GetToken115("test", "123456")
}

func Start() {
	log.Println("图联数据转发")

	updateToken()

	job := cron.New(
		cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	_, _ = job.AddFunc("0 0 0/12 * * *", updateToken)
	_, _ = job.AddFunc("0 */5 * * * *", updateData)
	//_, _ = job.AddFunc("0 */1 * * * *", updateData)

	job.Start()

}

func updateData() {

	for _, deviceID := range deviceIDs {
		log.Println(deviceID)
		data := getLastData(deviceID)
		if data == nil {
			continue
		}
		result := api.PostData115(data)
		log.Println(result)

		time.Sleep(1 * time.Second)
	}
}

func getLastData(deviceID int) map[string]int {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "http://101.34.116.221:8005/data/"+strconv.Itoa(deviceID), nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	req.Header.Set("token", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	result, _ := ioutil.ReadAll(resp.Body)
	var currentData api.CurrentData
	_ = json.Unmarshal(result, &currentData)

	if currentData.Datatime != "" {
		now := time.Now()
		datatime, _ := time.Parse("2006-01-02 15:04:05", currentData.Datatime)
		if datatime.After(now.Add(-time.Minute * 6)) {
			data := map[string]int{"facId": deviceID,
				"e1":  currentData.E1,
				"e2":  currentData.E2,
				"e3":  currentData.E3,
				"e4":  currentData.E4,
				"e5":  currentData.E5,
				"e6":  currentData.E6,
				"e7":  currentData.E7,
				"e8":  currentData.E8,
				"e9":  currentData.E9,
				"e10": currentData.E10,
				"e11": currentData.E11,
				"e12": currentData.E12,
				"e13": currentData.E13,
				"e14": currentData.E14,
				"e15": currentData.E15,
				"e16": currentData.E16,
			}
			return data
		}
	}
	return nil
}
