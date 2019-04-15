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

type HttpAccess struct {
	store *store.Store
}

type HttpResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (ha *HttpAccess) Listen(addr string) (err error) {
	err = http.ListenAndServe(addr, ha)
	return
}

func (ha *HttpAccess) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var result HttpResult //正常返回
	err := r.ParseForm()  //异常返回
	if err != nil {
		result = HttpResult{301, "Parse fail : " + err.Error()} //表单异常
	} else {
		switch strings.ToUpper(r.FormValue("do")) {
		case "POLL":
			//心跳
			result = HttpResult{200, "OK"}
		case "PUT":
			var ttl int64 = 0
			if r.FormValue("expire") != "" { //传入了expire
				seconds, err := strconv.Atoi(r.FormValue("expire"))
				if err != nil {
					result = HttpResult{302, "Convert fail : " + err.Error()} //过期时间异常
				} else {
					if seconds > 0 {
						ttl = time.Now().Add(time.Duration(seconds) * time.Second).Unix()
					}
				}
			}
			if result.Code == 0 {
				err = ha.store.Put(&store.Entry{r.FormValue("key"), r.FormValue("value"), ttl})
				if err != nil { //Put异常
					result = HttpResult{303, "Put fail : " + err.Error()}
				} else { //操作成功
					result = HttpResult{200, "Put success"}
				}
			}
		case "GET":
			value, err := ha.store.Get(r.FormValue("key"))
			if err != nil {
				result = HttpResult{303, "Get fail : " + err.Error()}
			} else {
				result = HttpResult{200, value}
			}
		case "DEL":
			err := ha.store.Del(r.FormValue("key"))
			if err != nil {
				result = HttpResult{303, "Del fail : " + err.Error()}
			} else {
				result = HttpResult{200, "Del success"}
			}
		default:
			result = HttpResult{404, "Unknow action"}
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
