package ghttp

import (
	"fmt"
	"github.com/chuck1024/gd/v2/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
	"sync"
)

var (
	RetMap = map[int]string{
		http.StatusBadRequest:          "bad request",
		http.StatusInternalServerError: "internal server error",
		http.StatusOK:                  "ok",
	}
)

type HttpResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

func SuccessJsonResp(c *gin.Context) {
	c.JSON(http.StatusOK, &HttpResponse{
		Code:    http.StatusOK,
		Message: RetMap[http.StatusOK],
	})
}

func SuccessJsonRespWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, &HttpResponse{
		Code:    http.StatusOK,
		Message: msg,
	})
}

func SuccessJsonRespWithData(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, &HttpResponse{
		Code:    http.StatusOK,
		Message: msg,
		Data:    data,
	})
}

func FailedJsonResp(c *gin.Context) {
	c.JSON(http.StatusOK, &HttpResponse{
		Code:    http.StatusInternalServerError,
		Message: RetMap[http.StatusInternalServerError],
	})
}

func FailedJsonRespWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, &HttpResponse{
		Code:    http.StatusInternalServerError,
		Message: msg,
	})
}

func FailedJsonRespWithData(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, &HttpResponse{
		Code:    http.StatusInternalServerError,
		Message: msg,
		Data:    data,
	})
}

func GinJsonResp(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, &HttpResponse{
		Code:    code,
		Message: msg,
	})
}

func GinJsonWrap(f interface{}) func(c *gin.Context) {
	return ginReturnWrap(f, "")
}

func getStructOrSliceFromBody(c *gin.Context, rType reflect.Type) (reflect.Value, error) {
	method := c.Request.Method
	switch method {
	case http.MethodGet:
		// 打印告警日志
	case http.MethodPost:
		if c.Request.ContentLength == 0 {
			return reflect.Value{}, fmt.Errorf("the request body length is 0")
		}
	}

	switch rType.Kind() {
	case reflect.Struct, reflect.Array, reflect.Slice:
	default:
		return reflect.Value{}, fmt.Errorf("only support struct data")
	}

	data := reflect.New(rType)
	err := c.BindJSON(data.Interface())
	if err != nil {
		// 是否为验证错误
		logger.Errorf("bind json failed:%v", err)
		return reflect.Value{}, err
	}
	return data, nil
}

// GinGdWrap 兼容gd框架返回值
func GinGdWrap(f interface{}) func(c *gin.Context) {
	return ginReturnWrap(f, "gd")
}

func ginReturnWrap(f interface{}, kind string) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				ctx.AbortWithStatusJSON(http.StatusOK, &HttpResponse{
					Code:    http.StatusInternalServerError,
					Message: RetMap[http.StatusInternalServerError],
				})
				return
			}
		}()
		valueOf := reflect.ValueOf(f)
		typeOf := valueOf.Type()
		if typeOf.Kind() != reflect.Func {
			return
		}

		numIn := typeOf.NumIn()

		inValueList := make([]reflect.Value, numIn)
		for idx := 0; idx < numIn; idx++ {
			in := typeOf.In(idx)
			var value reflect.Value
		InPtr:
			switch in.Kind() {
			case reflect.Ptr:
				// 如果是指针并且为 gin.Context
				if in.Elem() == reflect.TypeOf((*gin.Context)(nil)).Elem() {
					value = reflect.ValueOf(ctx)
					break
				}
				in = in.Elem()
				goto InPtr
			case reflect.Struct:
				// 如果是结构体并且为 gin.Context
				if in == reflect.TypeOf(gin.Context{}) {
					value = reflect.ValueOf(*ctx.Copy())
					break
				}
				data, err := getStructOrSliceFromBody(ctx, in)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusOK, &HttpResponse{
						Code:    http.StatusBadRequest,
						Message: RetMap[http.StatusBadRequest],
					})
					return
				}
				value = data
			case reflect.Array, reflect.Slice:
				data, err := getStructOrSliceFromBody(ctx, in)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusOK, &HttpResponse{
						Code:    http.StatusBadRequest,
						Message: RetMap[http.StatusBadRequest],
					})
					return
				}
				value = data
			default:
				// 其他类型
				value = reflect.Zero(in)
			}
			inValueList[idx] = value
		}

		retValues := valueOf.Call(inValueList)

		resp := &HttpResponse{
			Code:    http.StatusOK,
			Message: RetMap[http.StatusOK],
		}

		var (
			err  error
			msg  string
			data interface{}
			code int
		)
		for _, value := range retValues {
			retType := value.Type()
		RetPtr:
			switch retType.Kind() {
			case reflect.Ptr:
				retType = retType.Elem()
				value = value.Elem()
				goto RetPtr
			case reflect.Interface:
				if retType.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
					retErr, ok := value.Interface().(error)
					if !ok {
						err = fmt.Errorf("convert error")
					}
					err = retErr
				}
			case reflect.Array, reflect.Slice, reflect.Struct:
				data = value.Interface()
			case reflect.String:
				msg = value.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if value.CanInt() {
					i := value.Int()
					code = int(i)
				}
				if value.CanUint() {
					u := value.Uint()
					code = int(u)
				}
			default:
			}
		}

		if err != nil {
			r := &HttpResponse{
				Code:    http.StatusInternalServerError,
				Message: RetMap[http.StatusInternalServerError],
			}
			//var e error
			//if errors.As(e, io.EOF) {
			//	r.Code = http.StatusBadRequest
			//	r.Message = RetMap[http.StatusBadRequest]
			//}
			ctx.AbortWithStatusJSON(http.StatusOK, r)
			return
		}

		if code != 0 {
			resp.Code = code
		}

		if msg != "" {
			resp.Message = msg
		}

		if data != nil {
			resp.Data = data
			switch kind {
			case "gd":
				resp.Result = data
			default:
			}
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

type Validator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &Validator{}

func (v *Validator) ValidateStruct(obj interface{}) error {
	v.lazyInit()
	valueOf := reflect.ValueOf(obj)
	if valueOf.CanInterface() {
		obj = valueOf.Interface()
	}
	kind := kindOfData(obj)
	if kind == reflect.Struct {
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}

	return nil
}

func (v *Validator) Engine() interface{} {
	v.lazyInit()
	return v.validate
}

func (v *Validator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("validate")
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
