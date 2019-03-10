package amap

// key 默认配置的高德地图 kEY
const key = "d66e9f6d03180cfa2d54e6fd083bcf5c"
const gateway = "http://restapi.amap.com/v3"
const apiGeoCodeRe = "/geocode/regeo"

type response struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	InfoCode string `json:"infocode"`
}

// Point 一个地理位置的点
type Point struct {
	Latitude  string `json:"latitude"`  // 维度
	Longitude string `json:"longitude"` // 经度
}

func (t Point) String() string {
	return t.Longitude + "," + t.Latitude
}
