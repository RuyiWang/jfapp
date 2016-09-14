package constant

//基本常量
const (
	APP_SYSTEM_ERROR_MSG = "error"
	APP_PARAM_DATA       = "data"
	APP_RETCODE_VALUE    = "0000"
	APP_RETMSG_VALUE     = "正常响应"
	APP_RETMSG           = "retMsg"
	APP_RETCODE          = "retCode"
	OKT                  = "OKT"
	APP_SERVICE_REQ      = "app_service_req"
	APP_SERVICE_RESP     = "app_service_resp"
	APP_SERVICE_SIGNINFO = "app_service_signinfo"
	APP_SERVICE_DEVICEID = "deviceId"
	MERCHANTID           = "merchantId"
	USERID               = "userId"

	APP_BANK_List         = "bankList"
	APP_TURNOVER_GOTOURL  = "gotoUrl"
	APP_TURNOVER_BANKCODE = "bankCode"
	APP_SERVICE_BANKNAME  = "bankname"

	APP_AUTHENTICATE_CUSTNAME = "custName"
	APP_AUTHENTICATE_CARDID   = "cardId"
	APP_AUTHENTICATE_BANKCODE = "bankCode"
	APP_AUTHENTICATE_REFRESH  = "refresh"

	BEGINDATE = "beginDate"
	ENDDATE   = "endDate"
	DATE      = "date"

	LOGINURL       = "loginUrl"
	SMALLLOGINURL  = "smallLoginUrl"
	JQLOGO         = "jqLogo"
	APP_BANK_ARRAY = "bankArray"
	CARD_ARRAY     = "cardArray"
	ERRORMSG       = "errorMsg"

	SUBTRADELIST = "subTradeList"

	/** 流水有效规则：连续6个月 */
	EFFECTMONTH = 6

	STRING_YES = "1"
	STRING_NO  = "0"

	/** 流水有效规则：金额大于等于2000元 */
	SALARYNUM = 2000

	/** 流水类型判定天数：15天 */
	BASTYPE = 15

	LimitNumber    = "limitNumber"
	ReportSaveDays = "reportSaveDays"
	TimeSpan       = "timeSpan"

	EMPTYJSON = "{}"
)

//服务配置
const (
	//bcw
	JFAPP_MERCHID               = "0000002"
	App_turnOver_host           = "http://10.100.140.84:8091"
	App_turnOver_bankList       = "/ns-bcwForApp/tradeQueryForApp/api/queryBankInfoList"
	APP_TURNOVER_FIRSTACTION    = "load.htm"
	App_turnOver_historyList    = "/ns-bcwForApp/deviceRelevance/getHistoryList"
	App_turnOver_basInfo        = "/ns-bcwForApp/tradeQueryForApp/api/queryTradeDetailsByCard"
	App_turnOver_queryLimitInfo = "/ns-bcwForApp/userLimit/getUserLimit"

	//bas
	App_bas_host     = "http://10.100.140.125:9001"
	App_bas_loginUrl = "/bas/AppAction/load.htm"
	JQLOGO_PREFIX    = "/bas/pub/images/jq"
)

var SALARYFLAG []string = []string{"工资", "代发工资"}

type BusinessRetInfo struct {
	Code string
	Msg  string
}

var (
	APP_REQ_DATA             = BusinessRetInfo{"1001", "请求参数data必填"}
	APP_REQ_IDENTITYNO       = BusinessRetInfo{"1002", "身份证号不允许为空"}
	APP_REQ_MSGCODE          = BusinessRetInfo{"1003", "银行卡号不允许为空"}
	APP_REQ_DEVICEID         = BusinessRetInfo{"1004", "设备号必填"}
	APP_REQ_OKT              = BusinessRetInfo{"1005", "OKT必填"}
	APP_REQ_REQINFO          = BusinessRetInfo{"1006", "请求信息参数app_service_req必填"}
	APP_SYSTEM_ERROR         = BusinessRetInfo{"2001", "服务器错误！"}
	APP_WRONG_EMPTYDATA      = BusinessRetInfo{"3003", "服务端返回数据为空！"}
	APP_WRONG_EMPTY_GOTOURL  = BusinessRetInfo{"3010", "服务器未返回gotoUrl字段！"}
	APP_WRONG_EMPTY_BANKCODE = BusinessRetInfo{"3011", "服务器未返回bankCode字段！"}
	APP_RETURN_OK            = BusinessRetInfo{"0000", "正常响应"}
	APP_REQ_SIGNINFO         = BusinessRetInfo{"1007", "签名必填"}
	APP_WRONG_SIGNINFO       = BusinessRetInfo{"3001", "签名验证失败"}
	APP_WRONG_NOSIGN         = BusinessRetInfo{"3002", "服务端签名为空"}
	APP_REQ_CUSTNAME         = BusinessRetInfo{"1009", "客户姓名必填"}
	APP_REQ_BANKCODE         = BusinessRetInfo{"1010", "银行简称必填"}
	APP_REQ_REFRESH          = BusinessRetInfo{"1013", "刷新标识必填"}
	APP_WRONG_EMPTY_LIMIT    = BusinessRetInfo{"3019", "服务器返回限制信息为空！"}
	APP_REQ_DATE             = BusinessRetInfo{"1011", "月份必填"}
)
