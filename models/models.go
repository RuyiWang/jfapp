package models

type TradeData struct {
	LastBlance  string `json:"lastBlance"`  //交易后余额
	TradeBlance string `json:"tradeBlance"` //交易金额
	TradeDesc   string `json:"tradeDesc"`   //摘要
	TradeTime   string `json:"tradeTime"`   //交易时间
	QueryTime   string `json:"queryTime"`   //流水入库时间
	TradeType   string `json:"tradeType"`   //交易类型	0 支出 1 收入
	SalaryType  string `json:"salaryType"`  //工资标识	0 非工资 1工资
}

//月交易流水
type MonthTrade struct {
	//收入
	Income string `json:"income"`
	//支出
	Expenses string `json:"expenses"`
	//交易时间：格式（06月2015）
	TradeTime string `json:"tradeTime"`
	//月跨度：格式（0601-0630）
	MonthLenth string `json:"monthLenth"`
	Date       string `json:"date"`
}
