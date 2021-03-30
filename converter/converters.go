package converter

import (
	"fmt"
)

func StrToStr(pointer *string) string {
	if pointer == nil {
		return ""
	}
	return *pointer
}

func IntToStr(pointer *int) string {
	if pointer == nil {
		return ""
	}
	return fmt.Sprintf("%d", *pointer)
}
