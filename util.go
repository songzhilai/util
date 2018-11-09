package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// ResponseJSON encodes value to json and write as response.
func ResponseJSON(w http.ResponseWriter, errorCode int, msg string, v interface{}) error {
	r := map[string]interface{}{
		"msg":        msg,
		"error_code": errorCode,
		"data":       v,
	}

	// Encode JSON.
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}

	// Write response.
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return nil
}

//DoorHttpRequest 发送http请求
func DoorHttpRequest(method string, url string, data interface{}) (string, int) {
	var datas *bytes.Buffer
	rejson, err := json.Marshal(data)
	if err != nil {
		return "", -4
	}
	datas = bytes.NewBuffer(rejson)
	req, err := http.NewRequest(method, url, datas)
	// req.Header.Set("content-type", "application/json")
	// ba64authstr := base64.StdEncoding.EncodeToString([]byte(libs.ConfigValue("FrpsUser") + ":" + libs.ConfigValue("FrpsPasswd")))
	// req.Header.Set("Authorization", "Basic "+ba64authstr)
	if err != nil {
		return "", -1
	}
	Client := &http.Client{Timeout: 5 * time.Second}
	resp, err := Client.Do(req)
	if resp.Status != "200 OK" {
		return "", -4
	}
	if err != nil {
		return "", -2
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", -3
	}

	return string(body), 1
}

//RemoveDuplicatesAndEmpty 数组去重
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	alen := len(a)
	for i := 0; i < alen; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}
