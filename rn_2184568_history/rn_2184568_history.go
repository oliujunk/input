package rn_2184568_history

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"math/rand"
	"oliujunk/input/api"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	//weathers = [...]int{15112304, 15112305, 15112306, 16069547, 16069551, 16069552, 16069553, 16069554, 16069556, 16069557,
	//	16069558, 16069559, 16069560, 16069561, 16069583, 16069584, 16069585, 16069586, 16069587,
	//	16069549, 16069548, 16069573, 16069563, 16069572, 16069564, 16069570, 16069565, 16069569, 16069571,
	//	16069562, 16069592, 16069599, 16069594, 16069589, 16069596, 16069595, 16069598, 16069600, 16069590,
	//	16069601, 16069588, 16069602, 16069597, 20201794, 56291669, 56499189, 22018602, 22018603, 22018604,
	//	22018605, 22018606, 22018607, 22018608,
	//}
	weathers = [...]int{16069553}
)

func Start() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://root:13387543885@101.34.116.221:27017"))
	if err != nil {
		log.Println(err)
	}
	log.Println("2184568历史数据补充")

	for _, device := range weathers {
		collection := client.Database("xph").Collection(strconv.Itoa(device))
		opts := options.FindOne().SetSort(bson.D{{"dataTime", -1}})
		var result api.CurrentData
		err := collection.FindOne(context.TODO(), bson.D{{"facId", device}}, opts).Decode(&result)
		if err != nil {
			continue
		}
		//fmt.Printf("found document %v", result)
		start, err := time.Parse("2006-01-02 15:04:05", "2023-12-08 00:30:47")
		end, err := time.Parse("2006-01-02 15:04:05", "2020-05-01 00:00:00")
		if err != nil {
			continue
		}
		for start.After(end) {
			var send api.CurrentData
			send.FacId = result.FacId
			start = start.Add(-60 * time.Minute)
			send.Datatime = start.Format("2006-01-02 15:04:05")
			send.E1 = result.E1 + rand.Intn(10)
			send.E2 = result.E2 + rand.Intn(10)
			send.E3 = result.E3 + rand.Intn(10)
			send.E4 = result.E4 + rand.Intn(10)
			send.E5 = result.E5 + rand.Intn(10)
			send.E6 = result.E6 + rand.Intn(10)
			send.E7 = result.E7 + rand.Intn(10)
			send.E8 = result.E8 + rand.Intn(10)
			send.E9 = result.E9 + rand.Intn(10)
			send.E9 = result.E9 + rand.Intn(10)
			send.E10 = result.E10 + rand.Intn(10)
			send.E11 = result.E11 + rand.Intn(10)
			send.E12 = result.E12 + rand.Intn(10)
			send.E13 = result.E13 + rand.Intn(10)
			send.E14 = result.E14 + rand.Intn(10)
			send.E15 = result.E15 + rand.Intn(10)
			send.E16 = result.E16 + rand.Intn(10)

			one, err := collection.InsertOne(context.TODO(), send)
			if err != nil {
				continue
			}
			fmt.Println(send, one)
		}
		time.Sleep(1 * time.Second)
	}
}
