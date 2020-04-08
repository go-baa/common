package ipip

import (
	"log"

	ipdb "github.com/ipipdotnet/ipdb-go"
)

var db *ipdb.City

func init() {
	var err error
	db, err = ipdb.NewCity("./data/mydata4vipweek2.ipdb")
	if err != nil {
		log.Fatal(err)
	}
}

// InCity 判断一个IP是不是特定城市
func InCity(ip string, citys ...string) bool {
	info, err := db.FindInfo(ip, "CN")
	if err == nil && info != nil {
		for _, city := range citys {
			if info.RegionName == city {
				return true
			}
		}
	}
	return false
}

// InBeijing 判断一个IP是不是北京的
func InBeijing(ip string) bool {
	info, err := db.FindInfo(ip, "CN")
	if err == nil && info != nil && info.RegionName == "北京" {
		return true
	}
	return false
}

// GetCity 获取一个IP的城市
func GetCity(ip string) (string, error) {
	info, err := db.FindInfo(ip, "CN")
	if err != nil {
		return "", err
	}
	return info.CityName, nil
}
