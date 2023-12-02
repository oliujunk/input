package rn_2184568

import (
	"encoding/json"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"math/rand"
	"net/http"
	"oliujunk/input/api"
	"strconv"
	"time"
)

var (
	token    string
	weathers = [...]int{16069548, 16069547, 16069549}
	waters   = [...]int{16069551, 16069554, 16069557, 16069561, 16069562, 16069569, 16069571}
	soils    = [...]int{16069586, 16069583, 16069584, 16069585, 16069587, 16069589, 16069592, 16069594, 16069595, 16069597, 16069600, 16069601, 16069602}
)

func updateToken() {
	token = api.GetToken("2184568", "88888888")
}

func Start() {
	log.Println("2184568离线设备数据补充")
	updateToken()
	job := cron.New(
		cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	_, _ = job.AddFunc("0 0 0/12 * * *", updateToken)
	_, _ = job.AddFunc("0 */30 * * * *", updateData)

	job.Start()

}

func updateData() {

	// 气象站
	originalWeatherData := getLastData(weathers[0])
	for i := 1; i < len(weathers); i++ {
		weatherData := map[string]int{"facId": weathers[i],
			"e1":  originalWeatherData["e1"],
			"e2":  originalWeatherData["e2"],
			"e3":  originalWeatherData["e3"],
			"e4":  originalWeatherData["e4"],
			"e5":  originalWeatherData["e5"],
			"e6":  originalWeatherData["e6"],
			"e7":  originalWeatherData["e7"],
			"e8":  originalWeatherData["e8"],
			"e9":  originalWeatherData["e9"],
			"e10": originalWeatherData["e10"],
			"e11": originalWeatherData["e11"],
			"e12": originalWeatherData["e12"],
			"e13": originalWeatherData["e13"],
			"e14": originalWeatherData["e14"],
			"e15": originalWeatherData["e15"],
			"e16": originalWeatherData["e16"],
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		weatherData["e1"] = weatherData["e1"] + r.Intn(10)
		weatherData["e5"] = weatherData["e5"] + r.Intn(20)
		weatherData["e7"] = weatherData["e7"] + r.Intn(10)
		weatherData["e10"] = weatherData["e10"] + r.Intn(20)
		weatherData["e11"] = weatherData["e11"] + r.Intn(20)
		result := api.PostData(weatherData)
		log.Println(weathers[i], result)

		time.Sleep(1 * time.Second)
	}

	// 水位监测
	originalWaterData := getLastData(waters[0])
	for i := 1; i < len(waters); i++ {
		waterData := map[string]int{"facId": waters[i],
			"e1":  originalWaterData["e1"],
			"e2":  originalWaterData["e2"],
			"e3":  originalWaterData["e3"],
			"e4":  originalWaterData["e4"],
			"e5":  originalWaterData["e5"],
			"e6":  originalWaterData["e6"],
			"e7":  originalWaterData["e7"],
			"e8":  originalWaterData["e8"],
			"e9":  originalWaterData["e9"],
			"e10": originalWaterData["e10"],
			"e11": originalWaterData["e11"],
			"e12": originalWaterData["e12"],
			"e13": originalWaterData["e13"],
			"e14": originalWaterData["e14"],
			"e15": originalWaterData["e15"],
			"e16": originalWaterData["e16"],
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		waterData["e2"] = waterData["e2"] + r.Intn(500)
		result := api.PostData(waterData)
		log.Println(waters[i], result)
		time.Sleep(1 * time.Second)
	}

	// 墒情监测
	originalSoilData := getLastData(soils[0])
	for i := 1; i < len(soils); i++ {
		soilData := map[string]int{"facId": soils[i],
			"e1":  originalSoilData["e1"],
			"e2":  originalSoilData["e2"],
			"e3":  originalSoilData["e3"],
			"e4":  originalSoilData["e4"],
			"e5":  originalSoilData["e5"],
			"e6":  originalSoilData["e6"],
			"e7":  originalSoilData["e7"],
			"e8":  originalSoilData["e8"],
			"e9":  originalSoilData["e9"],
			"e10": originalSoilData["e10"],
			"e11": originalSoilData["e11"],
			"e12": originalSoilData["e12"],
			"e13": originalSoilData["e13"],
			"e14": originalSoilData["e14"],
			"e15": originalSoilData["e15"],
			"e16": originalSoilData["e16"],
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		soilData["e1"] = soilData["e1"] + r.Intn(10)
		soilData["e2"] = soilData["e2"] + r.Intn(50)
		soilData["e3"] = soilData["e3"] + r.Intn(50)
		soilData["e4"] = soilData["e4"] + r.Intn(10)
		soilData["e5"] = soilData["e5"] + r.Intn(20)
		soilData["e6"] = soilData["e6"] + r.Intn(50)
		soilData["e11"] = soilData["e11"] + r.Intn(50)
		soilData["e14"] = soilData["e14"] + r.Intn(20)
		soilData["e15"] = soilData["e15"] + r.Intn(20)
		result := api.PostData(soilData)
		log.Println(soils[i], result)
		time.Sleep(1 * time.Second)
	}

}

func getLastData(deviceID int) map[string]int {
	client := &http.Client{Timeout: 15 * time.Second}
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
	result, _ := io.ReadAll(resp.Body)
	var currentData api.CurrentData
	_ = json.Unmarshal(result, &currentData)

	if currentData.Datatime != "" {
		now := time.Now()
		datatime, _ := time.Parse("2006-01-02 15:04:05", currentData.Datatime)
		if datatime.After(now.Add(-time.Minute * 60)) {
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
