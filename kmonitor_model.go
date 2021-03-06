package gokids

import (
	"sync"
	"time"

)
type KMonitor struct {
	data map[string]*KMonitorKeyData
	sync.Mutex
	C      chan KMonitorReqStatus
	ticker *time.Ticker
}


// 每个Url的计数器
type KMonitorKeyData struct {
	Url         string
	TotalCount  int64
	SucessCount int64
	Duration    int64
	ProductKey	string		//应用中心产品唯一标志
}

// 每个请求的处理结果
type KMonitorReqStatus struct {
	Key      string // key
	Status   int    // 处理结果，成功=1，失败=0
	Duration int64  // 处理时间，ms
	ProductKey	string		//应用中心产品唯一标志
}
type KMonitorReportReqData struct {
	Code         string `json:"m"`  //服务码
	SKey         string `json:"s"`  //子key，监控点
	TotalCount   int64  `json:"tc"` //总数量
	SuccessCount int64  `json:"ts"` //成功数量
	Duration     int64  `json:"tt"` //总耗时
}

// 采样上报数据
type KMonitorReportReq struct {
	T string          `json:"t"` //采样时间
	P string          `json:"p"` //产品编码
	R []KMonitorReportReqData `json:"r"` //产品编码
}
