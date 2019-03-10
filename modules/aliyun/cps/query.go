package cps

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

const (
	// QueryTypeNew 新增
	QueryTypeNew = "NEW"
	// QueryTypeTotal 累计
	QueryTypeTotal = "TOTAL"
)

// DeviceStatResponse 设备统计查询响应
type DeviceStatResponse struct {
	Response
	AppDeviceStats AppDeviceStat
}

// AppDeviceStat 响应
type AppDeviceStat struct {
	AppDeviceStat []*AppDeviceStatsItem
}

// AppDeviceStatsItem 设备统计
type AppDeviceStatsItem struct {
	Time       string
	Count      int
	DeviceType string
}

// QueryDeviceStat 设备新增与留存数据查询
func (t *CPS) QueryDeviceStat(startTime, endTime time.Time, deviceType, queryType string) ([]*AppDeviceStatsItem, error) {
	params := t.getCommonParams()
	params["Action"] = "QueryDeviceStat"
	params["DeviceType"] = deviceType
	params["QueryType"] = queryType
	params["StartTime"] = startTime.Format("2006-01-02T") + "16:00:00Z"
	params["EndTime"] = endTime.Format("2006-01-02T") + "16:00:00Z"

	// 签名
	sign := t.sign(http.MethodGet, params)
	params["Signature"] = sign

	// 获取响应
	query := url.Values{}
	for k := range params {
		query.Add(k, params[k])
	}
	res, err := t.request(CPSGateway, query)
	if err != nil {
		return nil, err
	}

	// 解析响应
	ret := new(DeviceStatResponse)
	err = json.Unmarshal(res, ret)
	if err != nil {
		return nil, err
	}

	if ret.Code != "" {
		return nil, errors.New(ret.Message)
	}

	return ret.AppDeviceStats.AppDeviceStat, nil
}
