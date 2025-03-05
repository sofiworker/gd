package ghttp

import (
	"errors"
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
	OnlyWrapFuncError = errors.New("only wrap function")
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

// GinGdWrap 兼容gd框架返回值
func GinGdWrap(f interface{}) func(c *gin.Context) {
	return ginReturnWrap(f, "gd")
}

func ginReturnWrap(f interface{}, kind string) func(c *gin.Context) {
	valueOf := reflect.ValueOf(f)
	typeOf := valueOf.Type()
	if typeOf.Kind() != reflect.Func {
		panic(OnlyWrapFuncError)
	}
	return func(ctx *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				//debug.PrintStack()
				ctx.AbortWithStatusJSON(http.StatusOK, &HttpResponse{
					Code:    http.StatusInternalServerError,
					Message: RetMap[http.StatusInternalServerError],
				})
				return
			}
		}()

		numIn := typeOf.NumIn()
		inValueList := make([]reflect.Value, numIn)

		for idx := 0; idx < numIn; idx++ {
			paramType := typeOf.In(idx)
			if paramType.Kind() == reflect.Ptr && paramType.Elem() == reflect.TypeOf((*gin.Context)(nil)).Elem() {
				inValueList[idx] = reflect.ValueOf(ctx)
				continue
			}
			// 如果是结构体并且为 gin.Context
			if paramType == reflect.TypeOf(gin.Context{}) {
				inValueList[idx] = reflect.ValueOf(*ctx.Copy())
				continue
			}
			value, err := resolveParam(paramType, ctx)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusOK, &HttpResponse{
					Code:    http.StatusInternalServerError,
					Message: RetMap[http.StatusInternalServerError],
				})
				return
			}
			inValueList[idx] = value
		}

		retValues := valueOf.Call(inValueList)
		if len(retValues) == 0 {
			return
		}
		if len(retValues) > 4 {
			ctx.AbortWithStatusJSON(http.StatusOK, &HttpResponse{
				Code:    http.StatusInternalServerError,
				Message: RetMap[http.StatusInternalServerError],
			})
			return
		}

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
			val := value.Interface()
			switch v := val.(type) {
			case error:
				err = v
			case string:
				msg = v
			case int:
				code = v
			default:
				data = v
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

func resolveParam(paramType reflect.Type, c *gin.Context) (reflect.Value, error) {
	switch paramType.Kind() {
	case reflect.Ptr:
		elemType := paramType.Elem()
		// 递归处理元素类型
		elemValue, err := resolveParam(elemType, c)
		if err != nil {
			return reflect.Value{}, err
		}
		// 创建一个新的指针实例并设置其指向的值
		ptrValue := reflect.New(elemType)
		ptrValue.Elem().Set(elemValue)
		return ptrValue, nil
	case reflect.Struct, reflect.Map:
		// 如果是结构体并且为 gin.Context
		if paramType == reflect.TypeOf(gin.Context{}) {
			return reflect.ValueOf(c.Copy()), nil
		}
		instance := reflect.New(paramType).Interface()
		if err := c.ShouldBind(instance); err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(instance).Elem(), nil
	case reflect.Slice, reflect.Array:
		// 对于数组、切片尝试从body中解码json到实例
		instance := reflect.New(paramType).Interface()
		if err := c.ShouldBind(instance); err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(instance).Elem(), nil
	default:
		// 对于其他类型尝试从query、form中获取值
		return reflect.Zero(paramType), nil
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
