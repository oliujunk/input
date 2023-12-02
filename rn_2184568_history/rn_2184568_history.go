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
	weathers = [...]int{20201794, 56291669, 56499189, 22018602, 22018603, 22018604, 22018605, 22018606, 22018607, 22018608}
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

	for _, device := range weathers {
		collection := client.Database("xph").Collection(strconv.Itoa(device))
		opts := options.FindOne().SetSort(bson.D{{"dataTime", -1}})
		var result api.CurrentData
		err := collection.FindOne(context.TODO(), bson.D{{"facId", device}}, opts).Decode(&result)
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
