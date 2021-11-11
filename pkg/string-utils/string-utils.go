package string_utils

func ArrayToString(array []string, separator string) string {
	var str string
	for _, v := range array {
		str += v + separator
	}
	return str
}
