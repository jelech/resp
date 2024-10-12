package resp

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Code int    `json:"code"` // 错误code
	Msg  string `json:"msg"`  // 错误信息
}

var (
	OK            = Error{Code: 0, Msg: "OK"}
	BadRequest    = Error{Code: 0, Msg: "Bad Request"}
	Unauthorized  = Error{Code: 0, Msg: "Unauthorized"}
	Forbidden     = Error{Code: 0, Msg: "Forbidden"}
	NotFound      = Error{Code: 0, Msg: "Not found"}
	InternalError = Error{Code: 0, Msg: "Internal Error"}
)

type Response struct {
	Error
	aborted bool
}

func newResponse() *Response {
	return &Response{}
}

func (e *Response) defaultData(data interface{}) interface{} {
	if data == nil {
		data = Error{e.Code, e.Msg}
	}
	return data
}

func (e *Response) coverData(data Error) interface{} {
	if e.Code != 0 {
		data.Code = e.Code
	}

	if e.Msg != "" {
		data.Msg = e.Msg
	}
	return data
}

func (e *Response) finishContext(c *gin.Context, statusCode int, data interface{}) {
	if e.aborted {
		c.AbortWithStatusJSON(statusCode, data)
		return
	}
	c.JSON(statusCode, data)
}

func (e *Response) WithCode(c int) *Response {
	e.Code = c
	return e
}

func (e *Response) WithMsg(s ...interface{}) *Response {
	e.Msg = fmt.Sprint(s...)
	return e
}

func (e *Response) WithCodeAndMsg(ec Error) *Response {
	e.Code = ec.Code
	e.Msg = ec.Msg
	return e
}

func (e *Response) Log(err error) *Response {
	caller, file, line, ok := runtime.Caller(1)
	fmt.Print(caller, file, line, ok, err)
	return e
}

func (e *Response) Abort() *Response {
	e.aborted = true
	return e
}

func (e *Response) Try(err interface{}) bool {
	return err != nil
}

func (e *Response) InternalErr(c *gin.Context) *Response {
	e.aborted = true
	e.finishContext(c, http.StatusInternalServerError, e.coverData(InternalError))
	return e
}

func (e *Response) ForbiddenErr(c *gin.Context) *Response {
	e.aborted = true
	e.finishContext(c, http.StatusForbidden, e.coverData(Forbidden))
	return e
}

func (e *Response) NotFoundErr(c *gin.Context) *Response {
	e.aborted = true
	e.finishContext(c, http.StatusNotFound, e.coverData(NotFound))
	return e
}

func (e *Response) UnauthorizedErr(c *gin.Context) *Response {
	e.aborted = true
	e.finishContext(c, http.StatusUnauthorized, e.coverData(Unauthorized))
	return e
}

func (e *Response) BadRequestErr(c *gin.Context) *Response {
	e.aborted = true
	e.finishContext(c, http.StatusBadRequest, e.coverData(BadRequest))
	return e
}

func (e *Response) Success(c *gin.Context, data interface{}) {
	data = e.defaultData(data)

	e.finishContext(c, http.StatusOK, data)
}

func (e *Response) CheckInternalErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	e.aborted = true
	e.finishContext(c, http.StatusInternalServerError, e.coverData(InternalError))
	return true
}

// WithCode 会将改errCode绑定在返回的response中，不可用于Success返回
func WithCode(c int) *Response { r := newResponse(); r.Code = c; return r }

// WithMsg 会将这些信息包装到返回的response.msg中，不可用于Success返回
func WithMsg(s ...interface{}) *Response { r := newResponse(); r.Msg = fmt.Sprint(s...); return r }

func WithCodeAndMsg(e Error) *Response {
	r := newResponse()
	r.Code = e.Code
	r.Msg = e.Msg
	return r
}

func WithMsgLog(s ...interface{}) *Response {
	r := newResponse()
	r.Msg = fmt.Sprint(s...)
	if len(s) == 1 && s[0] == nil {
		return r
	}

	var err = errors.New(fmt.Sprintln(s...))
	caller, file, line, _ := runtime.Caller(1)
	fmt.Printf("[%v] [%v:%v:%v] %v", time.Now().Unix(), caller, file, line, err)
	return r
}

// Log 会将error信息打印出来，包括调用的位置、时间等信息
func Log(err error) *Response {
	r := newResponse()

	caller, file, line, _ := runtime.Caller(1)
	fmt.Printf("[%v] [%v:%v:%v] %v", time.Now().Unix(), caller, file, line, err)
	return r
}

func Success(ctx *gin.Context, data interface{}) { newResponse().Success(ctx, data) }

func InternalErr(c *gin.Context) *Response { return newResponse().InternalErr(c) }

func ForbiddenErr(c *gin.Context) *Response { return newResponse().ForbiddenErr(c) }

func NotFoundErr(c *gin.Context) *Response { return newResponse().NotFoundErr(c) }

func UnauthorizedErr(c *gin.Context) *Response { return newResponse().UnauthorizedErr(c) }

func BadRequestErr(c *gin.Context) *Response { return newResponse().BadRequestErr(c) }
