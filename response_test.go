package resp_test

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jelech/resp"
)

func TestMain(m *testing.M) {
	c := &gin.Context{}

	// 数据正常完成，没有其他信息返回
	resp.Success(c, resp.OK)

	// 数据正常完成，有自定义结构返回
	resp.Success(c, struct{}{})

	// 打印并返回错误
	resp.WithMsgLog("this is msg").InternalErr(c)

	// 添加信息，打印，如果err不为nil则返回
	err := fmt.Errorf("this is test error")
	if resp.WithMsgLog("some error", err).InternalErr(c).Try(err) {
		fmt.Println("error occur")
		// return
	}

	// 打印错误位置调用信息，这里为 response_test.go:24
	resp.Log(err)

	// 自定义code&msg
	resp.WithCodeAndMsg(resp.Error{
		Code: 403001,
		Msg:  "用户名",
	}).ForbiddenErr(c)

}
