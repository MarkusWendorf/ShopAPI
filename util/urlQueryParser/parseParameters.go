package urlQueryParser

import (
	"net/url"
	"strconv"
)

type UrlQueryParser struct {
	values url.Values
}

func New(values url.Values) UrlQueryParser {
	return UrlQueryParser{values: values}
}

func (u *UrlQueryParser) GetInt(key string, defaultVal int) int {

	val := u.values.Get(key)
	if val == "" {
		return defaultVal
	}

	toInt, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}

	return toInt
}

func (u *UrlQueryParser) GetString(key string, defaultVal string) string {

	val := u.values.Get(key)
	if val == "" {
		return defaultVal
	}

	return val
}
