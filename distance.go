package amaplbs

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"errors"
)

type DistanceResponse struct {
	Status  string           `json:"status"`  // 返回结果状态值，值为0或1，0表示请求失败；1表示请求成功
	Info    string           `json:"info"`    // 返回状态说明，status为0时，info返回错误原；否则返回"OK"
	Results []DistanceResult `json:"results"` // 距离信息列表
}

type DistanceResult struct {
	OriginId string `json:"origin_id"` // 起点坐标，起点坐标序列号（从１开始）
	DestId   string `json:"dest_id"`   // 终点坐标，终点坐标序列号（从１开始）
	Distance string `json:"distance"`  // 路径距离，单位：米
	Duration string `json:"duration"`  // 预计行驶时间，单位：秒
}

// @Title Get distance list
// @Description Get the distance between two points
// @Param	origins		[]string  	"出发点坐标列表"
// @Param	destination	string 		"目的地坐标"
// @Success distances([]int) nil(error)
// @Failure nil([]int) err(error)
func (a *AmapLbs) BeeLineDistance(origins []string, destination string) ([]int, error) {
	queryParam := make(url.Values, 4)
	queryParam.Add("key", a.Key)
	queryParam.Add("origins", strings.Join(origins, "|"))
	queryParam.Add("destination", destination)
	queryParam.Add("type", "0")

	var b bytes.Buffer
	b.WriteString("http://restapi.amap.com/v3/distance?")
	b.WriteString(queryParam.Encode())
	res, err := http.Get(b.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	data := new(DistanceResponse)
	if err = json.Unmarshal(resData, data); err != nil {
		return nil, err
	}

	// request failure
	if data.Status != "1" {
		return nil, errors.New(data.Info)
	}

	distances := make([]int, len(data.Results))
	for i := range data.Results {
		index, err := strconv.Atoi(data.Results[i].OriginId)
		if err != nil {
			return nil, err
		}
		distances[index-1], err = strconv.Atoi(data.Results[i].Distance)
		if err != nil {
			return nil, err
		}
	}

	return distances, nil
}
