package httpx

import "github.com/gin-gonic/gin"

// ErrorBody là định dạng lỗi thống nhất trả về cho client.
type ErrorBody struct {
	Error   string `json:"error"`             // thông báo lỗi cho người dùng
	Message string `json:"message,omitempty"` // mô tả chi tiết (tùy chọn)
}

// Error trả về lỗi với mã trạng thái HTTP cho trước.
func Error(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, ErrorBody{Error: message})
}

// ErrorWithDetail trả về lỗi kèm mô tả chi tiết.
func ErrorWithDetail(c *gin.Context, status int, message, detail string) {
	c.AbortWithStatusJSON(status, ErrorBody{Error: message, Message: detail})
}

// ErrorBodyWithCode là định dạng lỗi mở rộng, dùng khi client cần phân biệt
// loại lỗi bằng một mã `code` ổn định (ví dụ: INACTIVE_CODE_EXISTS) và cần
// kèm dữ liệu bổ sung (ví dụ: bản ghi cũ đã bị xóa mềm).
type ErrorBodyWithCode struct {
	Error   string      `json:"error"`
	Message string      `json:"message,omitempty"`
	Code    string      `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ConflictWithCode trả về HTTP 409 kèm mã lỗi và dữ liệu bổ sung cho client.
// Dùng khi mã đã tồn tại ở trạng thái ngừng hoạt động và client cần hỏi
// người dùng có muốn khôi phục bản ghi cũ hay không.
func ConflictWithCode(c *gin.Context, message, code string, data interface{}) {
	c.AbortWithStatusJSON(409, ErrorBodyWithCode{
		Error: message,
		Code:  code,
		Data:  data,
	})
}

// OK trả về dữ liệu thành công với HTTP 200.
func OK(c *gin.Context, data interface{}) {
	c.JSON(200, data)
}

// Created trả về dữ liệu vừa tạo với HTTP 201.
func Created(c *gin.Context, data interface{}) {
	c.JSON(201, data)
}
