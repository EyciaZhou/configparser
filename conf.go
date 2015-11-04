package configparser

import (
	"reflect"
	//"fmt"
	"errors"
	//"strings"
	"strconv"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

var (
	ErrorNotStruct = errors.New("excepted a pointer of struct")
)

func ToJson(_struct interface{}) (map[string]interface{}, error) {
	var (
		v reflect.Value
		t reflect.Type
	)

	result := make(map[string]interface{})

	if reflect.TypeOf(_struct).Kind() == reflect.Ptr {
		v = reflect.ValueOf(_struct).Elem()
		t = reflect.TypeOf(_struct).Elem()
	} else {
		v = reflect.ValueOf(_struct)
		t = reflect.TypeOf(_struct)
	}

	if t.Kind() != reflect.Struct {
		panic(errors.New("excepted a pointer of struct or a struct"))
	}

	for i := 0; i < v.NumField(); i++ {
		//fieldName := strings.ToLower(t.Field(i).Name)
		fieldName := t.Field(i).Name
		result[fieldName] = v.Field(i).Interface()
	}

	val, _ := json.MarshalIndent(result, "", "\t")

	fmt.Printf("%v\n", (string)(val))

	return result, nil
}

func LoadConfDefault(_struct interface{}) error {
	v := make(map[string]interface{})
	return LoadConfFromJson(_struct, v)
}

func LoadConfDir(_struct interface{}, confdir string) error {
	b, err := ioutil.ReadFile(confdir)
	if err != nil {
		return err
	}

	v := make(map[string]interface{})

	err = json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	return LoadConfFromJson(_struct, v)
}

func LoadConfFromJson(_struct interface{}, conf map[string]interface{}) error {
	if reflect.TypeOf(_struct).Kind() != reflect.Ptr {
		panic(ErrorNotStruct)
	}

	v := reflect.ValueOf(_struct).Elem()
	t := reflect.TypeOf(_struct).Elem()

	if t.Kind() != reflect.Struct {
		panic(ErrorNotStruct)
	}

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		//fieldName := strings.ToLower(t.Field(i).Name)

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

		case reflect.Int64, reflect.Int:

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

		case reflect.Float64, reflect.Float32:

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
			if ok {
				if reflect.TypeOf(mv).Kind() != reflect.Slice {
					return fmt.Errorf("In field %s, the type not match, should be slice")
				}
			} else {
				s := t.Field(i).Tag.Get("default")
				if s == "" {
					mv = []interface{}{}
				} else {
					err := json.Unmarshal(([]byte)(s), &mv)
					if err != nil {
						return fmt.Errorf("In field %s, unmarshal default error : %s", fieldName, err.Error())
					}
				}
			}
			//so the mv getted

			//tm := reflect.TypeOf(mv)
			vm := reflect.ValueOf(mv)

			tshould := t.Field(i).Type.Elem()

			le := vm.Len()
			tmpTransed := reflect.MakeSlice(reflect.SliceOf(tshould), le, le)

			for i := 0; i < le; i++ {
				if reflect.TypeOf(vm.Index(i).Interface()).ConvertibleTo(tshould) {
					tmpTransed.Index(i).Set( reflect.ValueOf(vm.Index(i).Interface()).Convert(tshould) )
				} else {
					return fmt.Errorf("In field %s, have a Un Convertible value, at postion %d", fieldName, i)
				}
			}

			v.Field(i).Set(tmpTransed)

		default:
			panic(fmt.Errorf("In field %s, unsupported type", fieldName))
		}
	}

	return nil
}
