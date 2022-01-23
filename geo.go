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
			FormatTedAddress interface{} `json:"formatted_address"` // 地址信息
			AddressComponent struct {
				Country  interface{} `json:"country"`        // 国家
				Province interface{} `json:"province"`       // 省
				City     interface{} `json:"city,omitempty"` // 市
				CityCode interface{} `json:"citycode"`       // 城市编码
				District interface{} `json:"district"`       // 区
				AdCode   interface{} `json:"adcode"`         // 区编码
				Township interface{} `json:"township"`       // 乡镇/街道
				TownCode interface{} `json:"towncode"`       // 乡镇街道编码
			} `json:"addressComponent"` // 地址元素
		} `json:"regeocode"` // 逆地理编码
	}

	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		return nil, err
	}

	if resData.Status != "1" {
		return nil, errors.New(resData.Info)
	}

	reGeoCode := new(ReGeoCode)
	reGeoCode.Country, _ = resData.ReGeoCode.AddressComponent.Country.(string)
	reGeoCode.Province, _ = resData.ReGeoCode.AddressComponent.Province.(string)
	reGeoCode.City, _ = resData.ReGeoCode.AddressComponent.City.(string)
	reGeoCode.CityCode, _ = resData.ReGeoCode.AddressComponent.CityCode.(string)
	reGeoCode.District, _ = resData.ReGeoCode.AddressComponent.District.(string)
	reGeoCode.AdCode, _ = resData.ReGeoCode.AddressComponent.AdCode.(string)
	reGeoCode.Address, _ = resData.ReGeoCode.FormatTedAddress.(string)

	return reGeoCode, nil
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
