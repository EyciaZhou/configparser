package configparser

import (
	"reflect"
	"errors"
	"strconv"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"flag"
	"time"
)

var (
	ErrorNotStruct = errors.New("excepted a pointer of struct")
)

func LoadConfFromJson(_struct interface{}, conf map[string]interface{}) {
	handleError(getDefaultMap(_struct))
	handleError(addMap(_struct, conf))
}

func AutoLoadConfig(moduleName string, _struct interface{}) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	err := loadConfPath(_struct, dir + "/conf/" + moduleName + ".conf")
	if err == nil {
		return
	}

	err = loadConfPath(_struct, dir + "/conf/" + moduleName + ".json")
	if err == nil {
		return
	}

	err = loadConfPath(_struct, dir + moduleName + ".conf")
	if err == nil {
		return
	}

	err = loadConfPath(_struct, dir + moduleName + ".json")
	if err == nil {
		return
	}

	err = loadConfPath(_struct, dir + moduleName + "conf.json")
	if err == nil {
		return
	}

	LoadConfDefault(_struct)
}

func LoadConfDefault(_struct interface{}) {
	handleError(getDefaultMap(_struct))
}

func LoadConfPath(_struct interface{}, confpath string) {
	handleError(loadConfPath(_struct, confpath))
}

func LoadConfFromFlag(_struct interface{}) {
	handleError(getDefaultMap(_struct))
	handleError(addFlagMap(_struct))
}


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

func handleError(err error) {
	if err != nil {
		fmt.Println("Error When Load Config\n" + err.Error())
		os.Exit(-2)
	}
}

func loadConfPath(_struct interface{}, confpath string) error {
	b, err := ioutil.ReadFile(confpath)

	handleError(err)

	v := make(map[string]interface{})

	err = json.Unmarshal(b, &v)

	handleError(err)

	if err := getDefaultMap(_struct); err != nil {
		return err
	}
	if err := addMap(_struct, v); err != nil {
		return err
	}

	return nil
}

func checkTypeAndPanic(_struct interface{}) (reflect.Value, reflect.Type) {
	if reflect.TypeOf(_struct).Kind() != reflect.Ptr {
		panic(ErrorNotStruct)
	}

	v := reflect.ValueOf(_struct).Elem()
	t := reflect.TypeOf(_struct).Elem()

	if t.Kind() != reflect.Struct {
		panic(ErrorNotStruct)
	}

	return v, t
}

type sliceValue struct {
	typ reflect.Type
	vals reflect.Value
	setted bool
	field string
}

func (s *sliceValue) String() string {
	return fmt.Sprintf("%v", s.vals)
}

func (sl *sliceValue) Set(s string) error {
	sl.setted = true

	var tmpv interface{}

	if s == "" {
		tmpv = []interface{}{}
	} else {
		err := json.Unmarshal(([]byte)(s), &tmpv)
		if err != nil {
			return fmt.Errorf("In field %s, unmarshal slice error : %s", sl.field, err.Error())
		}
	}

	vm := reflect.ValueOf(tmpv)

	tshould := sl.typ

	le := vm.Len()
	tmpTransed := reflect.MakeSlice(reflect.SliceOf(tshould), le, le)

	for i := 0; i < le; i++ {
		if reflect.TypeOf(vm.Index(i).Interface()).ConvertibleTo(tshould) {
			tmpTransed.Index(i).Set( reflect.ValueOf(vm.Index(i).Interface()).Convert(tshould) )
		} else {
			return fmt.Errorf("In field %s, have a Un Convertible value, at postion %d", sl.field, i)
		}
	}

	sl.vals = tmpTransed

	return nil
}

func addFlagMap(_struct interface{}) error {
	if flag.Parsed() {
		return errors.New("error when add flag arguments: parsed flags cann't parse twice")
	}

	v, t := checkTypeAndPanic(_struct)

	slicePostion := []int{}
	sliceValues := []*sliceValue{}

 	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name

		if !v.Field(i).CanSet() {
			continue;
		}

		pointer := v.Field(i).Addr().Interface()

		usage := t.Field(i).Tag.Get("usage")

		switch t.Field(i).Type.Kind() {

		case reflect.String:
			flag.StringVar(pointer.(*string), fieldName, v.Field(i).String(), usage)

		case reflect.Bool:
			flag.BoolVar(pointer.(*bool), fieldName, v.Field(i).Bool(), usage)

		case reflect.Int:
			flag.IntVar(pointer.(*int), fieldName, (int)(v.Field(i).Int()), usage)

		case reflect.Int64:
			if t.Field(i).Type == reflect.TypeOf(time.Second) {
				flag.DurationVar(pointer.(*time.Duration), fieldName,
					time.Duration(v.Field(i).Int()), usage)
			} else {
				flag.Int64Var(pointer.(*int64), fieldName, v.Field(i).Int(), usage)
			}

		case reflect.Float64:
			flag.Float64Var(pointer.(*float64), fieldName, v.Field(i).Float(), usage)

		case reflect.Slice:
			slicePostion = append(slicePostion, i)
			sliceValues = append(sliceValues, &sliceValue{
				typ: t.Field(i).Type.Elem(),
				field: fieldName,
			})

			flag.Var(sliceValues[len(sliceValues)-1], fieldName, usage)

		default:
			panic(fmt.Errorf("In field %s, unsupported type", fieldName))
		}
	}

	flag.Parse()

	for i, pos := range slicePostion {
		v.Field(pos).Set(sliceValues[i].vals)
	}

	return nil
}

func getDefaultMap(_struct interface{}) error {
	v, t := checkTypeAndPanic(_struct)

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name

		if !v.Field(i).CanSet() {
			continue
		}

		switch raw := t.Field(i).Tag.Get("default"); t.Field(i).Type.Kind() {
		case reflect.String:
			v.Field(i).SetString(raw)
		case reflect.Bool:
			if raw == "true" || raw == "TRUE" || raw == "True" || raw == "T" || raw == "t" || raw == "1" {
				v.Field(i).SetBool(true)
			} else if raw == "" || raw == "false" || raw == "FALSE" || raw == "False" || raw == "F" || raw == "f" || raw == "0" {
				v.Field(i).SetBool(false)
			} else {
				return fmt.Errorf("In field %s, default string can't parse to bool", fieldName)
			}

		case reflect.Int64, reflect.Int:
			if t.Field(i).Type == reflect.TypeOf(time.Duration(0)) {
				//if is time duration, progress by time.Duration

				dur, err := time.ParseDuration(raw)
				if err != nil {
					return fmt.Errorf("In field %s, default string can't parse to time.Duration", fieldName)
				}
				v.Field(i).SetInt((int64)(dur))
			} else {
				if raw == "" {
					v.Field(i).SetInt(0)
					continue
				}

				tmpv, e := strconv.ParseInt(raw, 10, 64)

				if e != nil {
					return fmt.Errorf("In field %s, default string can't be parse to int", fieldName)
				}

				v.Field(i).SetInt(tmpv)
			}

		case reflect.Float64:

			if raw == "" {
				v.Field(i).SetFloat(0.0)
				continue
			}

			tmpv, e := strconv.ParseFloat(t.Field(i).Tag.Get("default"), 64)

			if e != nil {
				return fmt.Errorf("In field %s, default string can't be parse to float", fieldName)
			}

			v.Field(i).SetFloat(tmpv)

		case reflect.Slice:
			sl := &sliceValue{
				typ:t.Field(i).Type.Elem(),
				field:fieldName,
			}

			if e := sl.Set(raw); e != nil {
				return e
			}

			v.Field(i).Set(sl.vals)

		default:
			panic(fmt.Errorf("In field %s, unsupported type", fieldName))
		}
	}

	return nil
}

func isDuration(typ reflect.Type) bool {
	return typ == reflect.TypeOf(time.Duration(1))
}

func addMap(_struct interface{}, conf map[string]interface{}) error {
	v, t := checkTypeAndPanic(_struct)

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name

		if !v.Field(i).CanSet() {
			continue
		}

		mv, ok := conf[fieldName]
		if !ok {
			continue
		}

		switch t.Field(i).Type.Kind() {
		case reflect.String:
			if reflect.TypeOf(mv).Kind() != reflect.String {
				return fmt.Errorf("In field %s, the type not match, should be string", fieldName)
			} else {
				v.Field(i).SetString(mv.(string))
			}

		case reflect.Bool:
			if reflect.TypeOf(mv).Kind() != reflect.Bool {
				return fmt.Errorf("In field %s, the type not match, should be bool", fieldName)
			} else {
				v.Field(i).SetBool(mv.(bool))
			}

		case reflect.Int64, reflect.Int:

			if isDuration(t.Field(i).Type) {
				if reflect.TypeOf(mv).Kind() == reflect.String {
					dur, err := time.ParseDuration(mv.(string))
					if err != nil {
						return fmt.Errorf("In field %s, Can't parse time duration in json: %s", fieldName, err.Error())
					}
					mv = dur
				} else if reflect.TypeOf(mv).Kind() != reflect.Int64 {
					return fmt.Errorf("In field %s, the type not match, should be number or time duration", fieldName)
				}
			}

			if reflect.TypeOf(mv).ConvertibleTo(t.Field(i).Type) {
				v.Field(i).Set(reflect.ValueOf(mv).Convert(t.Field(i).Type))
			} else {
				return fmt.Errorf("In field %s, the type not match, should be number", fieldName)
			}

		case reflect.Float64:
			if reflect.TypeOf(mv).ConvertibleTo(t.Field(i).Type) {
				v.Field(i).Set(reflect.ValueOf(mv).Convert(t.Field(i).Type))
			} else {
				return fmt.Errorf("In field %s, the type not match, should be number", fieldName)
			}

		case reflect.Slice:
			if reflect.TypeOf(mv).Kind() != reflect.Slice {
				return fmt.Errorf("In field %s, the type not match, should be slice", fieldName)
			}

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