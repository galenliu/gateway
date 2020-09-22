package accessory

import (
	"github.com/iancoleman/strcase"
)

func CamelCase(s string) string {
	return strcase.ToSnake(s)
}
