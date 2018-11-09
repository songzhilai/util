package util

var errorCode = map[int]string{
	0:    "ok",
	-1:   "系统异常，请稍后重试",
	1002: "不合法的参数",
	1003: "kpi数据点参数错误",
	1004: "历史数据时间参数错误",
	1005: "统计类型参数错误",
	1006: "历史数据指标过多，最多10个",
}

// ErrorCode 返回错误代码描述
func ErrorCode(code int) string {
	var codeString = ""
	if _, ok := errorCode[code]; ok {
		codeString = errorCode[code]
	}
	return codeString
}
