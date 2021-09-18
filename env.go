package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

var (
	ErrDestMustBeStructPointer = errors.New("destination must be pointer on struct")
	ErrEnvRequired             = errors.New("env required")
	ErrUnsupportedFieldType    = errors.New("unsupported field type")
)

var lookupEnv = os.LookupEnv

func Parse(prefix string, dest interface{}) error {
	if err := checkDest(dest); err != nil {
		return err
	}

	return lookup(prefix, dest)
}

func checkDest(dest interface{}) error {
	if dest == nil {
		return ErrDestMustBeStructPointer
	}

	t := reflect.TypeOf(dest)

	if t.Kind() != reflect.Ptr {
		return ErrDestMustBeStructPointer
	}

	if t.Elem().Kind() != reflect.Struct {
		return ErrDestMustBeStructPointer
	}

	return nil
}

func lookup(prefix string, dest interface{}) error {
	el := reflect.TypeOf(dest).Elem()

	for i := 0; i < el.NumField(); i++ {
		fieldDef := el.Field(i)

		fieldEnvName, fieldEnvNameDefined := fieldDef.Tag.Lookup("env")
		_, fieldEnvRequired := fieldDef.Tag.Lookup("required")
		fieldEnvDefault := fieldDef.Tag.Get("default")

		envName := getEnvName(prefix, fieldEnvName, fieldDef.Name, fieldEnvNameDefined)

		field := reflect.ValueOf(dest).Elem().Field(i)

		kind := field.Kind()

		switch kind {
		case reflect.Interface:
			// todo
		case reflect.Ptr:
			// todo
		case reflect.Struct:
			// todo
		}

		envValue, envExists := lookupEnv(envName)

		if !envExists && fieldEnvDefault != "" {
			envValue = fieldEnvDefault
		} else if !envExists && fieldEnvRequired {
			return fmt.Errorf("%w <%s>", ErrEnvRequired, envName)
		}

		switch kind {
		case reflect.Bool:
			v, err := strconv.ParseBool(envValue)
			if err != nil {
				return err
			}

			reflect.ValueOf(dest).Elem().Field(i).SetBool(v)
		case reflect.Int:
			v, err := strconv.ParseInt(envValue, 10, strconv.IntSize)
			if err != nil {
				return err
			}

			reflect.ValueOf(dest).Elem().Field(i).SetInt(v)
		case reflect.Int8:
			v, err := strconv.ParseInt(envValue, 10, 8)
			if err != nil {
				return err
			}

			reflect.ValueOf(dest).Elem().Field(i).SetInt(v)
		case reflect.Int16:
			v, err := strconv.ParseInt(envValue, 10, 16)
			if err != nil {
				return err
			}

			reflect.ValueOf(dest).Elem().Field(i).SetInt(v)
		case reflect.Int32:
			v, err := strconv.ParseInt(envValue, 10, 32)
			if err != nil {
				return err
			}

			reflect.ValueOf(dest).Elem().Field(i).SetInt(v)
		case reflect.Int64:
			v, err := strconv.ParseInt(envValue, 10, 64)
			if err != nil {
				return err
			}

			reflect.ValueOf(dest).Elem().Field(i).SetInt(v)
		case reflect.Uint:
			v, err := strconv.ParseUint(envValue, 10, strconv.IntSize)
			if err != nil {
				return err
			}

			reflect.ValueOf(dest).Elem().Field(i).SetUint(v)
		case reflect.Uint8:
			v, err := strconv.ParseUint(envValue, 10, 8)
			if err != nil {
				return err
			}

			reflect.ValueOf(dest).Elem().Field(i).SetUint(v)
		case reflect.Uint16:
			v, err := strconv.ParseUint(envValue, 10, 16)
			if err != nil {
				return err
			}

			reflect.ValueOf(dest).Elem().Field(i).SetUint(v)
		case reflect.Uint32:
			v, err := strconv.ParseUint(envValue, 10, 32)
			if err != nil {
				return err
			}

			reflect.ValueOf(dest).Elem().Field(i).SetUint(v)
		case reflect.Uint64:
			v, err := strconv.ParseUint(envValue, 10, 64)
			if err != nil {
				return err
			}

			reflect.ValueOf(dest).Elem().Field(i).SetUint(v)
		case reflect.Float32:
			v, err := strconv.ParseFloat(envValue, 32)
			if err != nil {
				return err
			}

			reflect.ValueOf(dest).Elem().Field(i).SetFloat(v)
		case reflect.Float64:
			v, err := strconv.ParseFloat(envValue, 64)
			if err != nil {
				return err
			}

			reflect.ValueOf(dest).Elem().Field(i).SetFloat(v)
		case reflect.Array:
			// todo
		case reflect.Map:
			// todo
		case reflect.Slice:
			// todo
		case reflect.String:
			reflect.ValueOf(dest).Elem().Field(i).SetString(envValue)

		default:
			return fmt.Errorf("%w <%s>", ErrUnsupportedFieldType, kind)
		}
	}

	return nil
}

func getEnvName(prefix, fieldEnvName, fieldName string, fieldEnvNameDefined bool) string {
	localName := fieldEnvName
	if !fieldEnvNameDefined {
		localName = fieldName
	}

	if prefix == "" {
		return localName
	}

	return prefix + "_" + localName
}
