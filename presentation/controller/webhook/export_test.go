package webhook

import (
	"github.com/duck8823/duci/application/service/executor"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/url"
	"reflect"
)

type Handler = handler

func (h *Handler) SetExecutor(executor executor.Executor) (reset func()) {
	tmp := h.executor
	h.executor = executor
	return func() {
		h.executor = tmp
	}
}

func URLMust(url *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	return url
}

func CmpOptsAllowFields(typ interface{}, names ...string) cmp.Option {
	return cmpopts.IgnoreFields(typ, func() []string {
		var ignoreFields []string

		t := reflect.TypeOf(typ)
		for i := 0; i < t.NumField(); i++ {
			name := t.Field(i).Name
			if !contains(names, name) {
				ignoreFields = append(ignoreFields, name)
			}
		}
		return ignoreFields
	}()...)
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
