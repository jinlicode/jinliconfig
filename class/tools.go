package class

import "regexp"

// CheckDomain 正则匹配输入的域名是否正确
func CheckDomain(domain string) bool {
	re := regexp.MustCompile(`^([\w\-\*]{1,100}\.){1,4}([\w\-]{1,24}|[\w\-]{1,24}\.[\w\-]{1,24})$`)
	if re.MatchString(domain) {
		return true
	}
	return false
}

// CheckEmail 验证邮箱是否正确
func CheckEmail(email string) bool {
	re := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	if re.MatchString(email) {
		return true
	}
	return false
}
