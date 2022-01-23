package amaplbs

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

type ReGeoCode struct {
	Country  string // 国家
	Province string // 省
	City     string // 市
	CityCode string // 城市编码
	District string // 区
	AdCode   string // 区域码
	Address  string // 地址
}

func (a *AmapLbs) ReGeoCode(longitude string, latitude string) (*ReGeoCode, error) {
	queryParam := make(url.Values, 4)
	queryParam.Add("key", a.Key)
	queryParam.Add("location", a.GenLocationTrim(longitude, 6)+","+a.GenLocationTrim(latitude, 6))
	queryParam.Add("extensions", "base")
	queryParam.Add("batch", "false")
	queryParam.Add("output", "JSON")

	var b bytes.Buffer
	b.WriteString("https://restapi.amap.com/v3/geocode/regeo?")
	b.WriteString(queryParam.Encode())
	res, err := http.Get(b.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resData struct {
		Status    string `json:"status"` // 返回结果状态值，值为0或1，0表示请求失败；1表示请求成功
		Info      string `json:"info"`   // 返回状态说明，status为0时，info返回错误原；否则返回"OK"
		ReGeoCode struct {
			FormatTedAddress string `json:"formatted_address"` // 地址信息
			AddressComponent struct {
				Country  string      `json:"country"`        // 国家
				Province string      `json:"province"`       // 省
				City     interface{} `json:"city,omitempty"` // 市
				CityCode string      `json:"citycode"`       // 城市编码
				District string      `json:"district"`       // 区
				AdCode   string      `json:"adcode"`         // 区编码
				Township string      `json:"township"`       // 乡镇/街道
				TownCode string      `json:"towncode"`       // 乡镇街道编码
			} `json:"addressComponent"` // 地址元素
		} `json:"regeocode"` // 逆地理编码
	}

	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		return nil, err
	}

	if resData.Status != "1" {
		return nil, errors.New(resData.Info)
	}

	city, _ := resData.ReGeoCode.AddressComponent.City.(string)

	return &ReGeoCode{
		Country:  resData.ReGeoCode.AddressComponent.Country,
		Province: resData.ReGeoCode.AddressComponent.Province,
		City:     city,
		CityCode: resData.ReGeoCode.AddressComponent.CityCode,
		District: resData.ReGeoCode.AddressComponent.District,
		AdCode:   resData.ReGeoCode.AddressComponent.AdCode,
		Address:  resData.ReGeoCode.FormatTedAddress,
	}, nil
}

func (a *AmapLbs) GenLocationTrim(location string, num int) string {
	var (
		locationLen = len(location)
		pointIndex  = strings.Index(location, ".")
	)
	if locationLen-pointIndex > 7 {
		location = location[:pointIndex+7]
	}
	return location
}
