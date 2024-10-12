<h1 align="center">Response Handler for Gin</h1>

<p align="center">
  <i>A response handling package for the Gin framework</i>
</p>

## Overview

This project is a response handling utility package for use with the [Gin](https://github.com/gin-gonic/gin) web framework. It provides a set of tools to create uniform and structured API responses, including error handling, success responses, and logging, making your API more consistent and easier to maintain.

## Features

- **Standardized Error Responses**: Easily return standard HTTP error codes like `BadRequest`, `Unauthorized`, `Forbidden`, etc., with descriptive messages.
- **Logging Capabilities**: Automatically log errors with contextual information, including file name and line number.
- **Convenient API Methods**: Use chainable methods like `WithCode()` and `WithMsg()` to set response codes and messages.

## Installation

```shell
$ go get github.com/jelech/resp
```

## Usage

Here's how you can use the package in a Gin handler:

```go
package resp_test

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jelech/resp"
)

func TestMain(m *testing.M) {
	c := &gin.Context{}

    // Data successfully completed without additional information
    resp.Success(c, resp.OK)

    // Data successfully completed with a custom structure
    resp.Success(c, struct{}{})

    // Print and return an error
    resp.WithMsgLog("this is msg").InternalErr(c)

    // Add information, print, and return if err is not nil
    err := fmt.Errorf("this is test error")
    if resp.WithMsgLog("some error", err).InternalErr(c).Try(err) {
        fmt.Println("error occur")
        // return
    }

    // Print error location and call information, in this case is `response_test.go:24`
    resp.Log(err)

    // Custom code & message
    resp.WithCodeAndMsg(resp.Error{
        Code: 403001,
        Msg:  "project permission denied",
    }).ForbiddenErr(c)
}
```

## API Reference

### Error Structs

- **OK**: Standard success response.
- **BadRequest**: Represents a `400` error.
- **Unauthorized**: Represents a `401` error.
- **Forbidden**: Represents a `403` error.
- **NotFound**: Represents a `404` error.
- **InternalError**: Represents a `500` error.

### Response Methods

- **`WithCode(int) *Response`**: Set custom error code for response.
- **`WithMsg(...interface{}) *Response`**: Set custom message for response.
- **`WithCodeAndMsg(Error) *Response`**: Set both error code and message.
- **`Log(error) *Response`**: Logs the error along with runtime details.
- **`Abort() *Response`**: Aborts the request after returning a response.
- **`Success(*gin.Context, interface{})`**: Return a successful response with optional data.
- **`InternalErr(*gin.Context) *Response`**: Sends a `500 Internal Server Error` response.
- **`ForbiddenErr(*gin.Context) *Response`**: Sends a `403 Forbidden` response.
- **`NotFoundErr(*gin.Context) *Response`**: Sends a `404 Not Found` response.
- **`UnauthorizedErr(*gin.Context) *Response`**: Sends a `401 Unauthorized` response.
- **`BadRequestErr(*gin.Context) *Response`**: Sends a `400 Bad Request` response.

### Logging

The package also provides extensive error logging for debugging purposes. When an error is logged using `Log(error)`, it prints details such as the file, line number, timestamp, and error message. Example:

```go
resp.Log(errors.New("an unexpected error occurred"))
```

## Contribution

If you would like to contribute, please fork the repository and use a feature branch. Pull requests are warmly welcome.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contact

For any questions, issues, or suggestions, feel free to open an issue on GitHub or reach out to the project maintainer.

