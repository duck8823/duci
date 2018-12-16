package container_test

import (
	"fmt"
	"github.com/duck8823/duci/domain/internal/container"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

type testInterface interface {
	Hoge()
}

type testImpl struct {
	Name string
}

func (*testImpl) Hoge() {}

type testInterfaceNothing interface {
	Fuga()
}

type hoge string

func TestSubmit(t *testing.T) {
	// given
	ins := &container.SingletonContainer{}
	defer ins.SetValues(map[string]interface{}{})()
	defer container.SetInstance(ins)()

	// and
	want := map[string]interface{}{
		"string": "test",
	}

	// when
	err := container.Submit("test")

	// then
	if err != nil {
		t.Errorf("must be nil, but got %+v", err)
	}

	// and
	if !cmp.Equal(ins.GetValues(), want) {
		t.Errorf("must be equal, but %+v", cmp.Diff(ins.GetValues(), want))
	}

	// when
	err = container.Submit("twice")

	// then
	if err == nil {
		t.Error("must not be nil")
	}

	// and
	if !cmp.Equal(ins.GetValues(), want) {
		t.Errorf("must be equal, but %+v", cmp.Diff(ins.GetValues(), want))
	}

	// when
	str := "ptr tri"
	err = container.Submit(&str)

	// then
	if err == nil {
		t.Error("must not be nil")
	}

	// and
	if !cmp.Equal(ins.GetValues(), want) {
		t.Errorf("must be equal, but %+v", cmp.Diff(ins.GetValues(), want))
	}

}

func TestOverride(t *testing.T) {
	// given
	ins := &container.SingletonContainer{}
	defer ins.SetValues(map[string]interface{}{
		"string": "hoge",
	})()
	defer container.SetInstance(ins)()

	// and
	want := map[string]interface{}{
		"string": "test",
	}

	// when
	container.Override("test")

	// then
	if !cmp.Equal(ins.GetValues(), want) {
		t.Errorf("must be equal, but %+v", cmp.Diff(ins.GetValues(), want))
	}

}

func TestGet(t *testing.T) {
	// given
	ins := &container.SingletonContainer{}
	defer ins.SetValues(map[string]interface{}{
		"string":   "value",
		"int":      1234,
		"float64":  12.34,
		"container_test.testImpl": &testImpl{Name: "hoge"},
	})()
	defer container.SetInstance(ins)()

	// where
	for _, tt := range []struct {
		in   interface{}
		want interface{}
		err  bool
	}{
		{
			in: new(string),
			want: "value",
		},
		{
			in: new(int),
			want: 1234,
		},
		{
			in: new(float64),
			want: 12.34,
		},
		{
			in: new(hoge),
			want: hoge(""),
			err: true,
		},
		{
			in: new(testImpl),
			want: testImpl{Name: "hoge"},
		},
		{
			in: new(testInterface),
			want: &testImpl{Name: "hoge"},
		},
		{
			in: new(testInterfaceNothing),
			err: true,
		},
	} {
		t.Run(fmt.Sprintf("type=%s", reflect.TypeOf(tt.in).String()), func(t *testing.T) {
			// when
			err := container.Get(tt.in)
			got := Value(tt.in)

			// then
			if tt.err && err == nil {
				t.Error("must not be nil")
			}
			if !tt.err && err != nil {
				t.Errorf("must be nil, but got %+v", err)
			}

			// and
			if !cmp.Equal(got, tt.want) {
				t.Errorf("must be equal, but %+v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func Value(v interface{}) interface{} {
	return reflect.Indirect(reflect.ValueOf(v)).Interface()
}