package api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cocktailcloud/acloud-monitoring-common/v2/logger"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/types"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/utils"
	"github.com/labstack/echo/v4"
	"github.com/zemirco/uid"
)

// 직접 전달 받는 경우는 없음(deprecated)
func (a *API) ReceiveAlarms(c echo.Context) error {
	request := c.Request()
	response := c.Response()

	requestId := request.Header.Get("request-id")
	if len(requestId) == 0 {
		requestId = uid.New(10)
	}

	logger.Infof("## Request >> %s | %s | GET %s | %s\n", time.Now().Format(utils.DefaultTimeFormat),
		request.RemoteAddr, request.URL.Path, requestId)
	if a.Config.RequestTrace {
		logger.Infof("[TRACE] %s %s %s %s\n", "receivehMetrics", "0", requestId, time.Now().Format(utils.DefaultTimeFormat))
	}

	clusterId := request.Header.Get("Cluster-Id")
	if len(clusterId) == 0 {
		if request.Body != nil {
			_, _ = ioutil.ReadAll(request.Body)
			_ = request.Body.Close()
		}
		logger.Debug("cluster-id not exists")
		response.WriteHeader(http.StatusBadRequest)
		return nil
	}
	logger.Info("alert came from cluster: ", clusterId)

	// auth check
	if !a.Config.DevMode || !a.Config.SkipAuth {
		ok, err := utils.CheckAcloudAuth(request, a.Config.AuthServerUrl)
		if err != nil {
			if request.Body != nil {
				_, _ = ioutil.ReadAll(request.Body)
				_ = request.Body.Close()
			}
			logger.Errorf("fail to authenticate request from %s: %s", clusterId, err.Error())
			if ok {
				response.WriteHeader(http.StatusInternalServerError)
			} else {
				response.WriteHeader(http.StatusBadRequest)
			}
			_, _ = response.Write([]byte(err.Error()))
			return nil
		} else if !ok {
			if request.Body != nil {
				_, _ = ioutil.ReadAll(request.Body)
				_ = request.Body.Close()
			}
			response.WriteHeader(http.StatusForbidden)
			return nil
		}
	}

	// unzip
	reader, err := gzip.NewReader(request.Body)
	if err != nil {
		logger.Error("fail to read gzip request: ", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(err.Error()))
		logger.Infof("[TRACE] %s %s %s %s\n", "receiveAlarm", "1", requestId, time.Now().Format(utils.DefaultTimeFormat))
		return nil
	}
	defer reader.Close()

	// read alarm
	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(reader)
	if err != nil {
		logger.Error("fail to read gzip content: ", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(err.Error()))
		logger.Infof("[TRACE] %s %s %s %s\n", "receiveAlarm", "2", requestId, time.Now().Format(utils.DefaultTimeFormat))
		return nil
	}

	decoder := json.NewDecoder(bytes.NewReader(buffer.Bytes()))
	var envelop types.TransmitEnvelop
	if err := decoder.Decode(&envelop); err != nil {
		logger.Error("fail to read response(decode): ", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(err.Error()))
		logger.Infof("[TRACE] %s %s %s %s\n", "receiveAlarm", "3", requestId, time.Now().Format(utils.DefaultTimeFormat))
		return nil
	}

	logger.Debugf("alarm from %s\n%s", clusterId, string(envelop.Alerts))

	// save alarm
	// inserted, err := a.Service.InsertAlerts(&envelop)
	inserted, err := a.Service.ReceiveAlarmsProcess(&envelop)
	if err != nil {
		logger.Error("fail to proess alerts: ", err.Error())
	}

	// notify if needed
	if a.Config.DoAlarmNotify && len(inserted) > 0 {
		buffer := bytes.Buffer{}
		buffer.WriteString(`{"data":[`)
		for i := 0; i < len(inserted); i++ {
			if i > 0 {
				buffer.WriteString(",")
			}
			buffer.Write(inserted[i])
		}
		buffer.WriteString("]}")

		go a.NofityAlarms(buffer.Bytes())
	}

	logger.Debug("alert after processing")
	logger.Debug(string(envelop.Alerts))

	if request.Body != nil {
		_, _ = ioutil.ReadAll(request.Body)
		_ = request.Body.Close()
	}

	logger.Infof("[TRACE] %s %s %s %s\n", "receiveAlarm", "4", requestId, time.Now().Format(utils.DefaultTimeFormat))

	return nil
}

func (a *API) NofityAlarms(alarms []byte) {
	requestUrl := a.Config.AlarmNotificationApiUrl
	resp, err := http.Post(requestUrl, "application/json", bytes.NewReader(alarms))
	if resp != nil && resp.Body != nil {
		_, _ = ioutil.ReadAll(resp.Body)
	}

	if err != nil {
		logger.Error(err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("request fail: alarm notify api returns no ok staus - %d", resp.StatusCode)
	} else {
		logger.Info("alarm notify ok")
	}
}
