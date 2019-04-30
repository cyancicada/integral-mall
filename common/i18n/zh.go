package i18n

const (
	ErrParam  string = "参数错误"
	ErrServer string = "服务器忙碌，请稍后再试"
	TimLayOut string = "2006-01-02 15:04:05"
)

var ZhMessage = map[string]string{
	"RegisterRequest.Mobile.required":   "手机号不能为空",
	"RegisterRequest.Password.required": "密码不能为空",
	"GoodsOrderRequest.Mobile.required": "手机号不能为空",
}
