package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"golang.org/x/exp/constraints"
)

// IsPhoneNumber checks if the given string is a valid phone number.
//
// phoneNumber: the string to be checked.
// Returns: a boolean value indicating if the string is a valid phone number.
//
// Example:
//
//	IsPhoneNumber("1234567890") // false
func IsPhoneNumber(phoneNumber string) bool {
	// 定义手机号码正则表达式
	pattern := `^(1[3-9])\d{9}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(phoneNumber)
}

// MD5 calculates the MD5 hash of the given text.
//
// It takes a string parameter called "text" which represents the text to be hashed.
// The function returns a string which represents the hexadecimal representation of the MD5 hash.
//
// Example:
//
//	MD5("Hello, World!") // "b10a8db164e0754105b7a99be72e3fe5"
func MD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func MD5Bytes(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

// IsChinese checks if the given string contains only Chinese characters.
//
// Parameter:
// str - the string to be checked.
//
// Return:
// bool - true if the string contains only Chinese characters, false otherwise.
//
// Example:
//
//	IsChinese("你好") // true
//	IsChinese("Hello") // false
func IsChinese(str string) bool {
	reg := regexp.MustCompile(`^[\u4e00-\u9fa5]+$`)
	return reg.MatchString(str)
}

func ToMap(s interface{}) (m map[string]interface{}, err error) {
	return cast.ToStringMapE(s)
}

// MapSlice函数，接受一个数组和一个映射函数f，返回一个新的数组
//
// Example:
//
//	MapSlice([]int{1, 2, 3}, func(x int) int { return x + 1 }) // [2, 3, 4]
//	MapSlice([]string{"a", "b", "c"}, func(x string) string { return strings.ToUpper(x) }) // ["A", "B", "C"]
func MapSlice[T any, U any](nums []T, f func(T) U) []U {
	result := make([]U, len(nums))
	for i, num := range nums {
		result[i] = f(num)
	}
	return result
}

func SliceToMap[T any, K comparable](slice []T, keyFunc func(T) K) map[K]T {
	result := make(map[K]T, len(slice))
	for _, item := range slice {
		result[keyFunc(item)] = item
	}
	return result
}

// Reduce函数，接受一个整数数组和一个归约函数f，返回归约结果
func Reduce(nums []int, f func(int, int) int, init int) int {
	result := init
	for _, num := range nums {
		result = f(result, num)
	}
	return result
}

// 生成指定长度的随机数字字符串
//
// Example:
//
//	GenerateCode(6) // "123456"
func GenerateCode(length int) string {
	code := ""
	for i := 0; i < length; i++ {
		code += fmt.Sprintf("%d", rand.Intn(10))
	}
	return code
}

// CalculateChecksum 计算给定请求的校验和
// 它将请求转换为字符串,计算 MD5 哈希,并将哈希转换为十六进制编码的字符串
func CalculateChecksum(request interface{}) string {
	// 使用 json.Marshal 来序列化请求
	data, err := json.Marshal(request)
	if err != nil {
		// 如果序列化失败，回退到原来的方法
		data = []byte(fmt.Sprintf("%v", request))
	}

	// 直接使用 MD5Bytes 函数
	return MD5Bytes(data)
}

// CompareRequests compares the checksum of multiple requests.
//
// It takes a slice of requests as input and returns a boolean value indicating if the checksums of all requests are equal.
//
// Example:
//
//	 type MyRequest struct {
//	  Name string
//	  Age  int
//	}
//
//	r1 := MyRequest{Name: "John", Age: 30}
//	r2 := MyRequest{Name: "John", Age: 30}
//	r3 := MyRequest{Name: "Jane", Age: 20}
//	CompareRequests(r1, r2, r3) // true
//	CompareRequests(r1, r2, r3, r3) // false
//	CompareRequests(r1, r3, r2) // true
//	CompareRequests(r1, r2, r3, r3) // false
//	CompareRequests(r1, r2, r3) // true
func CompareRequests(requests ...interface{}) bool {
	if len(requests) <= 1 {
		return true // No need to compare if there's only one request
	}

	firstChecksum := CalculateChecksum(requests[0])

	for _, request := range requests[1:] {
		checksum := CalculateChecksum(request)
		if checksum != firstChecksum {
			return false
		}
	}

	return true
}

// GenerateRandomNumber 生成指定范围内的随机整数
func GenerateRandomNumber(min, max int) int {
	if min >= max {
		panic("min must be less than max")
	}
	return rand.Intn(max-min+1) + min
}

// IsZeroStruct checks if the given struct is empty.
//
// Example:
//
//	 type MyStruct struct {
//	  Name string
//	  Age  int
//	}
//
//	s := MyStruct{}
//	IsZeroStruct(s) // true
//	s.Name = "John"
//	IsZeroStruct(s) // false
func IsZeroStruct(s any) bool {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			if !v.Field(i).IsZero() {
				return false
			}
		}
	}
	return true
}

// SetDefaults sets default values for struct fields tagged with "default"
// by reflecting over the struct. It handles setting defaults for string,
// int, float64 and bool struct fields based on the tag value.
//
// Example:
//
//	type MyStruct struct {
//	  Name string `default:"John"`
//	  Age  int    `default:"30"`
//	  Enabled bool `default:"true"`
//	}
//	SetDefaults(&MyStruct{})
//	// MyStruct will be set to:
//	MyStruct{Name: "John", Age: 30, Enabled: true}
func SetDefaults(data interface{}) {
	value := reflect.ValueOf(data).Elem()
	typ := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		tag := typ.Field(i).Tag.Get("default")
		if tag == "" || !field.IsZero() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(tag)
		case reflect.Int:
			intValue, _ := strconv.Atoi(tag)
			field.SetInt(int64(intValue))
		case reflect.Float64:
			v, _ := strconv.ParseFloat(tag, 64)
			field.SetFloat(v)
		case reflect.Bool:
			field.SetBool(tag == "true")
		default:
			panic(fmt.Sprintf("unsupported type: %s", field.Kind()))
		}
	}
}

// GenerateRandomString generates a random string of the given length.
// It does this by selecting random bytes from the set of alphanumeric characters.
func GenerateRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

// ToString converts the given interface{} value to a string.
// If the value cannot be converted to a string, an empty string is returned.
func ToString(data interface{}) string {
	return cast.ToString(data)
}

// ToJSON converts the given interface{} value to a JSON string.
func ToJSON(data interface{}) string {
	jsonBytes, _ := json.Marshal(data)

	return string(jsonBytes)
}

// ToInt converts the given interface{} value to an int.
// If the value cannot be converted to an int, 0 is returned.
func ToInt(data interface{}) int {
	return cast.ToInt(data)
}

// ToFloat64 converts the given interface{} value to a float64.
// If the value cannot be converted to a float64, 0.0 is returned.
func ToFloat64(data interface{}) float64 {
	return cast.ToFloat64(data)
}

// GetEnv returns the value of the environment variable named by the key.
// If the environment variable is not present, the fallback value is returned instead.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetUUID() string {
	return uuid.New().String()
}

// GenerateOrderNumber 生成订单号
func GenerateOrderNumber() (string, error) {
	node, err := snowflake.NewNode(1) // 1 是节点ID
	if err != nil {
		return "", err
	}
	id := node.Generate()
	return id.String(), nil
}

// ClearDirectory 清理指定目录中的所有文件
func ClearDirectory(dirPath string) error {
	d, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dirPath, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteOldFiles 删除超过指定天数的文件
func DeleteOldFiles(logDir string, days int) error {
	files, err := os.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("无法读取目录: %v", err)
	}

	// 获取当前时间
	now := time.Now()

	for _, file := range files {
		// 获取文件的完整路径
		filePath := filepath.Join(logDir, file.Name())

		// 获取文件信息
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Errorf("无法获取文件信息: %v", err)
			continue
		}

		// 计算文件的修改时间和当前时间的差值
		diff := now.Sub(fileInfo.ModTime())

		// 如果文件修改时间超过指定天数，则删除该文件
		if diff.Hours() > float64(days*24) {
			err := os.Remove(filePath)
			if err != nil {
				log.Errorf("无法删除文件: %v", err)
			} else {
				log.Infof("删除文件: %s\n", filePath)
			}
		}
	}

	return nil
}

// 将一种整数类型的切片转换为另一种整数类型的切片
//
//	使用示例:
//	ToIntSlice := ToSlice[int64, int]
//	ToInt64Slice := ToSlice[int32, int64]
func ToSlice[T, U constraints.Integer](slice []T) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = U(v)
	}
	return result
}

func ToPointer[T any](value T) *T {
	return &value
}

// 获取某个月份的最后一天的日期
func GetLastDayOfMonth(year int, month int) time.Time {
	return time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.Local).AddDate(0, 0, -1)
}

// GetDate 获取日期
func GetDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetDateStr 获取日期字符串
func GetDateStr(t time.Time) string {
	return t.Format("2006-01-02")
}

// 生成随机邮箱地址
func GenerateEmailAddress() string {
	username := GenerateRandomString(rand.Intn(10) + 5) // 随机生成5到15位的用户名
	domains := []string{"gmail.com", "yahoo.com", "hotmail.com", "outlook.com"}
	domain := domains[rand.Intn(len(domains))] // 随机选择一个域名

	return fmt.Sprintf("%s@%s", username, domain)
}

func PadNumber[T constraints.Integer](number T, length int) string {
	// 将数字转换为字符串
	numStr := strconv.Itoa(ToInt(number))

	// 计算需要补0的数量
	padding := length - len(numStr)

	// 如果需要补0的数量大于0，则在前面补0
	if padding > 0 {
		numStr = strings.Repeat("0", padding) + numStr
	}

	return numStr
}

func Timer[T any, R any](f func(T) R) func(T) R {
	return func(arg T) R {
		start := time.Now()
		result := f(arg)
		duration := time.Since(start)
		fmt.Printf("函数运行时间: %v\n", duration)
		return result
	}
}

// 定义一个类型，用于表示无参无返回值的函数
type Function func()

// MeasureTime 用于测量函数的运行时间
func MeasureTime(fn Function) time.Duration {
	start := time.Now()          // 获取当前时间
	fn()                         // 执行传入的函数
	elapsed := time.Since(start) // 计算运行时间
	return elapsed
}

// ContainsAny checks if the given string `s` contains any of the substrings in the `substrs` slice.
// It returns true if at least one of the substrings is found in `s`, and false otherwise.
func ContainsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// GenerateDates 生成两个日期之间的所有日期（包括这两个日期）日期格式："2006-01-02"
func GenerateDates(start, end string) ([]string, error) {
	// 解析开始和结束日期
	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		return nil, err
	}

	// 创建一个空的日期切片
	dates := []string{}

	// 循环从startDate到endDate，包括endDate
	current := startDate
	for !current.After(endDate) {
		// 将日期添加到切片中，格式为YYYY-MM-DD
		dates = append(dates, current.Format("2006-01-02"))

		// 移到下一个日期
		current = current.AddDate(0, 0, 1)
	}

	return dates, nil
}

// FormatDateWithTime 将string"2006-01-02"转换为"01-02"
func FormatDateWithTime(dateStr string) string {
	// 解析日期字符串
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}
	// 格式化日期为"MM-DD"
	formattedDate := parsedDate.Format("01-02")
	return formattedDate
}

// MapDecode 将map[string]any转换为结构体
func MapDecode(data map[string]any, v any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, v)
}

// GroupBy groups a slice of structs by a key function.
//
// Example:
//
//	type MyStruct struct {
//	  ID   int
//	  Name string
//	}
//
//	data := []MyStruct{{ID: 1, Name: "John"}, {ID: 2, Name: "Jane"}, {ID: 1, Name: "Jack"}}
//	result := GroupBy(data, func(s MyStruct) int { return s.ID })
//	// result will be map[int][]MyStruct{1: {{ID: 1, Name: "John"}, {ID: 1, Name: "Jack"}}, 2: {{ID: 2, Name: "Jane"}}}
func GroupBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range slice {
		key := keyFunc(item)
		result[key] = append(result[key], item)
	}
	return result
}
