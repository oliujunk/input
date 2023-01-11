package tiannankeji

import (
	"bytes"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	imeis = [...]string{"G14547781", "G95799357",
		"K09576986", "K09576994",
		"K09577008", "K09577011",
		"K09576966", "K09576970",
		"K09577010", "K09576956",
		"K55779945", "K68778231",
		"K68778241",
	}
	uids = [...]string{"31FFD9055642343035621343", "31FFD2054848313023811743",
		"19004300190000384753524E", "490036000F0000384753524E",
		"480036000F0000384753524E", "470040000F0000384753524E",
		"510029000F0000384753524E", "09001900140000384753524E",
		"29000900130000384753524E", "470043000F0000384753524E",
		"3E005F001238314D33504E53", "26005D001238314D33504E53",
		"2A005E001238314D33504E53",
	}
	usernames = [...]string{"田南科技"}
	passwords = [...]string{"888888"}
)

type PestData struct {
	Imei     string `json:"imei"`
	DataTime string `json:"dataTime"`
	E1       int    `json:"e1"`
	E2       int    `json:"e2"`
	E3       int    `json:"e3"`
	E4       int    `json:"e4"`
	E5       int    `json:"e5"`
	E6       int    `json:"e6"`
	E7       int    `json:"e7"`
	E8       int    `json:"e8"`
	E9       int    `json:"e9"`
	E10      int    `json:"e10"`
	E11      int    `json:"e11"`
	E12      int    `json:"e12"`
	E13      int    `json:"e13"`
	E14      int    `json:"e14"`
	E15      int    `json:"e15"`
	E16      int    `json:"e16"`
	E17      int    `json:"e17"`
	E18      int    `json:"e18"`
	E19      int    `json:"e19"`
	E20      int    `json:"e20"`
	E21      int    `json:"e21"`
	E22      int    `json:"e22"`
	E23      int    `json:"e23"`
	E24      int    `json:"e24"`
	E25      int    `json:"e25"`
	E26      int    `json:"e26"`
	E27      int    `json:"e27"`
	E28      int    `json:"e28"`
	E29      int    `json:"e29"`
	E30      int    `json:"e30"`
	E31      int    `json:"e31"`
	E32      int    `json:"e32"`
}

func Start() {
	log.Print("比昂虫情设备数据更新")

	job := cron.New(
		cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	_, _ = job.AddFunc("0 0 */1 * * *", getData)
	//_, _ = job.AddFunc("0 */2 * * * *", getData)
	_, _ = job.AddFunc("00 50 23 * * *", getPicture)
	//_, _ = job.AddFunc("0 */2 * * * *", getPicture)
	job.Start()
}

func getData() {
	for index, imei := range imeis {
		params := url.Values{}
		Url, _ := url.Parse("http://open.cdbeyond.com/open/Ycc2GetCurrData")
		params.Set("username", usernames[0])
		params.Set("password", passwords[0])
		params.Set("did", uids[index])
		Url.RawQuery = params.Encode()
		urlPath := Url.String()
		resp, err := http.Get(urlPath)
		if err != nil {
			log.Panicln(err)
		}
		result, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(result))
		buf := bytes.NewBuffer(result)
		message, err := simplejson.NewFromReader(buf)
		if err != nil {
			log.Panicln(err)
		}
		pestData := PestData{}
		pestData.Imei = imei
		stamp := message.Get("data").Get("time").MustInt64()
		dataTime := time.Unix(stamp/1000, 0)
		pestData.DataTime = dataTime.Format("2006-01-02 15:04:05")
		pestData.E1 = int(message.Get("data").Get("humidity").MustFloat64() * 10)
		pestData.E2 = int(message.Get("data").Get("temp").MustFloat64() * 10)
		pestData.E3 = message.Get("data").Get("lux").MustInt()
		pestData.E4 = int(message.Get("data").Get("lat").MustFloat64() * 10000)
		pestData.E5 = int(message.Get("data").Get("lng").MustFloat64() * 10000)
		pestData.E6 = message.Get("data").Get("num").MustInt()

		client := &http.Client{Timeout: 5 * time.Second}
		jsonStr, _ := json.Marshal(pestData)
		log.Println(string(jsonStr))
		resp, err = client.Post("http://101.34.116.221:8005/pest/data", "application/json", bytes.NewBuffer(jsonStr))
		if err != nil {
			log.Panicln(err)
		}

		result, _ = ioutil.ReadAll(resp.Body)
		log.Println(string(result))

		time.Sleep(5 * time.Second)
	}
}

func getPicture() {
	for _, imei := range imeis {
		params := url.Values{}
		Url, _ := url.Parse("http://open.cdbeyond.com/open/getTelemeter")
		params.Set("username", usernames[0])
		params.Set("password", passwords[0])
		params.Set("imei", imei)
		Url.RawQuery = params.Encode()
		urlPath := Url.String()
		resp, err := http.Get(urlPath)
		if err != nil {
			log.Panicln(err)
		}
		result, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(result))
		buf := bytes.NewBuffer(result)
		message, err := simplejson.NewFromReader(buf)
		if err != nil {
			log.Panicln(err)
		}
		images := message.Get("data").MustArray()
		for index := range images {
			image := message.Get("data").GetIndex(index)
			stamp, _ := image.Get("photoTime").Int64()
			photoTime := time.Unix(stamp/1000, 0)
			params := url.Values{}
			Url, _ := url.Parse("http://101.34.116.221:8005/pest/image/url")
			params.Set("imei", imei)
			params.Set("url", image.Get("picUrl").MustString())
			params.Set("time", photoTime.Format("2006-01-02 15:04:05"))
			log.Println(params)
			Url.RawQuery = params.Encode()
			urlPath := Url.String()
			resp, err := http.Post(urlPath, "application/json", nil)
			if err != nil {
				log.Panicln(err)
			}
			result, _ := ioutil.ReadAll(resp.Body)
			log.Println(string(result))
			time.Sleep(10 * time.Second)
		}
		time.Sleep(10 * time.Second)
	}
}
