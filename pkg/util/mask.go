package util

import (
	"regexp"
)

func Mask(str string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(str) {
		end = len(str)
	}
	mask := ""
	for i := start; i < end; i++ {
		mask += "*"
	}
	return str[:start] + mask + str[end:]
}

// MaskChineseName masks the Chinese name by replacing the characters with asterisks.
//
// It takes a string parameter `name` which represents the Chinese name to be masked.
// It returns a string which is the masked Chinese name.
func MaskChineseName(name string) string {
	reg := regexp.MustCompile(`([\p{Han}]{1})([\p{Han}]*)([\p{Han}]{1})`) // 匹配中文姓名
	return reg.ReplaceAllStringFunc(name, func(s string) string {
		return Mask(s, 1, len(s)-1)
	})
}

func MaskPhone(phone string) string {
	return Mask(phone, 3, 7)
}

func MaskEmail(email string) string {
	reg := regexp.MustCompile(`^([\w\.\-]+)@([\w\-]+\.)+([\w]{2,})$`) // 匹配邮箱
	return reg.ReplaceAllStringFunc(email, func(s string) string {
		parts := reg.FindStringSubmatch(s)
		username := Mask(parts[1], 1, len(parts[1])-1)
		return username + "@" + parts[2] + parts[3]
	})
}

// MaskUrl masks parts of a URL string to anonymize it.
// It takes a URL string as input, uses a regular expression to match the URL format,
// and calls Mask() to replace characters with "*" symbols while preserving the overall structure.
// The characters between the 4th and 2nd last position are masked.
func MaskUrl(url string) string {
	reg := regexp.MustCompile(`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`) // 匹配URL
	return reg.ReplaceAllStringFunc(url, func(s string) string {
		return Mask(s, 4, len(s)-1)
	})
}

func MaskIP(ip string) string {
	reg := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`) // 匹配IP
	return reg.ReplaceAllStringFunc(ip, func(s string) string {
		return Mask(s, 1, len(s)-1)
	})
}

func MaskIDCard(idCard string) string {
	reg := regexp.MustCompile(`^(\d{6})(\d{4})(\d{2})(\d{2})(\d{3})([0-9Xx])$`) // 匹配身份证号码
	return reg.ReplaceAllStringFunc(idCard, func(s string) string {
		return s[:6] + Mask(s[6:], 1, len(s[6:])-1)
	})
}
