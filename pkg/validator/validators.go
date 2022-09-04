package validator

import (
	"fmt"
	"reflect"
)

func IfNotNil(field any, fn func(fieldname string) error) ValidatorHandler {

	if reflect.ValueOf(field).IsNil() {
		return func(fieldname string) error { return nil }
	}
	return fn
}

func Required(field any) ValidatorHandler {
	return func(fieldname string) error {
		if reflect.ValueOf(field).IsNil() {
			return fmt.Errorf("Field (%s) cannot be empty", fieldname)
		}
		return nil
	}
}

func MinLength(field *string, min uint) ValidatorHandler {
	return func(fieldname string) error {
		if field == nil || uint(len(*field)) < min {
			return fmt.Errorf("Value of field (%s) cannot be less than %d", fieldname, min)
		}
		return nil
	}
}

func MaxLength(field *string, max uint) ValidatorHandler {
	return func(fieldname string) error {
		if field != nil && uint(len(*field)) > max {
			return fmt.Errorf("Value of field (%s) cannot be more than %d", fieldname, max)
		}
		return nil
	}
}

const (
	LETTERS_ONLY_TEMPLATE = "!#$%^&*()_+=-1234567890`~\\|[]{};:'\"/?.>,< "
)

func WithoutSymbols(field *string, symbols ...rune) ValidatorHandler {
	return func(fieldname string) error {
		if field == nil {
			return nil
		}

		if ExistSymbolsInString(symbols, *field) {
			return fmt.Errorf("Value of field (%s) cannot contain symbols: %s", fieldname, JoinRunes(symbols, ", "))
		}
		return nil
	}
}
