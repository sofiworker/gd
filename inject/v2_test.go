package inject

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProvide(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	t.Run("basic operations", func(t *testing.T) {
		container := New()

		// Test basic provides
		assert.NoError(container.ProvideWithName("testData", "testName"))
		assert.NoError(container.Provide("testData", WithGroup("group")))

		// Test error cases
		assert.ErrorIs(container.Provide(123), BasicTypeShouldWithNameError)
		assert.ErrorIs(container.Provide("data", WithName("name"), WithGroup("group")), NameOrGroupOnlyOneError)
		assert.ErrorIs(container.Provide(nil), NilError)
		var ptr *int
		assert.ErrorIs(container.Provide(ptr), NilPtrError)
	})

	t.Run("data types", func(t *testing.T) {
		container := New()
		cases := map[string]struct {
			value interface{}
			name  string
		}{
			"string":  {"test", "str"},
			"int":     {42, "int"},
			"float":   {3.14, "float"},
			"bool":    {true, "bool"},
			"slice":   {[]string{"a", "b"}, "slice"},
			"array":   {[2]int{1, 2}, "array"},
			"map":     {map[string]int{"a": 1}, "map"},
			"channel": {make(chan string), "chan"},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				assert.NoError(container.ProvideWithName(tc.value, tc.name))
			})
		}
	})

	t.Run("struct operations", func(t *testing.T) {
		container := New()

		// Test injectable struct
		type TestStruct struct {
			Field1 int
			Field2 string
		}
		assert.NoError(container.Provide(&TestStruct{Field1: 1, Field2: "test"}))

		// Test interface
		runner := &TestRunner{}
		require.NoError(container.Provide(runner, WithAs(new(Runner))))
		require.NoError(container.Invoke(func(r Runner) {
			assert.Equal("running", r.Run())
		}))

		// Test nested structs
		type Inner struct{ Value string }
		type Outer struct {
			Inner    *Inner
			OtherVal int
		}
		require.NoError(container.Provide(&Inner{Value: "test"}))
		require.NoError(container.Provide(&Outer{OtherVal: 42}))
		require.NoError(container.Invoke(func(o *Outer) {
			assert.NotNil(o.Inner)
			assert.Equal("test", o.Inner.Value)
		}))

		// Test tagged struct
		type TaggedStruct struct {
			Required string  `gd:"required"`
			Optional *int    `gd:"optional"`
			Ignored  float64 `gd:"-"`
			Untagged bool
		}
		val := 42
		tagged := &TaggedStruct{
			Required: "test",
			Optional: &val,
			Ignored:  3.14,
			Untagged: true,
		}
		assert.NoError(container.Provide(tagged))
	})
}

type Runner interface {
	Run() string
}

type TestRunner struct{}

func (tr *TestRunner) Run() string {
	return "running"
}

func TestInvoke(t *testing.T) {
	container := New()
	err := container.Provide(nil)
	require.NoError(t, err)

	err = container.Invoke(&TestRunner{})

}
