package rn_2184568_history

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"oliujunk/input/api"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	devices []api.Device

	token    string
	weathers = [...]int{16069548, 16069547, 16069549}
	waters   = [...]int{16069551, 16069554, 16069557, 16069561, 16069562, 16069569, 16069571}
	soils    = [...]int{16069586, 16069583, 16069584, 16069585, 16069587, 16069589, 16069592, 16069594, 16069595, 16069597, 16069600, 16069601, 16069602}
)

func updateToken() {
	token = api.GetToken("2184568", "88888888")
}

func updateDevices() {
	devices = api.GetDevices("2184568", token)
}

func Start() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://root:13387543885@101.34.116.221:27017"))
	if err != nil {
		log.Println(err)
	}
	log.Println("2184568历史数据补充")
	updateToken()
	updateDevices()

	for _, device := range devices {
		log.Println(device.DeviceID)
		if device.DeviceID != 16069573 {
			collection := client.Database("xph").Collection(strconv.Itoa(device.DeviceID))
			opts := options.FindOne().SetSort(bson.D{{"dataTime", -1}})
			var result api.CurrentData
			err := collection.FindOne(context.TODO(), bson.D{{"facId", device.DeviceID}}, opts).Decode(&result)
			if err != nil {
				continue
			}
			//fmt.Printf("found document %v", result)
			start, err := time.Parse("2006-01-02 15:04:05", "2023-04-01 00:00:00")
			end, err := time.Parse("2006-01-02 15:04:05", "2021-05-01 00:00:00")
			if err != nil {
				continue
			}
			for start.After(end) {
				start = start.Add(-60 * time.Minute)
				result.Datatime = start.Format("2006-01-02 15:04:05")
				one, err := collection.InsertOne(context.TODO(), result)
				if err != nil {
					continue
				}
				fmt.Println(result, one)
			}

			time.Sleep(1 * time.Second)
		}
	}
}
