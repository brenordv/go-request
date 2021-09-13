package parsers

import "fmt"

func InterfaceToString(it interface{}) string {
	str := fmt.Sprintf("%v", it)
	return str
}
