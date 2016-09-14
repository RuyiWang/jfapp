package aescrypt

import (
	"encoding/json"
//	"reflect"
	"strings"
	"sort"
	"encoding/hex"
	"crypto/md5"
	"log"
	"crypto/cipher"
	"crypto/aes"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"jfapp/util/oktutil"
	"github.com/bitly/go-simplejson"
)

const (
	LENGTH_BYTE_NUM = 5
	KEY_SPACE = "We should channel strategy to promote paperless into pieces, and plans to be in a three meter package to an app and its outreach efforts to: national, each working day, 7k-1w active;The new settlement platform has developed bank water, into a three meter (identity dimension tables, one table, so and repayment end functional P2P repayment packaging for an app;The initial version of APP 3 functions: query the user credit with water, no paper into pieces; the transfer of WAP checkout, active repayment. The next target: lending and financial information platform, through the accumulation of user credit information, for its accurate to find the loan and financing docking needs of the company."
	IGNORE_FILEDS = "__equalsCalc{]__hashCodeCalc{]serialVersionUID{]signInfo{]seed"

)

var	iv = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}


func Base64Encode(src []byte) string {  
	return base64.StdEncoding.EncodeToString(src)
}  
  
func Base64Decode(src string) ([]byte, error) {  
	return base64.StdEncoding.DecodeString(src)
} 



func getKey(strKey string) []byte {
	byteKey := []byte(strKey)
	key := sha256.Sum256(byteKey)
   	keyLen := len(key)
    if keyLen < 16 {
        panic("res key 长度不能小于16")
    }
    if keyLen >= 32 {
        //取前32个字节
        return key[:32]
    }
    if keyLen >= 24 {
        //取前24个字节
        return key[:24]
    }
//    取前16个字节
    return key[:16]
}



func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
   	padding := blockSize - len(ciphertext)%blockSize
   	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
   	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(plantText []byte, blockSize int) []byte {
   	length := len(plantText)
   	unpadding := int(plantText[length-1])
   	return plantText[:(length - unpadding)]
}


//加密
func Encrypt(key string, src []byte) (string, error) {
	keyBytes := getKey(key)
   	block, err := aes.NewCipher(keyBytes) //选择加密算法
   	if err != nil {
    	return "", err
   	}
   	src = PKCS7Padding(src, block.BlockSize())
   	blockModel := cipher.NewCBCEncrypter(block, iv)
   	ciphertext := make([]byte, len(src))
   	blockModel.CryptBlocks(ciphertext, src)
	res := Base64Encode(ciphertext) 
	
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Runtime error caught: %v", r)
		}
	}()
  	return res, nil
}



//解密
func Decrypt(key string, src []byte) (string, error) {
	keyBytes := getKey(key)
    block, err := aes.NewCipher(keyBytes) //选择加密算法
    if err != nil {
		return "", err
    }
   	blockModel := cipher.NewCBCDecrypter(block, iv)
   	plantText := make([]byte, len(src))
   	blockModel.CryptBlocks(plantText, src)
   	plantText = PKCS7UnPadding(plantText, block.BlockSize())
   	return string(plantText), nil
}

/* 通过Okt 进行解密 */
func DecryptByOkt(appKey, src string) (string, error) {
	srcDecode,_ := Base64Decode(src)
    return Decrypt(appKey, srcDecode)
}

//md5
func MD5Sign(src string) string {
	md5Ctx := md5.New()
    md5Ctx.Write([]byte(src))
    cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func GenerateSignInfo(jsonData *simplejson.Json , key string) string {
	
	data, err := jsonData.Map()
	
	if err != nil {
		panic("GenerateSignInfo failed ")
	}
	
	sorted_keys := make([]string, 0)
    for k, _ := range data {
        sorted_keys = append(sorted_keys, k)
    }
	sort.Strings(sorted_keys)
	var buf *oktutil.StringBuffer = oktutil.NewStringBuffer()
	for _, k := range sorted_keys {
		value := jsonData.Get(k)
		if !strings.Contains(IGNORE_FILEDS,k) {
				if s, ok := (value.Interface()).(string); ok {
					if s != "" {
						buf.AppendStr(k).AppendStr("=")
						buf.AppendStr(s)
						buf.AppendStr("&")
					}
				}  else  if s, ok := (value.Interface()).(json.Number); ok {
					buf.AppendStr(k).AppendStr("=")
					buf.AppendStr(s.String())
					buf.AppendStr("&")
				} else if _, ok := (value.Interface()).([]interface{}); ok {
					buf.AppendStr(k).AppendStr("=")
					byteJson, err := value.MarshalJSON()
					if err == nil {
						buf.AppendStr(string(byteJson))
						buf.AppendStr("&")
					}
				}
		}
    }	
	buf.AppendStr(key).ToString()
	return MD5Sign(buf.ToString())
}
