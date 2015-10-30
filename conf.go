package config

import (
	"reflect"
	//"fmt"
	"errors"
	"github.com/op/go-logging"
	"strings"
	"strconv"
	"fmt"
)

var log = logging.MustGetLogger("netease_news")

func LoadConfDir(_struct interface{}, confdir string) error {
	return nil
}

var (
	ErrorNotStruct = errors.New("excepted a point of struct")
	ErrorNotMatch = errors.New("type not match")
)

var (
	typInt64 = reflect.TypeOf((int64)(1))
	typInt = reflect.TypeOf((int)(1))
)

func LoadConfFromJson(_struct interface{}, conf map[string]interface{}) error {
	if reflect.TypeOf(_struct).Kind() != reflect.Ptr {
		panic(ErrorNotStruct)
	}

	v := reflect.ValueOf(_struct).Elem()
	t := reflect.TypeOf(_struct).Elem()

	if reflect.TypeOf(v).Kind() != reflect.Struct {
		panic(ErrorNotStruct)
	}

	for i := 0; i < v.NumField(); i++ {

		fieldName := strings.ToLower(t.Field(i).Name)

		if !v.Field(i).CanSet() {
			panic(fmt.Errorf("In field %s, unexported field", fieldName))
		}

		switch mv, ok := conf[fieldName]; t.Field(i).Type.Kind() {
		case reflect.String:

			if !ok {
				v.Field(i).SetString(t.Field(i).Tag.Get("default"))
			} else if reflect.TypeOf(mv).Kind() != reflect.String {
				return fmt.Errorf("In field %s, the type not match, should be string")
			} else {
				v.Field(i).SetString(mv.(string))
			}

		case reflect.Int64:

			if !ok {
				var e error

				s := t.Field(i).Tag.Get("default")

				if s == "" {
					v.Field(i).SetInt(0)
					continue
				}

				mv, e = strconv.ParseInt(t.Field(i).Tag.Get("default"), 10, 64)

				if e != nil {
					return fmt.Errorf("In field %s, default string can't be parse to int", fieldName)
				}
			}

			if reflect.TypeOf(mv).ConvertibleTo(t.Field(i).Type) {
				v.Field(i).Set(reflect.ValueOf(mv).Convert(t.Field(i).Type))
			} else {
				return fmt.Errorf("In field %s, the type not match, should be number")
			}

		case reflect.Float64:

			if !ok {
				var e error

				s := t.Field(i).Tag.Get("default")

				if s == "" {
					v.Field(i).SetFloat(0.0)
					continue
				}

				mv, e = strconv.ParseFloat(t.Field(i).Tag.Get("default"), 64)

				if e != nil {
					return fmt.Errorf("In field %s, default string can't be parse to float", fieldName)
				}

			}
			if reflect.TypeOf(mv).ConvertibleTo(t.Field(i).Type) {
				v.Field(i).Set(reflect.ValueOf(mv).Convert(t.Field(i).Type))
			} else {
				return fmt.Errorf("In field %s, the type not match, should be number")
			}

		case reflect.Slice:

			if !ok {
				//damafan==
			} else {

			}

		default:
			panic(fmt.Errorf("In field %s, unexported type", fieldName))
		}
		/*
		if tt.Type.Kind() == reflect.String {
			vv.SetString("aaaaaa")
		}

		if tt.Type.Kind() == reflect.Int {
			vv.SetInt(123231313)
		}

		if tt.Type.Kind() == reflect.Slice {
			log.Debug("this field's type is slice of %v", tt.Type.Elem().Kind())

			vv.Set(reflect.MakeSlice(reflect.SliceOf(tt.Type.Elem()), 10, 10))
		}
		*/
	}

	log.Debug("result : %v", v)

	return nil
}

func loadConfFromJson(_struct interface{}, conf map[string]interface{}) error {
	return nil
}

func init() {
	logging.SetFormatter(logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
	))
}