package core

import (
	"Calla/store"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//HTTPAccess 解析自request
type HTTPAccess struct {
	store *store.Store
}

//HTTPResult 作为http的json返回
type HTTPResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

//Listen 用于监听地址
func (ha *HTTPAccess) Listen(addr string) (err error) {
	err = http.ListenAndServe(addr, ha)
	return
}

func (ha *HTTPAccess) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var result HTTPResult //正常返回
	err := r.ParseForm()  //异常返回
	if err != nil {
		result = HTTPResult{301, "Parse fail : " + err.Error()} //表单异常
	} else {
		switch strings.ToUpper(r.FormValue("do")) {
		case "PUT":
			var ttl int64
			if r.FormValue("expire") != "" { //传入了expire
				seconds, err := strconv.Atoi(r.FormValue("expire"))
				if err != nil {
					result = HTTPResult{302, "Convert fail : " + err.Error()} //过期时间异常
				} else {
					if seconds > 0 {
						ttl = time.Now().Add(time.Duration(seconds) * time.Second).Unix()
					}
				}
			}
			if result.Code == 0 {
				err = ha.store.Put(&store.Entry{Key: r.FormValue("key"), Value: r.FormValue("value"), Expire: ttl})
				if err != nil { //Put异常
					result = HTTPResult{303, "Put fail : " + err.Error()}
				} else { //操作成功
					result = HTTPResult{200, "Put success"}
				}
			}
		case "GET":
			value, err := ha.store.Get(r.FormValue("key"))
			if err != nil {
				result = HTTPResult{303, "Get fail : " + err.Error()}
			} else {
				result = HTTPResult{200, value}
			}
		case "DEL":
			err := ha.store.Del(r.FormValue("key"))
			if err != nil {
				result = HTTPResult{303, "Del fail : " + err.Error()}
			} else {
				result = HTTPResult{200, "Del success"}
			}
		default:
			result = HTTPResult{404, "Unknow action"}
		}
	}
	if str, err := json.Marshal(result); err != nil {
		w.WriteHeader(500)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			fmt.Println(err)
		}
	} else {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		if _, err = w.Write(str); err != nil {
			fmt.Println(err)
		}
	}
}
