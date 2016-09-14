package controllers

import (
	"encoding/json"
	"jfapp/constant"
	"jfapp/models"
	"jfapp/util/commonutil"
	"jfapp/util/dateutils"
	"jfapp/util/netutil"
	"jfapp/util/oktutil"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

/**
 *宜信金服app接口 v1.1.0版本
 */
type AppServerController struct {
	BaseController
}

//银行列表
func (this *AppServerController) BankList() {
	defer this.Recover()
	reqJson, respData := this.Decrypt()
	if respData == nil {
		deviceId, err := reqJson.Get(constant.APP_SERVICE_DEVICEID).String()
		if err != nil || deviceId == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_REQ_DEVICEID)
			this.ServeJSON()
			return
		}

		//请求参数拼接
		paramBuffer := oktutil.NewStringBuffer()
		paramBuffer.AppendStr("?").AppendStr(constant.APP_SERVICE_DEVICEID).AppendStr("=").AppendStr(deviceId)
		paramBuffer.AppendStr("&").AppendStr(constant.MERCHANTID).AppendStr("=").AppendStr(constant.JFAPP_MERCHID)
		paramBuffer.AppendStr("&").AppendStr(constant.USERID).AppendStr("=").AppendStr(deviceId)

		urlBuffer := oktutil.NewStringBufferWithStr(constant.App_turnOver_host).AppendStr(constant.App_turnOver_bankList)
		urlBuffer.AppendStr(paramBuffer.ToString())

		commonutil.LogInstance().Info(urlBuffer.ToString())

		resultStr := netutil.CallInterface(urlBuffer.ToString(), "") //调用bcw服务

		if resultStr == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_WRONG_EMPTYDATA)
			this.ServeJSON()
			return
		}

		if constant.APP_SYSTEM_ERROR_MSG == resultStr {
			this.Data["json"] = this.SetRetInfo(constant.APP_SYSTEM_ERROR)
			this.ServeJSON()
			return
		}

		respInfo := simplejson.New()
		resultJson, _ := simplejson.NewJson([]byte(resultStr))
		for _, v := range resultJson.Get(constant.APP_BANK_List).MustArray() {
			if bank, ok := v.(map[string]interface{}); ok {
				var gotoUrl, bankCode string
				gotoUrl = (bank[constant.APP_TURNOVER_GOTOURL]).(string)

				if gotoUrl == "" {
					this.Data["json"] = this.SetRetInfo(constant.APP_WRONG_EMPTY_GOTOURL)
					this.ServeJSON()
					return
				}

				bankCode = (bank[constant.APP_TURNOVER_BANKCODE]).(string)
				if bankCode == "" {
					this.Data["json"] = this.SetRetInfo(constant.APP_WRONG_EMPTY_BANKCODE)
					this.ServeJSON()
					return
				} else {
					bankCode = strings.ToLower(bankCode)
				}

				firstParam := oktutil.NewStringBuffer().AppendStr(constant.APP_TURNOVER_FIRSTACTION).AppendStr("?")
				firstParam.AppendStr(constant.APP_SERVICE_BANKNAME).AppendStr("=").AppendStr(bankCode)
				firstUrl := gotoUrl + firstParam.ToString()
				bank[constant.APP_TURNOVER_BANKCODE] = bankCode
				bank[constant.APP_TURNOVER_BANKCODE] = firstUrl

				loginUrl := constant.App_bas_host + constant.App_bas_loginUrl
				bank[constant.LOGINURL] = loginUrl
				bank[constant.JQLOGO] = constant.App_bas_host + constant.JQLOGO_PREFIX + bankCode + ".png"

			}
		}
		respInfo.Set(constant.APP_BANK_ARRAY, resultJson.Get(constant.APP_BANK_List))
		respInfo = this.Encrypt(respInfo)
		this.Data["json"] = this.getRetJson(constant.APP_RETURN_OK, respInfo)
		this.ServeJSON()

	} else {
		this.Data["json"] = respData
		this.ServeJSON()
	}
}

//卡列表
func (this *AppServerController) HistoryList() {
	defer this.Recover()
	reqJson, respData := this.Decrypt()
	if respData == nil {
		deviceId, err := reqJson.Get(constant.APP_SERVICE_DEVICEID).String()
		if err != nil || deviceId == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_REQ_DEVICEID)
			this.ServeJSON()
			return
		}

		//请求参数拼接
		paramBuffer := oktutil.NewStringBuffer()
		paramBuffer.AppendStr("?").AppendStr(constant.APP_SERVICE_DEVICEID).AppendStr("=").AppendStr(deviceId)
		paramBuffer.AppendStr("&").AppendStr(constant.MERCHANTID).AppendStr("=").AppendStr(constant.JFAPP_MERCHID)
		paramBuffer.AppendStr("&").AppendStr(constant.USERID).AppendStr("=").AppendStr(deviceId)

		urlBuffer := oktutil.NewStringBufferWithStr(constant.App_turnOver_host).AppendStr(constant.App_turnOver_historyList)
		urlBuffer.AppendStr(paramBuffer.ToString())

		resultStr := netutil.CallInterface(urlBuffer.ToString(), "") //调用bcw服务
		commonutil.LogInstance().Info("宜信金服|流水爬取|已查询卡列表接口返回信息为：" + resultStr)
		if resultStr == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_WRONG_EMPTYDATA)
			this.ServeJSON()
			return
		}

		if constant.APP_SYSTEM_ERROR_MSG == resultStr {
			this.Data["json"] = this.SetRetInfo(constant.APP_SYSTEM_ERROR)
			this.ServeJSON()
			return
		}

		respInfo, _ := simplejson.NewJson([]byte(resultStr))

		loginUrl := oktutil.NewStringBufferWithStr(constant.App_bas_host).AppendStr(constant.App_bas_loginUrl).ToString()

		for _, v := range respInfo.Get(constant.CARD_ARRAY).MustArray() {
			if card, ok := v.(map[string]interface{}); ok {
				errorMsg := (card[constant.ERRORMSG]).(string)
				if errorMsg != "" && "单日爬取次数超过限制，请明天再试" == errorMsg {
					errorMsg = "本日查询次数超过限制，请明天再试"
					card[constant.ERRORMSG] = errorMsg
				}
				card[constant.LOGINURL] = loginUrl
			}
		}
		respInfo = this.Encrypt(respInfo)
		this.Data["json"] = this.getRetJson(constant.APP_RETURN_OK, respInfo)
		this.ServeJSON()

	} else {
		this.Data["json"] = respData
		this.ServeJSON()
	}
}

//流水详情接口
func (this *AppServerController) BasInfo() {
	defer this.Recover()
	reqJson, respData := this.Decrypt()
	if respData == nil {
		deviceId, err := reqJson.Get(constant.APP_SERVICE_DEVICEID).String()
		if err != nil || deviceId == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_REQ_DEVICEID)
			this.ServeJSON()
			return
		}

		custName, err := reqJson.Get(constant.APP_AUTHENTICATE_CUSTNAME).String()
		if err != nil || custName == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_REQ_CUSTNAME)
			this.ServeJSON()
			return
		}

		cardId, err := reqJson.Get(constant.APP_AUTHENTICATE_CARDID).String()
		if err != nil || cardId == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_REQ_MSGCODE)
			this.ServeJSON()
			return
		}

		bankCode, err := reqJson.Get(constant.APP_AUTHENTICATE_BANKCODE).String()
		if err != nil || bankCode == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_REQ_BANKCODE)
			this.ServeJSON()
			return
		}

		refresh, err := reqJson.Get(constant.APP_AUTHENTICATE_REFRESH).String()
		if err != nil || refresh == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_REQ_REFRESH)
			this.ServeJSON()
			return
		}

		endMonth := dateutils.GetCurrentDateStr(dateutils.YYYY_MM)
		endDate := dateutils.GetCurrentDateStr(dateutils.YYYY_MM_DD)
		beginMonth := dateutils.GetDateStr(time.Now().AddDate(0, -7, 0), dateutils.YYYY_MM)
		beginDate := beginMonth + "-01"

		//请求参数拼接
		paramBuffer := oktutil.NewStringBuffer()
		paramBuffer.AppendStr("?").AppendStr(constant.APP_SERVICE_DEVICEID).AppendStr("=").AppendStr(deviceId)
		paramBuffer.AppendStr("&").AppendStr(constant.MERCHANTID).AppendStr("=").AppendStr(constant.JFAPP_MERCHID)
		paramBuffer.AppendStr("&").AppendStr(constant.USERID).AppendStr("=").AppendStr(deviceId)
		paramBuffer.AppendStr("&").AppendStr(constant.APP_AUTHENTICATE_CARDID).AppendStr("=").AppendStr(cardId)
		paramBuffer.AppendStr("&").AppendStr(constant.APP_AUTHENTICATE_CUSTNAME).AppendStr("=").AppendStr(custName)
		paramBuffer.AppendStr("&").AppendStr(constant.APP_AUTHENTICATE_BANKCODE).AppendStr("=").AppendStr(bankCode)
		paramBuffer.AppendStr("&").AppendStr(constant.BEGINDATE).AppendStr("=").AppendStr(beginDate)
		paramBuffer.AppendStr("&").AppendStr(constant.ENDDATE).AppendStr("=").AppendStr(endDate)

		basInUrlBuffer := oktutil.NewStringBufferWithStr(constant.App_turnOver_host).AppendStr(constant.App_turnOver_basInfo)
		basInUrlBuffer.AppendStr(paramBuffer.ToString())

		resultStr := netutil.CallInterface(basInUrlBuffer.ToString(), "") //调用bcw服务
		commonutil.LogInstance().Info("宜信金服|流水爬取|流水详情接口返回信息为：" + resultStr)

		if resultStr == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_WRONG_EMPTYDATA)
			this.ServeJSON()
			return
		}

		if constant.APP_SYSTEM_ERROR_MSG == resultStr {
			this.Data["json"] = this.SetRetInfo(constant.APP_SYSTEM_ERROR)
			this.ServeJSON()
			return
		}

		respInfo, _ := simplejson.NewJson([]byte(resultStr))

		subTradeJson, ok := respInfo.CheckGet(constant.SUBTRADELIST)

		if !ok {
			this.Data["json"] = this.SetRetInfo(constant.APP_SYSTEM_ERROR)
			this.ServeJSON()
			return
		}

		tradeDataList := make([]*models.TradeData, 0)
		for i := 0; i < len(subTradeJson.MustArray()); i++ {
			tradeDataBytes, _ := subTradeJson.GetIndex(i).MarshalJSON()
			tradeData := &models.TradeData{}
			json.Unmarshal(tradeDataBytes, tradeData)
			tradeDataList = append(tradeDataList, tradeData)
		}
		var createTime string
		if tradeDataList[0] != nil {
			createTime = tradeDataList[0].QueryTime
		}

		resultJson := this.HandleTradeData(tradeDataList, endMonth, constant.EFFECTMONTH)

		effectiveType, _ := resultJson.Get("status").String()
		commonutil.LogInstance().Info("卡号|" + cardId + "的流水有效性标记为：" + effectiveType)
		var basType string

		if createTime != "" {
			queryDate := dateutils.Str2Date(createTime, dateutils.YYYYMMDDHHMMSS).Unix()
			currentTimeSec := time.Now().Unix()
			cha := int((currentTimeSec - queryDate) / 60 / 60 / 24)
			if cha <= constant.BASTYPE && effectiveType == constant.STRING_YES {
				basType = "优质流水"
			} else {
				basType = "普通流水"
			}
		}
		respInfo.Del(constant.SUBTRADELIST)
		respInfo.Set("basType", basType)
		respInfo.Set("lastTradeTime", endMonth)
		respInfo.Set("tradeArray", resultJson.Get("monthTradeList"))

		smallLoginUrl := oktutil.NewStringBufferWithStr(constant.App_bas_host).AppendStr(constant.App_bas_loginUrl).ToString()
		respInfo.Set(constant.SMALLLOGINURL, smallLoginUrl)
		respInfo.Set(constant.APP_AUTHENTICATE_REFRESH, refresh)

		limitParam := oktutil.NewStringBuffer()
		limitParam.AppendStr("?").AppendStr(constant.APP_SERVICE_DEVICEID).AppendStr("=").AppendStr(deviceId)
		limitParam.AppendStr("&").AppendStr(constant.MERCHANTID).AppendStr("=").AppendStr(constant.JFAPP_MERCHID)
		limitParam.AppendStr("&").AppendStr(constant.USERID).AppendStr("=").AppendStr(deviceId)

		limitUrl := oktutil.NewStringBufferWithStr(constant.App_turnOver_host).AppendStr(constant.App_turnOver_queryLimitInfo)
		limitUrl.AppendStr(limitParam.ToString())
		limitInfo := netutil.CallInterface(limitUrl.ToString(), "")

		if limitInfo == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_WRONG_EMPTY_LIMIT)
			this.ServeJSON()
			return
		}
		limitJson, _ := simplejson.NewJson([]byte(limitInfo))
		respInfo.Set(constant.LimitNumber, limitJson.Get(constant.LimitNumber))
		respInfo.Set(constant.ReportSaveDays, limitJson.Get(constant.ReportSaveDays))
		respInfo.Set(constant.TimeSpan, limitJson.Get(constant.TimeSpan))

		respInfo = this.Encrypt(respInfo)
		this.Data["json"] = this.getRetJson(constant.APP_RETURN_OK, respInfo)
		this.ServeJSON()

	} else {
		this.Data["json"] = respData
		this.ServeJSON()
	}
}

//月交易详情接口
func (this *AppServerController) MonthBasInfo() {
	defer this.Recover()
	reqJson, respData := this.Decrypt()
	if respData == nil {
		deviceId, err := reqJson.Get(constant.APP_SERVICE_DEVICEID).String()
		if err != nil || deviceId == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_REQ_DEVICEID)
			this.ServeJSON()
			return
		}
		cardId, err := reqJson.Get(constant.APP_AUTHENTICATE_CARDID).String()
		if err != nil || cardId == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_REQ_MSGCODE)
			this.ServeJSON()
			return
		}
		date, err := reqJson.Get(constant.DATE).String()
		if err != nil || date == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_REQ_DATE)
			this.ServeJSON()
			return
		}

		maxDay := dateutils.GetMonthActualMaximum(dateutils.Str2Date(date, dateutils.YYYY_MM))
		beginDate := date + "-01"
		endDate := date + "-" + strconv.Itoa(maxDay)

		//请求参数拼接
		paramBuffer := oktutil.NewStringBuffer()
		paramBuffer.AppendStr("?").AppendStr(constant.APP_SERVICE_DEVICEID).AppendStr("=").AppendStr(deviceId)
		paramBuffer.AppendStr("&").AppendStr(constant.MERCHANTID).AppendStr("=").AppendStr(constant.JFAPP_MERCHID)
		paramBuffer.AppendStr("&").AppendStr(constant.USERID).AppendStr("=").AppendStr(deviceId)
		paramBuffer.AppendStr("&").AppendStr(constant.APP_AUTHENTICATE_CARDID).AppendStr("=").AppendStr(cardId)
		paramBuffer.AppendStr("&").AppendStr(constant.BEGINDATE).AppendStr("=").AppendStr(beginDate)
		paramBuffer.AppendStr("&").AppendStr(constant.ENDDATE).AppendStr("=").AppendStr(endDate)

		url := oktutil.NewStringBufferWithStr(constant.App_turnOver_host).AppendStr(constant.App_turnOver_basInfo)
		url.AppendStr(paramBuffer.ToString())

		resultStr := netutil.CallInterface(url.ToString(), "")

		commonutil.LogInstance().Info("宜信金服|流水爬取|月交易详情接口返回信息为：" + resultStr)
		if resultStr == "" || constant.EMPTYJSON == resultStr {
			this.Data["json"] = this.SetRetInfo(constant.APP_WRONG_EMPTYDATA)
			this.ServeJSON()
			return
		}

		if constant.APP_SYSTEM_ERROR_MSG == resultStr {
			this.Data["json"] = this.SetRetInfo(constant.APP_SYSTEM_ERROR)
			this.ServeJSON()
			return
		}

		respInfo := simplejson.New()
		resultJson, _ := simplejson.NewJson([]byte(resultStr))

		if errJson, ok := resultJson.CheckGet(constant.ERRORMSG); ok {
			errMsg, _ := errJson.String()
			this.Data["json"] = this.SetRetJson2("9999", errMsg, errJson)
			this.ServeJSON()
			return
		}
		subTradeJson, _ := resultJson.CheckGet(constant.SUBTRADELIST)
		var returnInfo string
		if subTradeJson == nil {
			returnInfo = "[]"
		} else {
			returnInfoBytes, _ := subTradeJson.MarshalJSON()
			returnInfo = strings.Replace(string(returnInfoBytes), "tradeBlance", "tradePrice", -1)
		}

		returnJson, _ := simplejson.NewJson([]byte(returnInfo))
		respInfo.Set("tradeArray", returnJson)
		respInfo = this.Encrypt(respInfo)
		this.Data["json"] = this.getRetJson(constant.APP_RETURN_OK, respInfo)
		this.ServeJSON()

	} else {
		this.Data["json"] = respData
		this.ServeJSON()
	}

}

func (this *AppServerController) SalaryBasInfo() {
	defer this.Recover()
	reqJson, respData := this.Decrypt()
	if respData == nil {
		deviceId, err := reqJson.Get(constant.APP_SERVICE_DEVICEID).String()
		if err != nil || deviceId == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_REQ_DEVICEID)
			this.ServeJSON()
			return
		}

		custName, err := reqJson.Get(constant.APP_AUTHENTICATE_CUSTNAME).String()
		if err != nil || custName == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_REQ_CUSTNAME)
			this.ServeJSON()
			return
		}

		cardId, err := reqJson.Get(constant.APP_AUTHENTICATE_CARDID).String()
		if err != nil || cardId == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_REQ_MSGCODE)
			this.ServeJSON()
			return
		}

		bankCode, err := reqJson.Get(constant.APP_AUTHENTICATE_BANKCODE).String()
		if err != nil || bankCode == "" {
			this.Data["json"] = this.SetRetInfo(constant.APP_REQ_BANKCODE)
			this.ServeJSON()
			return
		}

		//		endDate := dateutils.GetCurrentDateStr(dateutils.YYYY_MM_DD)
	} else {
		this.Data["json"] = respData
		this.ServeJSON()
	}
}
