package ids

import (
	"strconv"
	"strings"
)

const separator = "."

func makeId(parentId string, number int) string {
	if len(parentId) > 0 {
		parentId += separator
	}
	return parentId + strconv.Itoa(number)
}

func splitId(id string) (parent string, number string) {
	index := strings.LastIndex(id, separator)
	if index < 0 {
		return "", id
	}
	return id[0:index], id[index+1:]
}
