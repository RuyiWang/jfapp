package test

import (
	//	"encoding/json"
	"math/big"
	//	"os"
	"reflect"
	//	"strconv"
	"time"
	//	"encoding/hex"
	//	"crypto/md5"
	"fmt"
	"jfapp/constant"
	"jfapp/util/aescrypt"
	"testing"

	"jfapp/models"
	//	"jfapp/util/oktutil"

	"github.com/bitly/go-simplejson"
	//	"jfapp/mysql"
	"jfapp/util/dateutils"
	"jfapp/util/netutil"

	//	pb "github.com/hyperledger/fabric/protos"
)

func init() {
}

func testAescrypt() {
	//	msg := "{\"deviceId\":\"13927430344\",\"cardId\":\"6227000012860329232\",\"bankCode\":\"CCB\",\"monthCount\":\"24\",\"custName\":\"\",\"endMonth\":\"2016-04\",\"app_client_type\":\"iOS\"}"

	msg2 := "{\"deviceId\":\"13927430344\",\"cardId\":\"6227000012860329232\",\"bankCode\":\"CCB\",\"monthCount\":24,\"custName\":\"\",\"endMonth\":\"2016-04\",\"app_client_type\":[{\"aa\":\"bb\"}]}"

	json, _ := simplejson.NewJson([]byte(msg2))
	json.Array()
	fmt.Println(aescrypt.GenerateSignInfo(json, "aa"))

}

func testNet() {
	url := "http://10.100.140.84:8091/ns-bcwForApp/tradeQueryForApp/api/queryBankInfoList?deviceId=13927430344&merchantId=0000002&userId=13927430344"
	data := ""
	resultStr := netutil.CallInterface(url, data)
	fmt.Println(resultStr)
	resultJson, _ := simplejson.NewJson([]byte(resultStr))

	for _, v := range resultJson.Get(constant.APP_BANK_List).MustArray() {
		if bankCode, ok := v.(map[string]interface{}); ok {
			//			bankCodeStr := (bankCode["bankCode"]).(string)
			fmt.Println(reflect.TypeOf(bankCode))
		}
	}

	slice := make([]*models.TradeData, 5, 10)

	fmt.Println(slice[0] == nil)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))

	bigRat := big.NewRat(0, 1)
	fmt.Println(bigRat.Float64())

	fmt.Println(time.Now().Day())

	fmt.Println(dateutils.GetMonthActualMaximum(time.Now().AddDate(0, -2, 0)))
	queryDate := dateutils.Str2Date("2016-04-19 18:16:58", dateutils.YYYYMMDDHHMMSS)
	fmt.Println(queryDate.Unix())
	fmt.Println(time.Unix(1461061018, 0))
}

func TestMain(t *testing.T) {
	//	fileName := "E:/huhuagaoshou.txt"
	//	fout, _ := os.Create(fileName)
	//	defer fout.Close()

	//	for i := 1; i <= 1475; i++ {
	//		urlBuffer := oktutil.NewStringBufferWithStr("http://www.huhuagaoshouzaidushi.com/")
	//		urlBuffer.AppendStr(strconv.Itoa(i)).AppendStr(".html")
	//		fmt.Println(urlBuffer.ToString())
	//		resultStr := netutil.CallInterface(urlBuffer.ToString(), "")
	//		fmt.Println(resultStr)
	//		fout.WriteString(resultStr)
	//	}
	fmt.Println("aa")
	//	var chaincodeCtorJSON = "{\"Function\":\"init\", \"Args\": [\"a\",\"100\",\"b\", \"200\"]}"
	//	input := &pb.ChaincodeInput{}
	//	err := json.Unmarshal([]byte(chaincodeCtorJSON), &input)
	//	fmt.Println("Chaincode argument error: %s", err)

}
