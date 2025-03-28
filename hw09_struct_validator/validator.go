package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrValidator = errors.New("wrong validator")
	ErrStruct    = errors.New("it is not a structure")
	ErrTag       = errors.New("wrong tag")
	ErrIn        = errors.New("error validating the IN tag")
	ErrMin       = errors.New("error validating the MIN tag")
	ErrMax       = errors.New("error validating the MAX tag")
	ErrLen       = errors.New("error validating the LEN tag")
	ErrRegexp    = errors.New("error validating the REGEXP tag")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func checkInInt(allowedNums string, val int) error {
	allowedNumsSplit := strings.Split(allowedNums, ",")
	for _, num := range allowedNumsSplit {
		numInt, err := strconv.Atoi(num)
		if err != nil {
			return err
		}

		if numInt == val {
			return nil
		}
	}

	return ErrIn
}

func checkMinInt(allowedNum string, val int) error {
	allowedNumInt, err := strconv.Atoi(allowedNum)
	if err != nil {
		return err
	}

	if val < allowedNumInt {
		return ErrMin
	}

	return nil
}

func checkMaxInt(allowedNum string, val int) error {
	allowedNumInt, err := strconv.Atoi(allowedNum)
	if err != nil {
		return err
	}

	if val > allowedNumInt {
		return ErrMax
	}

	return nil
}

func checkInStr(allowedStrs string, val string) error {
	allowedNumsSplit := strings.Split(allowedStrs, ",")
	for _, allowedStr := range allowedNumsSplit {
		if allowedStr == val {
			return nil
		}
	}

	return ErrIn
}

func checkLenStr(allowedLen string, val string) error {
	allowedLenInt, err := strconv.Atoi(allowedLen)
	if err != nil {
		return err
	}

	if allowedLenInt != len(val) {
		return ErrLen
	}

	return nil
}

func checkRegexpStr(allowedReg string, val string) error {
	re, err := regexp.Compile(allowedReg)
	if err != nil {
		return err
	}

	if !re.MatchString(val) {
		return ErrRegexp
	}

	return nil
}

func ValidateInt(tags []string, value int) error {
	functions := map[string]func(string, int) error{
		"in":  checkInInt,
		"min": checkMinInt,
		"max": checkMaxInt,
	}

	for _, tag := range tags {
		validator := strings.Split(tag, ":")
		if len(validator) < 2 {
			return ErrValidator
		}

		if f, ok := functions[validator[0]]; ok {
			if err := f(validator[1], value); err != nil {
				return err
			}
		} else {
			return ErrValidator
		}
	}

	return nil
}

func ValidateStr(tags []string, value string) error {
	functions := map[string]func(string, string) error{
		"in":     checkInStr,
		"len":    checkLenStr,
		"regexp": checkRegexpStr,
	}
	for _, tag := range tags {
		validator := strings.Split(tag, ":")
		if len(validator) < 2 {
			return ErrValidator
		}

		if f, ok := functions[validator[0]]; ok {
			if err := f(validator[1], value); err != nil {
				return err
			}
		} else {
			return ErrValidator
		}
	}

	return nil
}

func ValidateSlice(field reflect.StructField, tags []string, valueOf reflect.Value) ValidationErrors {
	var valErrs ValidationErrors
	elemType := field.Type.Elem()

	switch elemType.Kind() {
	case reflect.Int:
		if slice, ok := valueOf.FieldByName(field.Name).Interface().([]int); ok {
			for _, elem := range slice {
				if err := ValidateInt(tags, elem); err != nil {
					valErrs = append(valErrs, ValidationError{Field: field.Name, Err: err})
				}
			}
		}
	case reflect.String:
		if slice, ok := valueOf.FieldByName(field.Name).Interface().([]string); ok {
			for _, elem := range slice {
				if err := ValidateStr(tags, elem); err != nil {
					valErrs = append(valErrs, ValidationError{Field: field.Name, Err: err})
				}
			}
		}
	case reflect.Invalid, reflect.Bool, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Array, reflect.Chan,
		reflect.Func, reflect.Interface, reflect.Map, reflect.Slice, reflect.Struct, reflect.UnsafePointer, reflect.Pointer:
		// unsupported types
	}

	return valErrs
}

func Validate(v interface{}) error {
	typeOf := reflect.TypeOf(v)
	if typeOf.Kind() != reflect.Struct {
		return ErrStruct
	}

	valueOf := reflect.ValueOf(v)

	fields := reflect.VisibleFields(typeOf)
	valErrs := make(ValidationErrors, 0)
	for _, field := range fields {
		if !field.IsExported() {
			continue
		}

		if val, ok := field.Tag.Lookup("validate"); ok {
			if len(val) == 0 {
				return ErrTag
			}

			tags := strings.Split(val, "|")
			fieldName := field.Name

			switch field.Type.Kind() {
			case reflect.Int:
				v := valueOf.FieldByName(fieldName).Interface().(int)
				if err := ValidateInt(tags, v); err != nil {
					valErrs = append(valErrs, ValidationError{Field: fieldName, Err: err})
				}
			case reflect.String:
				v := valueOf.FieldByName(fieldName).String()
				if err := ValidateStr(tags, v); err != nil {
					valErrs = append(valErrs, ValidationError{Field: field.Name, Err: err})
				}
			case reflect.Slice:
				errs := ValidateSlice(field, tags, valueOf)
				valErrs = append(valErrs, errs...)
			case reflect.Invalid, reflect.Bool, reflect.Int8, reflect.Int16, reflect.Int32,
				reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
				reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64,
				reflect.Complex128, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
				reflect.Map, reflect.Struct, reflect.UnsafePointer, reflect.Pointer:
				// unsupported types
			}
		}
	}

	if len(valErrs) > 0 {
		return valErrs
	}

	return nil
}
