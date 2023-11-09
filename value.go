package httpc

import (
	"strconv"
)

type Value func(values)

type values interface {
	Add(key, val string)
}

func String(key, value string) Value {
	return func(u values) {
		u.Add(key, value)
	}
}

func Stringp(key string, value *string) Value {
	return func(u values) {
		if value != nil {
			u.Add(key, *value)
		}
	}
}

func Int64(key string, value int64) Value {
	return func(u values) {
		u.Add(key, strconv.FormatInt(value, 10))
	}
}

func Int(key string, value int) Value {
	return func(u values) {
		u.Add(key, strconv.Itoa(value))
	}
}

func Intp(key string, value *int) Value {
	return func(u values) {
		if value != nil {
			u.Add(key, strconv.Itoa(*value))
		}
	}
}

func Values(val map[string][]string) Value {
	return func(u values) {
		for key, vals := range val {
			for _, value := range vals {
				u.Add(key, value)
			}
		}
	}
}
