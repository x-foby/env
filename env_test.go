package env

import (
	"errors"
	"testing"
)

func Test_checkDest(t *testing.T) {
	cases := []struct {
		dest interface{}
		err  error
	}{
		{
			dest: "",
			err:  ErrDestMustBeStructPointer,
		},
		{
			dest: interface{}(0),
			err:  ErrDestMustBeStructPointer,
		},
		{
			dest: struct{}{},
			err:  ErrDestMustBeStructPointer,
		},
		{
			dest: &struct{}{},
		},
	}

	for _, c := range cases {
		if err := checkDest(c.dest); !errors.Is(err, c.err) {
			t.Errorf("checkDest(%v) = %v; want %v", c.dest, err, c.err)
		}
	}
}

type envMap map[string]string

func (e envMap) lookupEnv(key string) (string, bool) {
	v, ok := e[key]

	return v, ok
}

func Test_lookup(t *testing.T) {
	mocks := envMap{
		"ONE": "one",
	}

	type testStruct struct {
		One string `env:"ONE"`
	}

	type testStructWithDefault struct {
		String  string  `default:"one"`
		Bool    bool    `default:"true"`
		Int     int     `default:"-1"`
		Int8    int8    `default:"-1"`
		Int16   int16   `default:"-1"`
		Int32   int32   `default:"-1"`
		Int64   int64   `default:"-1"`
		Uint    uint    `default:"1"`
		Uint8   uint8   `default:"1"`
		Uint16  uint16  `default:"1"`
		Uint32  uint32  `default:"1"`
		Uint64  uint64  `default:"1"`
		Float32 float32 `default:"1.23"`
		Float64 float64 `default:"1.23"`
	}

	type testStructWithRequired struct {
		One string `env:"UNEXISTS" required:"true"`
	}
	type testStructWithRequiredAndDefault struct {
		One string `env:"UNEXISTS" required:"true" default:"one"`
	}

	lookupEnv = mocks.lookupEnv

	cases := []struct {
		prefix string
		dest   interface{}
		err    error
	}{
		{
			dest: &testStruct{},
		},
		{
			dest: &testStructWithDefault{},
		},
		{
			dest: &testStructWithRequired{},
			err:  ErrEnvRequired,
		},
	}

	for _, c := range cases {
		err := lookup("", c.dest)

		if !errors.Is(err, c.err) {
			t.Errorf("lookup(%s, %v) = %v; want %v", c.prefix, c.dest, err, c.err)
		}

		switch typed := c.dest.(type) {
		case *testStruct:
			if typed.One != "one" {
				t.Errorf("lookup(%s, %v) failed; want dest.One = %s", c.prefix, c.dest, "one")
			}
		case *testStructWithDefault:
			if typed.String != "one" {
				t.Errorf("lookup(%s, %v) failed; want dest.One = %s", c.prefix, c.dest, "one")
			}

			if typed.Bool != true {
				t.Errorf("lookup(%s, %v) failed; want dest.Bool = %v", c.prefix, c.dest, true)
			}

			if typed.Int != -1 {
				t.Errorf("lookup(%s, %v) failed; want dest.Int = %v", c.prefix, c.dest, -1)
			}

			if typed.Int8 != -1 {
				t.Errorf("lookup(%s, %v) failed; want dest.Int8 = %v", c.prefix, c.dest, -1)
			}

			if typed.Int16 != -1 {
				t.Errorf("lookup(%s, %v) failed; want dest.Int16 = %v", c.prefix, c.dest, -1)
			}

			if typed.Int32 != -1 {
				t.Errorf("lookup(%s, %v) failed; want dest.Int32 = %v", c.prefix, c.dest, -1)
			}

			if typed.Int64 != -1 {
				t.Errorf("lookup(%s, %v) failed; want dest.Int64 = %v", c.prefix, c.dest, -1)
			}

			if typed.Uint != 1 {
				t.Errorf("lookup(%s, %v) failed; want dest.Uint = %v", c.prefix, c.dest, 1)
			}

			if typed.Uint8 != 1 {
				t.Errorf("lookup(%s, %v) failed; want dest.Uint8 = %v", c.prefix, c.dest, 1)
			}

			if typed.Uint16 != 1 {
				t.Errorf("lookup(%s, %v) failed; want dest.Uint16 = %v", c.prefix, c.dest, 1)
			}

			if typed.Uint32 != 1 {
				t.Errorf("lookup(%s, %v) failed; want dest.Uint32 = %v", c.prefix, c.dest, 1)
			}

			if typed.Uint64 != 1 {
				t.Errorf("lookup(%s, %v) failed; want dest.Uint64 = %v", c.prefix, c.dest, 1)
			}

			if typed.Float32 != 1.23 {
				t.Errorf("lookup(%s, %v) failed; want dest.Float32 = %v", c.prefix, c.dest, 1.23)
			}

			if typed.Float64 != 1.23 {
				t.Errorf("lookup(%s, %v) failed; want dest.Float64 = %v", c.prefix, c.dest, 1.23)
			}
		case *testStructWithRequiredAndDefault:
			if typed.One != "one" {
				t.Errorf("lookup(%s, %v) failed; want dest.One = %s", c.prefix, c.dest, "one")
			}
		}

	}
}

func Test_getEnvName(t *testing.T) {
	cases := []struct {
		prefix, fieldEnvName, fieldName string
		fieldEnvNameDefined             bool
		want                            string
	}{
		{
			prefix:              "",
			fieldEnvName:        "FIELDENVNAME",
			fieldName:           "FieldName",
			fieldEnvNameDefined: true,
			want:                "FIELDENVNAME",
		},
		{
			prefix:              "PREFIX",
			fieldEnvName:        "FIELDENVNAME",
			fieldName:           "FieldName",
			fieldEnvNameDefined: true,
			want:                "PREFIX_FIELDENVNAME",
		},
		{
			prefix:              "",
			fieldEnvName:        "",
			fieldName:           "FieldName",
			fieldEnvNameDefined: false,
			want:                "FieldName",
		},
		{
			prefix:              "PREFIX",
			fieldEnvName:        "",
			fieldName:           "FieldName",
			fieldEnvNameDefined: false,
			want:                "PREFIX_FieldName",
		},
	}

	for _, c := range cases {
		if got := getEnvName(c.prefix, c.fieldEnvName, c.fieldName, c.fieldEnvNameDefined); got != c.want {
			t.Errorf("getEnvName(%s, %s, %s, %v) = %s; want %s", c.prefix, c.fieldEnvName, c.fieldName, c.fieldEnvNameDefined, got, c.want)
		}
	}

}
