package atom

var RetCodeMap map[int]string

func init() {
	RetCodeMap = make(map[int]string)
	RetCodeMap[SUCCESS] = "成功"
	RetCodeMap[ERROR_CODE_PARAM] = "参数错误"
	RetCodeMap[ERROR_CODE_NO_TOKEN] = "token为空"
	RetCodeMap[ERROR_CODE_JSON] = "json处理失败"
	RetCodeMap[ERROR_CODE_EXCEPTION] = "处理异常失败"
	RetCodeMap[ERROR_CODE_SMS] = "短信错误"
	RetCodeMap[ERROR_CODE_DATA_NOT_FIND] = "数据不存在"
}

func GetMsgByCode(code int) string {
	if v, ok := RetCodeMap[code]; ok {
		return v
	}
	return ""
}
