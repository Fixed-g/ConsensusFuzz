package tbft

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"time"
)

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
		fmt.Println("It's not a number type!")
		return false, false
	}
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

func strictly_increasing_generate_random_number_with_type(number interface{}) interface{} {
	_, maxNumber := generate_edge_value_with_type(number)
	return generate_random_number_with_range(number, maxNumber)
}

func strictly_decreasing_generate_random_number_with_type(number interface{}) interface{} {
	minNumber, _ := generate_edge_value_with_type(number)
	return generate_random_number_with_range(minNumber, number)
}

func generate_random_number_with_type(number interface{}) interface{} {
	minNumber, maxNumber := generate_edge_value_with_type(number)
	if _, ok := minNumber.(int64); ok {
		return generate_random_signed_long_integer_number_with_range(minNumber.(int64), maxNumber.(int64))
	}
	return generate_random_number_with_range(minNumber, maxNumber)
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
		fmt.Println("It's not a number type!")
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

/* random mutation for bytes */
func random_mutate_bytes(bytes []byte) []byte {
	length := len(bytes)

	if length <= 0 {
		return bytes
	}

	mutate_times := generate_random_signed_integer_number_with_range(1, MAX_MUTATE_ITER+1)

	res := bytes

	for i := 0; i < mutate_times; i++ {
		pos := seededRand.Intn(length)
		res[pos] = byte(generate_random_unsigned_integer_number_with_range(0, math.MaxUint8))
	}
	return res
}

/* random mutation for string */
func random_mutate_string(str string) string {
	length := len(str)

	if length <= 0 {
		return str
	}

	mutate_times := generate_random_signed_integer_number_with_range(1, MAX_MUTATE_ITER+1)

	res := str
	for i := 0; i < mutate_times; i++ {
		pos := seededRand.Intn(length)
		res = res[:pos] + string(charset[seededRand.Intn(len(charset))]) + res[pos+1:]
	}
	return res
}

func handleSlice(v reflect.Value) ([]interface{}, error) {
	if v.Type().Kind() != reflect.Slice {
		return nil, errors.New("incorrect config type: config type should be slice")
	}

	if v.Len() <= 0 || !v.IsValid() || v.IsNil() {
		return nil, nil
	}

	res := make([]interface{}, v.Len())

	switch v.Index(0).Type().Kind() {
	case reflect.Uint8:
		value := make([]byte, v.Len())
		for i := 0; i < v.Len(); i++ {
			value[i] = uint8(v.Index(i).Uint())
		}
		value = random_mutate_bytes(value)
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
	case reflect.Interface:
		for i := 0; i < v.Len(); i++ {
			switch v.Index(i).Interface().(type) {
			case map[string]interface{}:
				value, err := MutateMap(v.Index(i).Interface().(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				res[i] = value
			case byte:
				value := make([]byte, v.Len())
				for i := 0; i < v.Len(); i++ {
					value[i] = v.Index(i).Interface().(uint8)
				}
				value = random_mutate_bytes(value)
				for i := 0; i < len(value); i++ {
					res[i] = value[i]
				}
			default:
				errMsg := fmt.Sprintf("unknow data type %s\n", v.Index(0).Type().Kind().String())
				return nil, errors.New(errMsg)
			}
		}
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
		case reflect.Int:
			input[k] = generate_random_number_with_type(value)
		case reflect.Int8:
			input[k] = generate_random_number_with_type(value)
		case reflect.Int16:
			input[k] = generate_random_number_with_type(value)
		case reflect.Int32:
			input[k] = generate_random_number_with_type(value)
		case reflect.Int64:
			input[k] = generate_random_number_with_type(value)
		case reflect.Uint:
			input[k] = generate_random_number_with_type(value)
		case reflect.Uint8:
			input[k] = generate_random_number_with_type(value)
		case reflect.Uint16:
			input[k] = generate_random_number_with_type(value)
		case reflect.Uint32:
			input[k] = generate_random_number_with_type(value)
		case reflect.Uint64:
			input[k] = generate_random_number_with_type(value)
		case reflect.Float32:
			input[k] = generate_random_number_with_type(value)
		case reflect.Float64:
			input[k] = generate_random_number_with_type(value)
		case reflect.Bool:
			input[k] = Generate_random_bool()
		case reflect.String:
			input[k] = random_mutate_string(value.(string))
		case reflect.Slice:
			input[k], err = handleSlice(reflect.ValueOf(value))
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
