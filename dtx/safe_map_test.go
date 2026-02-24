package dtx_test

import (
	"reflect"
	"testing"

	"github.com/mikeschinkel/go-dt/dtx"
)

type Key string
type Value struct {
	Value string
}

func TestNewSafeMap(t *testing.T) {
	type args struct {
		cap int
	}
	type testCase[K comparable, V any] struct {
		name string
		args args
		want *dtx.SafeMap[K, V]
	}
	tests := []testCase[Key, Value]{
		{
			name: "create empty map with cap 0",
			args: args{cap: 0},
			want: dtx.NewSafeMap[Key, Value](0),
		},
		{
			name: "create map with cap 10",
			args: args{cap: 10},
			want: dtx.NewSafeMap[Key, Value](10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dtx.NewSafeMap[Key, Value](tt.args.cap)
			if got == nil {
				t.Errorf("NewSafeMap() returned nil")
			}
			if got.Len() != 0 {
				t.Errorf("NewSafeMap() initial length = %v, want 0", got.Len())
			}
		})
	}
}

func TestSafeMap_Delete(t *testing.T) {
	type args[K comparable] struct {
		k K
	}
	type testCase[K comparable, V any] struct {
		name      string
		setupFunc func() *dtx.SafeMap[K, V]
		args      args[K]
		wantLen   int
	}
	tests := []testCase[Key, Value]{
		{
			name: "delete existing key",
			setupFunc: func() *dtx.SafeMap[Key, Value] {
				sm := dtx.NewSafeMap[Key, Value](10)
				sm.Set("key1", Value{Value: "value1"})
				sm.Set("key2", Value{Value: "value2"})
				return sm
			},
			args:    args[Key]{k: "key1"},
			wantLen: 1,
		},
		{
			name: "delete non-existing key",
			setupFunc: func() *dtx.SafeMap[Key, Value] {
				sm := dtx.NewSafeMap[Key, Value](10)
				sm.Set("key1", Value{Value: "value1"})
				return sm
			},
			args:    args[Key]{k: "key2"},
			wantLen: 1,
		},
		{
			name: "delete from empty map",
			setupFunc: func() *dtx.SafeMap[Key, Value] {
				return dtx.NewSafeMap[Key, Value](10)
			},
			args:    args[Key]{k: "key1"},
			wantLen: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := tt.setupFunc()
			sm.Delete(tt.args.k)
			if got := sm.Len(); got != tt.wantLen {
				t.Errorf("After Delete(), Len() = %v, want %v", got, tt.wantLen)
			}
		})
	}
}

func TestSafeMap_Get(t *testing.T) {
	type args[K comparable] struct {
		k K
	}
	type testCase[K comparable, V any] struct {
		name      string
		setupFunc func() *dtx.SafeMap[K, V]
		args      args[K]
		wantV     V
		wantOk    bool
	}
	tests := []testCase[Key, Value]{
		{
			name: "get existing key",
			setupFunc: func() *dtx.SafeMap[Key, Value] {
				sm := dtx.NewSafeMap[Key, Value](10)
				sm.Set("key1", Value{Value: "value1"})
				return sm
			},
			args:   args[Key]{k: "key1"},
			wantV:  Value{Value: "value1"},
			wantOk: true,
		},
		{
			name: "get non-existing key",
			setupFunc: func() *dtx.SafeMap[Key, Value] {
				sm := dtx.NewSafeMap[Key, Value](10)
				sm.Set("key1", Value{Value: "value1"})
				return sm
			},
			args:   args[Key]{k: "key2"},
			wantV:  Value{},
			wantOk: false,
		},
		{
			name: "get from empty map",
			setupFunc: func() *dtx.SafeMap[Key, Value] {
				return dtx.NewSafeMap[Key, Value](10)
			},
			args:   args[Key]{k: "key1"},
			wantV:  Value{},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := tt.setupFunc()
			gotV, gotOk := sm.Get(tt.args.k)
			if !reflect.DeepEqual(gotV, tt.wantV) {
				t.Errorf("Get() gotV = %v, want %v", gotV, tt.wantV)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Get() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestSafeMap_Iter(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name      string
		setupFunc func() *dtx.SafeMap[K, V]
		wantPairs map[K]V
	}
	tests := []testCase[Key, Value]{
		{
			name: "iterate over populated map",
			setupFunc: func() *dtx.SafeMap[Key, Value] {
				sm := dtx.NewSafeMap[Key, Value](10)
				sm.Set("key1", Value{Value: "value1"})
				sm.Set("key2", Value{Value: "value2"})
				sm.Set("key3", Value{Value: "value3"})
				return sm
			},
			wantPairs: map[Key]Value{
				"key1": {Value: "value1"},
				"key2": {Value: "value2"},
				"key3": {Value: "value3"},
			},
		},
		{
			name: "iterate over empty map",
			setupFunc: func() *dtx.SafeMap[Key, Value] {
				return dtx.NewSafeMap[Key, Value](10)
			},
			wantPairs: map[Key]Value{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := tt.setupFunc()
			got := make(map[Key]Value)
			for k, v := range sm.Iter() {
				got[k] = v
			}
			if !reflect.DeepEqual(got, tt.wantPairs) {
				t.Errorf("Iter() produced %v, want %v", got, tt.wantPairs)
			}
		})
	}
}

func TestSafeMap_Len(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name      string
		setupFunc func() *dtx.SafeMap[K, V]
		want      int
	}
	tests := []testCase[Key, Value]{
		{
			name: "empty map",
			setupFunc: func() *dtx.SafeMap[Key, Value] {
				return dtx.NewSafeMap[Key, Value](10)
			},
			want: 0,
		},
		{
			name: "map with one item",
			setupFunc: func() *dtx.SafeMap[Key, Value] {
				sm := dtx.NewSafeMap[Key, Value](10)
				sm.Set("key1", Value{Value: "value1"})
				return sm
			},
			want: 1,
		},
		{
			name: "map with multiple items",
			setupFunc: func() *dtx.SafeMap[Key, Value] {
				sm := dtx.NewSafeMap[Key, Value](10)
				sm.Set("key1", Value{Value: "value1"})
				sm.Set("key2", Value{Value: "value2"})
				sm.Set("key3", Value{Value: "value3"})
				return sm
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := tt.setupFunc()
			if got := sm.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeMap_Set(t *testing.T) {
	type args[K comparable, V any] struct {
		k K
		v V
	}
	type testCase[K comparable, V any] struct {
		name      string
		setupFunc func() *dtx.SafeMap[K, V]
		args      args[K, V]
		wantLen   int
		wantValue V
	}
	tests := []testCase[Key, Value]{
		{
			name: "set new key",
			setupFunc: func() *dtx.SafeMap[Key, Value] {
				return dtx.NewSafeMap[Key, Value](10)
			},
			args:      args[Key, Value]{k: "key1", v: Value{Value: "value1"}},
			wantLen:   1,
			wantValue: Value{Value: "value1"},
		},
		{
			name: "overwrite existing key",
			setupFunc: func() *dtx.SafeMap[Key, Value] {
				sm := dtx.NewSafeMap[Key, Value](10)
				sm.Set("key1", Value{Value: "old_value"})
				return sm
			},
			args:      args[Key, Value]{k: "key1", v: Value{Value: "new_value"}},
			wantLen:   1,
			wantValue: Value{Value: "new_value"},
		},
		{
			name: "set multiple keys",
			setupFunc: func() *dtx.SafeMap[Key, Value] {
				sm := dtx.NewSafeMap[Key, Value](10)
				sm.Set("key1", Value{Value: "value1"})
				return sm
			},
			args:      args[Key, Value]{k: "key2", v: Value{Value: "value2"}},
			wantLen:   2,
			wantValue: Value{Value: "value2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := tt.setupFunc()
			sm.Set(tt.args.k, tt.args.v)
			if got := sm.Len(); got != tt.wantLen {
				t.Errorf("After Set(), Len() = %v, want %v", got, tt.wantLen)
			}
			gotV, gotOk := sm.Get(tt.args.k)
			if !gotOk {
				t.Errorf("After Set(), Get() returned ok = false")
			}
			if !reflect.DeepEqual(gotV, tt.wantValue) {
				t.Errorf("After Set(), Get() = %v, want %v", gotV, tt.wantValue)
			}
		})
	}
}

func TestSafeMap_Values(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name       string
		setupFunc  func() *dtx.SafeMap[K, V]
		wantValues []V
	}
	tests := []testCase[Key, Value]{
		{
			name: "iterate over values in populated map",
			setupFunc: func() *dtx.SafeMap[Key, Value] {
				sm := dtx.NewSafeMap[Key, Value](10)
				sm.Set("key1", Value{Value: "value1"})
				sm.Set("key2", Value{Value: "value2"})
				sm.Set("key3", Value{Value: "value3"})
				return sm
			},
			wantValues: []Value{
				{Value: "value1"},
				{Value: "value2"},
				{Value: "value3"},
			},
		},
		{
			name: "iterate over empty map",
			setupFunc: func() *dtx.SafeMap[Key, Value] {
				return dtx.NewSafeMap[Key, Value](10)
			},
			wantValues: []Value{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := tt.setupFunc()
			var got []Value
			for v := range sm.Values() {
				got = append(got, v)
			}
			// Convert slices to maps for comparison (order doesn't matter)
			gotMap := make(map[string]bool)
			for _, v := range got {
				gotMap[v.Value] = true
			}
			wantMap := make(map[string]bool)
			for _, v := range tt.wantValues {
				wantMap[v.Value] = true
			}
			if !reflect.DeepEqual(gotMap, wantMap) {
				t.Errorf("Values() produced %v, want %v", got, tt.wantValues)
			}
			if len(got) != len(tt.wantValues) {
				t.Errorf("Values() count = %v, want %v", len(got), len(tt.wantValues))
			}
		})
	}
}
