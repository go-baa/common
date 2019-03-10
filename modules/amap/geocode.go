package amap

import (
	"encoding/json"
	"strings"

	"git.code.tencent.com/xinhuameiyu/common/util"
)

// RegeoParams 逆地理编码 参数
type RegeoParams struct {
	Location  Point    // 传入内容规则：经度在前，纬度在后，经纬度间以“,”分割，经纬度小数点后不要超过 6 位。如果需要解析多个经纬度的话，请用"|"进行间隔，并且将 batch 参数设置为 true，最多支持传入 20 对坐标点。每对点坐标之间用"|"分割。
	PoiType   []string // 支持传入POI TYPECODE及名称；支持传入多个POI类型，多值间用“|”分隔
	Radius    int      // 查询POI的半径范围。取值范围：0~3000,单位：米
	RoadLevel int      // 可选值：0|1，当roadlevel=0时，显示所有道路，当roadlevel=1时，过滤非主干道路，仅输出主干道路数据
}

// RegeoResponse 逆地理编码 接口返回
type RegeoResponse struct {
	response
	RegeoCode *RegeoCode `json:"regeocode"`
}

// RegeoCode 逆地理编码 返回结构
type RegeoCode struct {
	FormattedAddress string                `json:"formatted_address"`
	AddressComponent regeoAddressComponent `json:"addressComponent"`
}

type regeoAddressComponent struct {
	Country       string                   `json:"country"`
	Province      string                   `json:"province"`
	City          string                   `json:"city"`
	CityCode      string                   `json:"citycode"`
	District      string                   `json:"district"`
	Adcode        string                   `json:"adcode"`
	TownShip      string                   `json:"township"`
	TownCode      string                   `json:"towncode"`
	Neighborhood  regeoNeighborhood        `json:"neighborhood"`
	Building      regeoBuilding            `json:"building"`
	StreetNumber  regeoStreetNumber        `json:"streetNumber"`
	BusinessAreas []regeoCodeBusinessAreas `json:"businessAreas"`
}

type regeoNeighborhood struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type regeoBuilding struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type regeoStreetNumber struct {
	Street    string `json:"street"`
	Number    string `json:"number"`
	Location  string `json:"location"`
	Direction string `json:"direction"`
	Distance  string `json:"distance"`
}

type regeoCodeBusinessAreas struct {
	Location string `json:"location"`
	Name     string `json:"name"`
	ID       string `json:"id"`
}

// Regeo 逆地理编码
// ref: http://lbs.amap.com/api/webservice/guide/api/georegeo#regeo
// 不支持 extentions 参数，默认为 base返回类型
// 不支持 batch 参数，默认为 false
func Regeo(params RegeoParams) (*RegeoCode, error) {
	m := make(map[string]interface{})
	m["key"] = key
	m["output"] = "json"
	m["location"] = params.Location.String()
	if params.PoiType != nil {
		m["poitype"] = strings.Join(params.PoiType, "|")
	}
	m["radius"] = params.Radius
	m["roadlevel"] = params.RoadLevel

	querys := util.HTTPBuildQuery(m, util.QUERY_RFC3986)
	body, err := util.HTTPGet(gateway+apiGeoCodeRe+"?"+querys, 3)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("%s\n", body)

	res := new(RegeoResponse)
	res.RegeoCode = new(RegeoCode)
	json.Unmarshal(body, &res)
	if res != nil && res.RegeoCode != nil {
		return res.RegeoCode, nil
	}

	return new(RegeoCode), nil
}
