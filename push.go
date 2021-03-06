package gokids

import (
    "strconv"
    "time"
    "net/url"
    "encoding/json"
    "errors"
)

func PushSendOcmTask(mqPushUrl string, mqKey string, pushServiceUrl string, taskCodePrefix string, taskCode string, startTime time.Time, sourceSystem string,
    aimSystem string, content string, uids []uint64) (error) {
    /*
        1、发送任务
        2、配置任务和用户到mq(推送)
    */
    taskId, b := PushApiAddTask(pushServiceUrl, taskCodePrefix, taskCode, startTime, sourceSystem, aimSystem, content)
    if !b {
        return errors.New("gokids push add task api error" + pushServiceUrl + taskCodePrefix + taskCode + startTime.String() + sourceSystem + aimSystem + content)
    }
    err := PushAddTaskToMq(mqPushUrl, mqKey, taskId, uids)
    if err != nil {
        return err
    }

    return nil
}

//添加推送任务
func PushApiAddTask(serviceUrl string, taskCodePrefix string, taskCode string, startTime time.Time, sourceSystem string,
    aimSystem string, content string) (int, bool) {
    queryData := ModelPushGetApiTaskInfo(taskCodePrefix, taskCode, startTime, sourceSystem, aimSystem, content)

    taskCodeInt, _ := strconv.Atoi(queryData.TaskCode)

    params := make(url.Values)
    params.Set("taskCode", queryData.TaskCode)
    params.Set("name", queryData.Name)
    params.Set("startTime", queryData.StartTime)
    params.Set("failureTime", queryData.FailureTime)
    params.Set("content", queryData.Content)
    params.Set("memo", queryData.Memo)
    params.Set("sourceSystem", queryData.SourceSystem)
    params.Set("aimSystem", queryData.AimSystem)
    params.Set("taskType", queryData.TaskType)
    header := []string{"Accept:", "Content-Type:charset=UTF-8"}
    paramsStr := params.Encode()
    serviceUrl += "?" + paramsStr
    ret, err := UtilCurlGet(serviceUrl, header)
    if err != nil {
        return taskCodeInt, false
    }
    data := ModelPushApiResponseInfo{}
    err = json.Unmarshal(ret, &data)

    if err != nil {
        log.Error("push neibu api add task error query: %s error: %s", paramsStr, err.Error())
    }

    return taskCodeInt, data.Success
}

//配置推送任务到mq
func PushAddTaskToMq(mqPushUrl string, mqKey string, taskId int, uids []uint64) (error) {
    mqData := ModelPushMqOcmTaskInfo{}
    mqData.TaskCode = taskId
    mqData.DeviceType = -1
    for _, v := range uids {
        info := ModelPushMqOcmTaskRowsInfo{}
        info.CustomerId = v
        info.Badge = 1
        info.Sound = "default"
        mqData.Rows = append(mqData.Rows, info)
    }
    mqData.Count = len(mqData.Rows)

    return UtilMQSend(mqPushUrl, mqKey, mqData)
}