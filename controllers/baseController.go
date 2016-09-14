package controllers

import (
	"fmt"
	"jfapp/constant"
	"jfapp/models"
	"jfapp/util/aescrypt"
	"jfapp/util/commonutil"
	"jfapp/util/dateutils"
	"jfapp/util/oktutil"
	"log"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

//json映射 struct
type RespData struct {
	RetCode string      `json:"retCode"`
	RetMsg  string      `json:retMsg`
	Data    interface{} `json:"data"`
}

type BaseController struct {
	beego.Controller
}

//解密数据
func (this *BaseController) Decrypt() (*simplejson.Json, *RespData) {
	data := this.GetString(constant.APP_PARAM_DATA)
	if data == "" {
		return nil, this.SetRetInfo(constant.APP_REQ_DATA)
	}
	json, _ := simplejson.NewJson([]byte(data))
	reqInfo, _ := json.Get(constant.APP_SERVICE_REQ).String()
	if reqInfo == "" {
		return nil, this.SetRetInfo(constant.APP_REQ_REQINFO)
	}

	okt_req, _ := json.Get(constant.OKT).String()
	commonutil.LogInstance().Info("客户端的密钥标识解码前为：" + okt_req)
	if okt_req == "" {
		return nil, this.SetRetInfo(constant.APP_REQ_OKT)
	}
	okt_decode, _ := aescrypt.Base64Decode(okt_req)
	commonutil.LogInstance().Info("客户端的密钥标识解码前为：" + string(okt_decode))

	oktInt, err := strconv.Atoi(string(okt_decode))
	if err != nil {
		panic("okt_decode 不能转为 int ")
	}
	keyIndex := oktutil.GetIndex(oktInt)
	keyLength := oktutil.GetLength(oktInt)

	appKey_req := oktutil.Substr2(aescrypt.KEY_SPACE, keyIndex, keyIndex+keyLength)

	commonutil.LogInstance().Info("客户端的AES key为=>" + appKey_req)

	commonutil.LogInstance().Info("解密前的字符串为：" + reqInfo)
	result, err := aescrypt.DecryptByOkt(appKey_req, reqInfo)
	commonutil.LogInstance().Info("解密后的字符串为：" + result)

	if err != nil || result == "" {
		return nil, this.SetRetInfo(constant.APP_REQ_REQINFO)
	}

	reqJson, _ := simplejson.NewJson([]byte(result))

	signInfo_req, _ := json.Get(constant.APP_SERVICE_SIGNINFO).String()

	if signInfo_req == "" {
		return nil, this.SetRetInfo(constant.APP_REQ_SIGNINFO)
	}

	commonutil.LogInstance().Info("客户端的签名为：" + signInfo_req)
	signInfo_service := aescrypt.GenerateSignInfo(reqJson, appKey_req)
	commonutil.LogInstance().Info("服务端生成的签名为：" + signInfo_service)

	if signInfo_service == "" {
		return nil, this.SetRetInfo(constant.APP_WRONG_NOSIGN)
	} else if signInfo_service != signInfo_req {
		return nil, this.SetRetInfo(constant.APP_WRONG_SIGNINFO)
	}

	return reqJson, nil
}

//对json数据进行加密，签名
func (this *BaseController) Encrypt(json *simplejson.Json) *simplejson.Json {

	resultJson := simplejson.New()

	keyIndex := oktutil.MakeRandomIndex(len(aescrypt.KEY_SPACE))
	keyLength := oktutil.MakeRandomLength()
	appKey := oktutil.Substr2(aescrypt.KEY_SPACE, keyIndex, keyIndex+keyLength)
	okt := oktutil.MakeOKTCode(keyIndex, keyLength)

	okt_encode := aescrypt.Base64Encode([]byte(strconv.Itoa(okt)))

	data, _ := json.MarshalJSON()
	fmt.Println(string(data))
	dataEecode, err := aescrypt.Encrypt(appKey, data)
	signInfo := aescrypt.GenerateSignInfo(json, appKey)
	if err != nil {
		return nil
	}

	resultJson.Set(constant.OKT, okt_encode)
	resultJson.Set(constant.APP_SERVICE_RESP, dataEecode)
	resultJson.Set(constant.APP_SERVICE_SIGNINFO, signInfo)
	return resultJson
}

//设置响应数据
func (this *BaseController) SetRetInfo(respInfo constant.BusinessRetInfo) *RespData {
	respJson := simplejson.New()
	return &RespData{respInfo.Code, respInfo.Msg, respJson}
}

//返回数据
func (this *BaseController) SetRetJson2(code, msg string, encryptJson *simplejson.Json) *RespData {
	//	resultJson := simplejson.New()
	//	resultJson.Set(constant.APP_RETCODE, code)
	//	resultJson.Set(constant.APP_RETMSG, msg)
	//	resultJson.Set(constant.APP_PARAM_DATA, encryptJson)
	return &RespData{code, msg, encryptJson}
}

//简要请求必要参数
func (this *BaseController) checkParams() (*RespData, bool) {
	reqData := this.GetString(constant.APP_PARAM_DATA)
	var res *RespData
	var isPass bool = true
	if reqData == "" {
		res = this.SetRetInfo(constant.APP_REQ_DATA)
		isPass = false
	} else {
		reqJson, _ := simplejson.NewJson([]byte(reqData))
		dataJson := reqJson.Get(constant.APP_SERVICE_REQ)
		okt_json := reqJson.Get(constant.OKT)
		if dataJson.Interface() == nil {
			res = this.SetRetInfo(constant.APP_REQ_REQINFO)
			isPass = false
		}
		if okt_json.Interface() == nil {
			res = this.SetRetInfo(constant.APP_REQ_OKT)
			isPass = false
		}
	}
	return res, isPass
}

//返回数据
func (this *BaseController) getRetJson(respInfo constant.BusinessRetInfo, encryptJson *simplejson.Json) *RespData {
	//	resultJson := simplejson.New()
	//	resultJson.Set(constant.APP_RETCODE, respInfo.Code)
	//	resultJson.Set(constant.APP_RETMSG, respInfo.Msg)
	//	resultJson.Set(constant.APP_PARAM_DATA, encryptJson)
	return &RespData{respInfo.Code, respInfo.Msg, encryptJson}
}

func (this *BaseController) Recover() { //异常恢复
	if r := recover(); r != nil {
		log.Printf("Runtime error caught: %v", r)
		this.Data["json"] = this.SetRetInfo(constant.APP_SYSTEM_ERROR)
		this.ServeJSON()
	}
}

func (this *BaseController) HandleTradeData(subTradeList []*models.TradeData, endMonth string, effectMonth int) *simplejson.Json {
	resultJson := simplejson.New()

	monthTradeList := make([]models.MonthTrade, 0)
	// status 0 默认流水类型：无效
	var status, bothEnds, effectiveType string = "0", "0", "1"
	for i := 0; i < (effectMonth + 1); i++ {
		monthTrade := new(models.MonthTrade)
		var income, experss float64
		var monthTradeTime, monthLenth string

		beforeMonth := dateutils.AddMonth(dateutils.Str2Date(endMonth, dateutils.YYYY_MM), -i)
		monthStr_i := dateutils.GetDateStr(beforeMonth, dateutils.YYYY_MM)

		data_i := make([]models.TradeData, 0)

		for _, tradeData := range subTradeList {
			var tradeTime_month string
			tradeTime := tradeData.TradeTime
			if tradeTime != "" {
				tradeTime_month = dateutils.GetDateStr(dateutils.Str2Date(tradeTime, dateutils.YYYY_MM_DD_HH_MM_SS), dateutils.YYYY_MM)
			}
			if tradeTime_month == monthStr_i {
				data_i = append(data_i, *tradeData)
				tradeType := tradeData.TradeType
				if tradeType != "" {
					tradeBlance, _ := strconv.ParseFloat(tradeData.TradeBlance, 64)
					if tradeType == constant.STRING_YES {
						income = income + tradeBlance
					} else if tradeType == constant.STRING_NO {
						experss = experss + tradeBlance
					}
				}
			}

		}

		if i == 0 || i == effectMonth {
			effect := this.validateEffect(data_i)
			if effect == constant.STRING_YES {
				bothEnds = constant.STRING_YES
			}
		} else {
			effect := this.validateEffect(data_i)
			if effect == constant.STRING_NO {
				effectiveType = constant.STRING_NO
			}
		}

		dateArr := strings.Split(monthStr_i, "-")
		if dateArr != nil && len(dateArr) == 2 {
			monthTradeTime = dateArr[1] + "月" + dateArr[0]
			maxDay := dateutils.GetMonthActualMaximum(beforeMonth)
			monthLenth = dateArr[1] + ".01-" + dateArr[1] + "." + strconv.Itoa(maxDay)
		}

		monthTrade.Income = strconv.FormatFloat(income, 'f', 2, 64)
		monthTrade.Expenses = strconv.FormatFloat(experss, 'f', 2, 64)
		monthTrade.TradeTime = monthTradeTime
		monthTrade.Date = monthStr_i
		monthTrade.MonthLenth = monthLenth
		monthTradeList = append(monthTradeList, *monthTrade)
	}

	if bothEnds == constant.STRING_YES && effectiveType == constant.STRING_YES {
		status = constant.STRING_YES
	}
	resultJson.Set("status", status)
	resultJson.Set("monthTradeList", monthTradeList)
	return resultJson
}

func (this *BaseController) validateEffect(data_i []models.TradeData) string {
	effect := constant.STRING_NO
	var sum float64
	for _, tradeData := range data_i {
		this.validateSalary(&tradeData)
		salaryType := tradeData.SalaryType
		if salaryType == constant.STRING_YES {
			tradeBlance, _ := strconv.ParseFloat(tradeData.TradeBlance, 64)
			sum = sum + tradeBlance
		}
	}
	if sum >= constant.SALARYNUM {
		effect = constant.STRING_YES
	}
	return effect
}

func (this *BaseController) validateSalary(tradeData *models.TradeData) {
	tradeDesc := tradeData.TradeDesc
	if tradeDesc != "" {
		for _, salaryFlag := range constant.SALARYFLAG {
			if salaryFlag == tradeDesc {
				tradeData.SalaryType = constant.STRING_YES
				break
			}
		}
	}
	salaryType := tradeData.SalaryType

	if salaryType != "" {
		tradeData.SalaryType = constant.STRING_NO
	}

}
