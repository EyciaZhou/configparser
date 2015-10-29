package config

import (
	"reflect"
	//"fmt"
	"errors"
	"github.com/op/go-logging"
	"strings"
)

var log = logging.MustGetLogger("netease_news")

func LoadConfDir(_struct interface{}, confdir string) error {
	return nil
}

var (
	ErrorNotStruct = errors.New("excepted a point of struct")
	ErrorNotMatch = errors.New("type not match")
)

func LoadConfFromJson(_struct interface{}, conf map[string]interface{}) error {
	log.Debug("_struct's type is %v", reflect.TypeOf(_struct).Kind())

	if reflect.TypeOf(_struct).Kind() != reflect.Ptr {
		return ErrorNotStruct
	}

	v := reflect.ValueOf(_struct).Elem()
	t := reflect.TypeOf(_struct).Elem()

	log.Debug("point of _struct's type is %v", reflect.TypeOf(v).Kind())

	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return ErrorNotStruct
	}

	log.Debug("_struct is a struct = = pass.")
	log.Debug("starting read tags")

	for i := 0; i < v.NumField(); i++ {

		fieldName := strings.ToLower(t.Field(i).Name)

		log.Debug("---------------------------------------------------------------------------------")
		log.Debug("Field %d:", i)
		log.Debug("FieldName : %s ; Type : %s", fieldName, t.Field(i).Type.Kind())

		if !v.Field(i).CanSet() {
			continue
		}

		switch t.Field(i).Type.Kind() {
		case reflect.String:
			mv, ok := conf[fieldName]
			if !ok || reflect.TypeOf(mv).Kind() != reflect.String {

				if !ok {
					log.Debug("conf not have this field using default value")
					v.Field(i).SetString(t.Field(i).Tag.Get("default"))
					log.Debug("setted : %v", t.Field(i).Tag.Get("default"))
				} else {
					log.Debug("conf's this field's type not match")
					return ErrorNotMatch
				}

			} else {

				log.Debug("this field will be set to %v", mv)

				v.Field(i).SetString(mv.(string))

				log.Debug("setted : %v", v.Field(i))
			}
		case reflect.Int:
		case reflect.Slice:
		default:
			//the type unsupported just skip
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