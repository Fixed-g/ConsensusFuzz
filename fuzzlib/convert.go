package fuzzlib

import (
	consensuspb "chainmaker.org/chainmaker/pb-go/v2/consensus"
	tbftpb "chainmaker.org/chainmaker/pb-go/v2/consensus/tbft"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/iancoleman/strcase"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strconv"
)

const (
	tagMeta = "meta"
)

// StructToMap convert struct to map[string]interface{}
// if value of field is nil, result doesn't contain related field and value
// for example: input is struct {a: 1, b: "abc", c: nil},
// convert to a--1, b--"abc", result doesn't contain field c
// notice: field must be exported, unexported field will panic
func StructToMap(config interface{}) (map[string]interface{}, error) {

	if config == nil {
		return nil, nil
	}

	if reflect.TypeOf(config).Kind() != reflect.Struct {
		return nil, errors.New("incorrect config type: config type should be struct")
	}

	result := make(map[string]interface{})
	configValue := reflect.ValueOf(config)

	for i := 0; i < configValue.NumField(); i++ {
		field := reflect.TypeOf(config).Field(i)
		if !parseMetaTag(field) {
			continue
		}
		configField := field.Name
		snakeConfigField := strcase.ToSnake(configField)
		rv := configValue.Field(i)
		// changed here
		if rv.Kind() == reflect.Interface {
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Int:
			result[snakeConfigField] = parseInt(rv)
		case reflect.Int8:
			result[snakeConfigField] = parseInt8(rv)
		case reflect.Int16:
			result[snakeConfigField] = parseInt16(rv)
		case reflect.Int32:
			result[snakeConfigField] = parseInt32(rv)
		case reflect.Int64:
			result[snakeConfigField] = parseInt64(rv)
		case reflect.Uint:
			result[snakeConfigField] = parseUint(rv)
		case reflect.Uint8:
			result[snakeConfigField] = parseUint8(rv)
		case reflect.Uint16:
			result[snakeConfigField] = parseUint16(rv)
		case reflect.Uint32:
			result[snakeConfigField] = parseUint32(rv)
		case reflect.Uint64:
			result[snakeConfigField] = parseUint64(rv)
		case reflect.Float32:
			result[snakeConfigField] = parseFloat32(rv)
		case reflect.Float64:
			result[snakeConfigField] = parseFloat64(rv)
		case reflect.String:
			result[snakeConfigField] = parseString(rv)
		case reflect.Bool:
			result[snakeConfigField] = parseBool(rv)
		case reflect.Ptr:
			v, err := parsePtr(rv)
			if err != nil {
				errMsg := errors.New(fmt.Sprintf("structToMap fail, field is %s, value is %v, err is %s",
					configField, rv, err))
				return nil, errMsg
			}
			if v == nil {
				continue
			}
			result[snakeConfigField] = v
		case reflect.Map:
			v := parseMap(rv)
			if v == nil {
				continue
			}
			result[snakeConfigField] = v
		case reflect.Slice:
			v, err := parseSlice(rv)
			if err != nil {
				errMsg := errors.New(fmt.Sprintf("structToMap fail, field is %s, value is %v, err is %s",
					configField, rv, err))
				return nil, errMsg
			}
			if v == nil {
				continue
			}
			result[snakeConfigField] = v
		case reflect.Struct:
			v, err := StructToMap(rv.Interface())
			if err != nil {
				return nil, err
			}
			result[snakeConfigField] = v
		default:
			errMsg := fmt.Sprintf("structToMap fail, unknow value type, type is %s, value is %v\n",
				rv.Kind(), rv)
			return nil, errors.New(errMsg)
		}
	}
	return result, nil
}

func parseInt(value reflect.Value) int {
	return int(value.Int())
}

func parseInt8(value reflect.Value) int8 {
	return int8(value.Int())
}

func parseInt16(value reflect.Value) int16 {
	return int16(value.Int())
}

func parseInt32(value reflect.Value) int32 {
	return int32(value.Int())
}

func parseInt64(value reflect.Value) int64 {
	return value.Int()
}

func parseUint(value reflect.Value) uint {
	return uint(value.Uint())
}

func parseUint8(value reflect.Value) uint8 {
	return uint8(value.Uint())
}

func parseUint16(value reflect.Value) uint16 {
	return uint16(value.Uint())
}

func parseUint32(value reflect.Value) uint32 {
	return uint32(value.Uint())
}

func parseUint64(value reflect.Value) uint64 {
	return value.Uint()
}

func parseString(value reflect.Value) string {
	return value.String()
}

func parseFloat32(value reflect.Value) float32 {
	return float32(value.Float())
}

func parseFloat64(value reflect.Value) float64 {
	return value.Float()
}

func parseBool(value reflect.Value) bool {
	return value.Bool()
}

func parsePtr(v reflect.Value) (map[string]interface{}, error) {
	if v.IsNil() {
		return nil, nil
	}
	v = v.Elem()
	if !v.IsValid() {
		return nil, nil
	}
	return StructToMap(v.Interface())
}

func parseMap(v reflect.Value) map[string]interface{} {
	if v.IsNil() {
		return nil
	}
	result := make(map[string]interface{})
	for _, key := range v.MapKeys() {
		val := v.MapIndex(key)
		result[key.String()] = val.Interface()
	}
	return result
}

func parseSlice(v reflect.Value) ([]interface{}, error) {
	if v.Type().Kind() != reflect.Slice {
		return nil, errors.New("incorrect config type: config type should be slice")
	}

	if v.Len() <= 0 || !v.IsValid() || v.IsNil() {
		return nil, nil
	}

	res := make([]interface{}, v.Len())

	for i := 0; i < v.Len(); i++ {
		switch v.Index(i).Type().Kind() {
		case reflect.Int:
			res[i] = parseInt(v.Index(i))
		case reflect.Int8:
			res[i] = parseInt8(v.Index(i))
		case reflect.Int16:
			res[i] = parseInt16(v.Index(i))
		case reflect.Int32:
			res[i] = parseInt32(v.Index(i))
		case reflect.Int64:
			res[i] = parseInt64(v.Index(i))
		case reflect.Uint:
			res[i] = parseUint(v.Index(i))
		case reflect.Uint8:
			res[i] = parseUint8(v.Index(i))
		case reflect.Uint16:
			res[i] = parseUint16(v.Index(i))
		case reflect.Uint32:
			res[i] = parseUint32(v.Index(i))
		case reflect.Uint64:
			res[i] = parseUint64(v.Index(i))
		case reflect.Float32:
			res[i] = parseFloat32(v.Index(i))
		case reflect.Float64:
			res[i] = parseFloat64(v.Index(i))
		case reflect.String:
			res[i] = parseString(v.Index(i))
		case reflect.Bool:
			res[i] = parseBool(v.Index(i))
		case reflect.Ptr:
			value, err := parsePtr(v.Index(i))
			if err != nil {
				return nil, err
			}
			res[i] = value
		case reflect.Struct:
			curStruct, err := StructToMap(v.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			res[i] = curStruct
		default:
			errMsg := fmt.Sprintf("unknow slice type %s", v.Index(i).Type().Kind().String())
			return nil, errors.New(errMsg)
		}
	}
	return res, nil
}

func parseMetaTag(f reflect.StructField) bool {
	metaTag := f.Tag.Get(tagMeta)
	parseCurrentField, err := strconv.ParseBool(metaTag)
	if err != nil {
		return true
	}
	if parseCurrentField {
		return true
	}
	return false
}

func MapToVoteMsg(m map[string]interface{}) *tbftpb.TBFTMsg {
	vote := &tbftpb.Vote{}
	err := mapstructure.Decode(m["Msg"].(map[string]interface{}), vote)
	if err != nil {
		return nil
	}
	return &tbftpb.TBFTMsg{
		Type: m["Type"].(tbftpb.TBFTMsgType),
		Msg:  mustMarshal(vote),
	}
}

func MapToProposalMsg(m map[string]interface{}) *tbftpb.TBFTMsg {
	proposal := &tbftpb.Proposal{}
	err := mapstructure.Decode(m["Msg"].(map[string]interface{}), proposal)
	if err != nil {
		return nil
	}
	return &tbftpb.TBFTMsg{
		Type: m["Type"].(tbftpb.TBFTMsgType),
		Msg:  mustMarshal(proposal),
	}
}

func MapToTxs(m map[string]interface{}) *consensuspb.RwSetVerifyFailTxs {
	txs := &consensuspb.RwSetVerifyFailTxs{}
	err := mapstructure.Decode(m, txs)
	if err != nil {
		return nil
	}
	return txs
}

func mustMarshal(msg proto.Message) (data []byte) {
	var err error
	defer func() {
		// while first marshal failed, retry marshal again
		if recover() != nil {
			data, err = proto.Marshal(msg)
			if err != nil {
				panic(err)
			}
		}
	}()
	data, err = proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return
}
