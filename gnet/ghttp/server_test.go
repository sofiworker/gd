package ghttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type TestRequest struct {
	Name string `json:"name"`
}

func TestServer(t *testing.T) {
	server := NewServer()
	server.Port = 8080
	server.Address = "127.0.0.1"
	server.Post("/1", func(ctx *gin.Context) {
		fmt.Println(1)
	})
	server.Post("/2", GinJsonWrap(func(ctx gin.Context) {
		fmt.Println(2)
	}))
	server.Post("/21", GinJsonWrap(func(ctx *gin.Context, req *TestRequest) {
		fmt.Println(2, req.Name)
	}))
	server.Get("/3", GinJsonWrap(func() error {
		fmt.Println(3)
		return nil
	}))
	err := server.ListenAndServe()
	if err != nil {
		t.Fatal(err)
	}
}

// 测试数据结构
type TestStruct struct {
	Name string `json:"name"`
}

// 测试ginReturnWrap的参数验证
func TestGinReturnWrap_InvalidParam(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, OnlyWrapFuncError, r)
		}
	}()

	// 传入非函数参数应该触发panic
	ginReturnWrap("not a function", "default")
	t.Error("Expected panic but didn't occur")
}

// 测试panic恢复机制
func TestGinReturnWrap_PanicRecovery(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	// 创建一个会panic的处理函数
	f := func() { panic("test panic") }
	wrapped := ginReturnWrap(f, "default")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	wrapped(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp HttpResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, RetMap[http.StatusInternalServerError], resp.Message)
}

// 测试参数解析
func TestGinReturnWrap_ParamResolution(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	// 测试结构体参数解析
	f := func(s TestStruct) {}
	wrapped := ginReturnWrap(f, "default")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := bytes.NewBufferString(`{"name": "test"}`)
	c.Request = httptest.NewRequest("POST", "/", body)
	c.Request.Header.Set("Content-Type", "application/json")

	wrapped(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 测试返回值处理
func TestGinReturnWrap_ReturnValues(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	tests := []struct {
		name       string
		handler    interface{}
		kind       string
		expectCode int
		expectMsg  string
		expectData interface{}
	}{
		{
			name: "error return",
			handler: func() error {
				return errors.New("test error")
			},
			expectCode: http.StatusInternalServerError,
			expectMsg:  RetMap[http.StatusInternalServerError],
		},
		{
			name: "data return",
			handler: func() interface{} {
				return map[string]string{"key": "value"}
			},
			expectData: map[string]string{"key": "value"},
		},
		{
			name: "code and message",
			handler: func() (int, string) {
				return http.StatusCreated, "custom message"
			},
			expectCode: http.StatusCreated,
			expectMsg:  "custom message",
		},
		{
			name: "gd kind result",
			handler: func() interface{} {
				return "result data"
			},
			kind:       "gd",
			expectData: "result data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := ginReturnWrap(tt.handler, tt.kind)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/", nil)

			wrapped(c)

			var resp HttpResponse
			json.Unmarshal(w.Body.Bytes(), &resp)

			if tt.expectCode != 0 {
				assert.Equal(t, tt.expectCode, resp.Code)
			}
			if tt.expectMsg != "" {
				assert.Equal(t, tt.expectMsg, resp.Message)
			}
			if tt.expectData != nil {
				if tt.kind == "gd" {
					assert.Equal(t, tt.expectData, resp.Result)
				} else {
					assert.Equal(t, tt.expectData, resp.Data)
				}
			}
		})
	}
}

// 测试resolveParam函数
func TestResolveParam(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	tests := []struct {
		name     string
		param    interface{}
		body     string
		expected interface{}
	}{
		{
			name:     "basic type",
			param:    0,
			expected: 0,
		},
		{
			name:     "basic type",
			param:    1.0,
			expected: float64(0),
		},
		{
			name:     "struct type",
			param:    TestStruct{},
			body:     `{"name": "test"}`,
			expected: TestStruct{Name: "test"},
		},
		//{
		//	name:     "struct type",
		//	param:    gin.Context{},
		//	body:     "",
		//	expected: gin.Context{},
		//},
		{
			name:     "pointer type",
			param:    &TestStruct{},
			body:     `{"name": "test"}`,
			expected: &TestStruct{Name: "test"},
		},
		{
			name:     "pointer type",
			param:    &gin.Context{},
			body:     "",
			expected: &gin.Context{},
		},
		{
			name:  "map type",
			param: map[string]string{},
			body:  `{"name": "test"}`,
			expected: map[string]string{
				"name": "test",
			},
		},
		{
			name:  "slice type",
			param: []TestStruct{},
			body:  `[{"name": "test"}]`,
			expected: []TestStruct{
				{Name: "test"},
			},
		},
		{
			name:  "slice type",
			param: []*TestStruct{},
			body:  `[{"name": "test"}]`,
			expected: []*TestStruct{
				{Name: "test"},
			},
		},
		{
			name:  "array type",
			param: [2]TestStruct{},
			body:  `[{"name": "test"}, {"name": "test11111"}]`,
			expected: [2]TestStruct{
				{Name: "test"},
				{Name: "test11111"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paramType := reflect.TypeOf(tt.param)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/", strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")

			value, err := resolveParam(paramType, c)
			if err != nil {
				t.Error(err)
			}
			result := value.Interface()

			assert.Equal(t, tt.expected, result)
		})
	}
}

// 测试错误返回值处理
func TestGinReturnWrap_ErrorHandling(t *testing.T) {
	// 测试错误类型返回
	f := func() error {
		return errors.New("test error")
	}
	wrapped := ginReturnWrap(f, "default")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	wrapped(c)

	var resp HttpResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, RetMap[http.StatusInternalServerError], resp.Message)
}
