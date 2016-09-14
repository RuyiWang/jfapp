package oktutil

import (
	"bytes"
	"math/rand"
	"time"
)

const (
	LENGTH_BYTE_NUM = 5
)

func MakeRandomLength() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	a := (1 << (LENGTH_BYTE_NUM - 1)) - 1
	b := r.Float64()
	res := (b * float64(a)) + float64(a)
	return int(res)
}

func MakeRandomIndex(maxNum int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	a := maxNum - (1 << LENGTH_BYTE_NUM)
	b := r.Float64()
	res := b * float64(a)
	return int(res)
}

func MakeOKTCode(index, length int) int {
	return (index << LENGTH_BYTE_NUM) | length
}

func GetIndex(OKTCode int) int {
	return OKTCode >> LENGTH_BYTE_NUM
}

func GetLength(OKTCode int) int {
	return OKTCode & ((1 << LENGTH_BYTE_NUM) - 1)
}

//截取字符串 start 起点下标 end 终点下标(不包括)
func Substr2(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

//截取字符串 start 起点下标 length 需要截取的长度
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

type StringBuffer struct {
	buffer *bytes.Buffer
}

func NewStringBufferWithStr(src string) *StringBuffer {
	var buffer *bytes.Buffer = bytes.NewBufferString(src)
	return &StringBuffer{buffer}
}

func NewStringBuffer() *StringBuffer {
	var buffer *bytes.Buffer = bytes.NewBufferString("")
	return &StringBuffer{buffer}
}

func (this *StringBuffer) ToString() string {
	return this.buffer.String()
}

func (this *StringBuffer) AppendStr(str string) *StringBuffer {
	this.buffer.WriteString(str)
	return this
}
