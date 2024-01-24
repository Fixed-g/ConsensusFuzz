package tbft

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	tbftpb "chainmaker.org/chainmaker/pb-go/v2/consensus/tbft"
)

type Numeric interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

var MAX_MUTATE_ITER = 5

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

func Generate_random_bool() bool {
	return seededRand.Intn(2) != 0
}

func generate_edge_value_with_type(number interface{}) (interface{}, interface{}) {
	switch number.(type) {
	case uint:
		return uint(0), uint(math.MaxUint32)
	case uint8:
		return uint8(0), uint8(math.MaxUint8)
	case uint16:
		return uint16(0), uint16(math.MaxUint16)
	case uint32:
		return uint32(0), uint32(math.MaxUint32)
	case uint64:
		return uint64(0), uint64(math.MaxUint64)
	case int:
		return int(math.MinInt32), int(math.MaxInt32)
	case int8:
		return int8(math.MinInt8), int8(math.MaxInt8)
	case int16:
		return int16(math.MinInt16), int16(math.MaxInt16)
	case int32:
		return int32(math.MinInt32), int32(math.MaxInt32)
	case int64:
		return int64(math.MinInt64), int64(math.MaxInt64)
	case float32:
		return -math.MaxFloat32, math.MaxFloat32
	case float64:
		return float64(-math.MaxFloat64), float64(math.MaxFloat64)
	default:
		fmt.Printf("%s is an unknown number type!\n", reflect.TypeOf(number).String())
		return false, false
	}
}

func generate_random_number_with_range(lower interface{}, upper interface{}) interface{} {
	switch lower.(type) {
	case uint:
		return uint(generate_random_unsigned_integer_number_with_range(uint64(lower.(uint)), uint64(upper.(uint))))
	case uint8:
		return uint8(generate_random_unsigned_integer_number_with_range(uint64(lower.(uint8)), uint64(upper.(uint8))))
	case uint16:
		return uint16(generate_random_unsigned_integer_number_with_range(uint64(lower.(uint16)), uint64(upper.(uint16))))
	case uint32:
		return uint32(generate_random_unsigned_integer_number_with_range(uint64(lower.(uint32)), uint64(upper.(uint32))))
	case uint64:
		return uint64(generate_random_unsigned_integer_number_with_range(uint64(lower.(uint64)), uint64(upper.(uint64))))
	case int:
		return int(generate_random_signed_integer_number_with_range(int(lower.(int)), int(upper.(int))))
	case int8:
		return int8(generate_random_signed_integer_number_with_range(int(lower.(int8)), int(upper.(int8))))
	case int16:
		return int16(generate_random_signed_integer_number_with_range(int(lower.(int16)), int(upper.(int16))))
	case int32:
		return int32(generate_random_signed_integer_number_with_range(int(lower.(int32)), int(upper.(int32))))
	case int64:
		return int64(generate_random_signed_long_integer_number_with_range(int64(lower.(int64)), int64(upper.(int64))))
	case float32:
		return generate_random_float32_number_with_range(lower.(float32), upper.(float32))
	case float64:
		return generate_random_float64_number_with_range(lower.(float64), upper.(float64))
	default:
		fmt.Printf("%s is an unknown number type!\n", reflect.TypeOf(lower).String())
		return false
	}
}

func generate_random_float32_number_with_range(min, max float32) float32 {
	return seededRand.Float32()*(max-min) + min
}

func generate_random_float64_number_with_range(min, max float64) float64 {
	return seededRand.Float64()*(max-min) + min
}

func generate_random_signed_integer_number_with_range(min, max int) int {
	return int(seededRand.Int63n(int64(max)-int64(min)) + int64(min))
}

func generate_random_signed_long_integer_number_with_range(min, max int64) int64 {
	n := uint64(max - min)
	return int64(generate_random_unsigned_integer_number_with_range(0, n)) + min
}

func generate_random_unsigned_integer_number_with_range(min, max uint64) uint64 {
	n := max - min
	x := seededRand.Uint64()

	if n < math.MaxInt64 {
		return uint64(rand.Int63n(int64(max-min))) + min
	}
	for x > n {
		x = seededRand.Uint64()
	}
	return x + min
}

func generate_random_number_with_type(number interface{}) interface{} {
	minNumber, maxNumber := generate_edge_value_with_type(number)
	return generate_random_number_with_range(minNumber, maxNumber)
}

func edge_value_mutate_with_type(number interface{}) interface{} {
	minNumber, maxNumber := generate_edge_value_with_type(number)
	randChoice := Generate_random_bool()
	switch randChoice {
	case false:
		return minNumber
	case true:
		return maxNumber
	default:
		return 0
	}
}

func generate_random_enum(enumValueMap map[string]int32) int32 {
	var enumRandomPool []int32
	for _, value := range enumValueMap {
		enumRandomPool = append(enumRandomPool, value)
	}

	return enumRandomPool[seededRand.Intn(len(enumRandomPool))]
}

// func strictly_increasing_generate_random_number_with_type(number interface{}) interface{} {
// 	_, maxNumber := generate_edge_value_with_type(number)
// 	return generate_random_number_with_range(number, maxNumber)
// }

// func strictly_decreasing_generate_random_number_with_type(number interface{}) interface{} {
// 	minNumber, _ := generate_edge_value_with_type(number)
// 	return generate_random_number_with_range(minNumber, number)
// }

func small_change_fundation_number[T Numeric](number T) T {
	changeRange := seededRand.Intn(5) + 1
	randChoice := Generate_random_bool()

	minNumber, maxNumber := generate_edge_value_with_type(number)

	if randChoice {
		if number > maxNumber.(T)-T(changeRange) {
			return maxNumber.(T)
		}
		return number + T(changeRange)
	}
	if number < minNumber.(T)+T(changeRange) {
		return minNumber.(T)
	}
	return number - T(changeRange)
}

func mutate_foundation_number_type[T Numeric](number T) T {
	randChoice := seededRand.Intn(100)
	if randChoice < 10 {
		return edge_value_mutate_with_type(number).(T)
	} else if randChoice < 90 {
		return T(small_change_fundation_number(number))
	} else {
		return generate_random_number_with_type(number).(T)
	}
}

func mutate_other_number_type(number interface{}) interface{} {
	switch number.(type) {
	case tbftpb.TBFTMsgType:
		return tbftpb.TBFTMsgType(generate_random_enum(tbftpb.TBFTMsgType_value))
	case tbftpb.VoteType:
		return tbftpb.VoteType(generate_random_enum(tbftpb.VoteType_value))
	case common.TxType:
		return common.TxType(generate_random_enum(common.TxType_value))
	default:
		// fmt.Printf("unknow number data type %s\n", reflect.TypeOf(number).String())
		return number
	}
}

func handle_number_mutate(number interface{}) (interface{}, error) {
	reflect.TypeOf(number)
	switch number := number.(type) {
	case uint:
		return mutate_foundation_number_type(number), nil
	case uint8:
		return mutate_foundation_number_type(number), nil
	case uint16:
		return mutate_foundation_number_type(number), nil
	case uint32:
		return mutate_foundation_number_type(number), nil
	case uint64:
		return mutate_foundation_number_type(number), nil
	case int:
		return mutate_foundation_number_type(number), nil
	case int8:
		return mutate_foundation_number_type(number), nil
	case int16:
		return mutate_foundation_number_type(number), nil
	case int32:
		return mutate_foundation_number_type(number), nil
	case int64:
		return mutate_foundation_number_type(number), nil
	case float32:
		return mutate_foundation_number_type(number), nil
	case float64:
		return mutate_foundation_number_type(number), nil
	default:
		return mutate_other_number_type(number), nil
	}
}

func handle_bytes_mutate(bytes []byte, name string) []byte {
	switch name {
	case "MemberInfo":
		return bytes
	default:
		return random_mutate_bytes(bytes)
	}
}

/* random mutation for bytes */
func random_mutate_bytes(bytes []byte) []byte {
	length := len(bytes)

	if length <= 0 {
		return bytes
	}

	mutate_times := generate_random_signed_integer_number_with_range(0, MAX_MUTATE_ITER+1)

	res := bytes

	for i := 0; i < mutate_times; i++ {
		pos := seededRand.Intn(length)
		res[pos] = byte(generate_random_unsigned_integer_number_with_range(0, math.MaxUint8))
	}
	return res
}

func handle_string_mutate(str string, name string) string {
	switch name {
	case "Voter":
		return str
	default:
		return random_mutate_string(str)
	}
}

/* random mutation for string */
func random_mutate_string(str string) string {
	length := len(str)

	if length <= 0 {
		return str
	}

	mutate_times := generate_random_signed_integer_number_with_range(0, MAX_MUTATE_ITER+1)

	res := str
	for i := 0; i < mutate_times; i++ {
		pos := seededRand.Intn(length)
		res = res[:pos] + string(charset[seededRand.Intn(len(charset))]) + res[pos+1:]
	}
	return res
}

func handleSlice(v reflect.Value, name string) ([]interface{}, error) {
	// if seededRand.Intn(100) < 80 {
	// 	return v.Interface().([]interface{}), nil
	// }

	if v.Type().Kind() != reflect.Slice {
		return nil, errors.New("incorrect config type: config type should be slice")
	}

	if v.Len() <= 0 || !v.IsValid() || v.IsNil() {
		return nil, nil
	}

	res := make([]interface{}, v.Len())

	switch reflect.ValueOf(v.Index(0).Interface()).Kind() {
	case reflect.Uint8:
		value := make([]byte, v.Len())
		for i := 0; i < v.Len(); i++ {
			value[i] = uint8(v.Index(i).Interface().(uint8))
		}
		value = handle_bytes_mutate(value, name)
		for i := 0; i < len(value); i++ {
			res[i] = value[i]
		}
	case reflect.Map:
		for i := 0; i < v.Len(); i++ {
			value, err := MutateMap(v.Index(i).Interface().(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			res[i] = value
		}
	case reflect.String:
		for i := 0; i < v.Len(); i++ {
			res[i] = handle_string_mutate(v.Index(i).Interface().(string), name)
		}
	// case reflect.Interface:
	// 	for i := 0; i < v.Len(); i++ {
	// 		switch v.Index(i).Interface().(type) {
	// 		case map[string]interface{}:
	// 			value, err := MutateMap(v.Index(i).Interface().(map[string]interface{}))
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			res[i] = value
	// 		case byte:
	// 			value := make([]byte, v.Len())
	// 			for i := 0; i < v.Len(); i++ {
	// 				value[i] = v.Index(i).Interface().(uint8)
	// 			}
	// 			value = handle_bytes_mutate(value, name)
	// 			for i := 0; i < len(value); i++ {
	// 				res[i] = value[i]
	// 			}
	// 		default:
	// 			errMsg := fmt.Sprintf("unknow data type %s\n", v.Index(0).Type().Kind().String())
	// 			return nil, errors.New(errMsg)
	// 		}
	// 	}
	default:
		errMsg := fmt.Sprintf("unknow data type %s\n", v.Index(0).Type().Kind().String())
		return nil, errors.New(errMsg)
	}

	return res, nil
}

func MutateMap(input map[string]interface{}) (map[string]interface{}, error) {
	var err error

	err = nil
	for k, value := range input {
		ident := reflect.ValueOf(value).Kind()
		switch ident {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			input[k], err = handle_number_mutate(value)
		case reflect.Bool:
			input[k] = Generate_random_bool()
		case reflect.String:
			input[k] = handle_string_mutate(value.(string), k)
		case reflect.Slice:
			input[k], err = handleSlice(reflect.ValueOf(value), k)
			if err != nil {
				fmt.Println(err)
			}
		case reflect.Map:
			input[k], err = MutateMap(value.(map[string]interface{}))
		default:
			errMsg := fmt.Sprintf("mutateMap fail, unknow value type, type is %s, value is %v\n",
				ident, value)
			return nil, errors.New(errMsg)
		}
	}
	return input, err
}
