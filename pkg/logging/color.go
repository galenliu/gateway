package logging

import "fmt"

func Json(json string) string {
	return fmt.Sprintf("%s"+json, "\u001B[36m")
}
