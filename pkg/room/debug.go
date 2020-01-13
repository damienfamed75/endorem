package room

import (
	"fmt"
	"strings"
)

func debugMap(data string) {
	data = strings.Replace(data, "X", "\033[1;34mX\033[0m", -1)
	data = strings.Replace(data, "Y", "\033[1;31mY\033[0m", -1)
	data = strings.Replace(data, ":", "\033[1;32m:\033[0m", -1)
	data = strings.Replace(data, "^", "\033[1;33m^\033[0m", -1)
	data = strings.Replace(data, "-", "\033[1;33m-\033[0m", -1)
	data = strings.Replace(data, "}", "\033[1;33m}\033[0m", -1)
	data = strings.Replace(data, "~", "\033[1;34m~\033[0m", -1)
	data = strings.Replace(data, "=", "\033[1;34m=\033[0m", -1)

	fmt.Println(data)
}
