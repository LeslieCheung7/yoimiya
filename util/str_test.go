package util

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestStrToFloat64(t *testing.T) {
	testCases := []struct {
		name   string
		val    string
		expVal float64
		expErr bool
	}{
		{
			name:   "valid value",
			val:    "777.7777777",
			expVal: 777.7777777,
			expErr: false,
		},
		{
			name:   "out of range",
			val:    "1.7e+309",
			expVal: math.Inf(1),
			expErr: true,
		},
		{
			name:   "invalid value",
			val:    "invalid",
			expVal: 0,
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := StrToFloat64(tc.val)
			if tc.expErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, res, tc.expVal)
		})
	}
}

func TestFloat64ToStr(t *testing.T) {
	testCases := []struct {
		name   string
		val    float64
		expVal string
	}{
		{
			name:   "min value",
			val:    0,
			expVal: "0",
		},
		{
			name: "max value",
			val:  1.7e+308,
			expVal: "170000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
				"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
				"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
				"0000000000000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := Float64ToStr(tc.val)
			assert.Equal(t, res, tc.expVal)
		})
	}
}

func TestStrToInt64(t *testing.T) {
	testCases := []struct {
		name   string
		val    string
		expVal int64
		expErr bool
	}{
		{
			name:   "valid",
			val:    "7777777",
			expVal: 7777777,
			expErr: false,
		},
		{
			name:   "out of range",
			val:    "9243372036854775809",
			expVal: math.MaxInt64,
			expErr: true,
		},
		{
			name:   "invalid",
			val:    "invalid",
			expVal: 0,
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := StrToInt64(tc.val)
			if tc.expErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, res, tc.expVal)
		})
	}
}

func TestStrToUint64(t *testing.T) {
	testCases := []struct {
		name   string
		val    string
		expVal uint64
		expErr bool
	}{
		{
			name:   "valid",
			val:    "7777777",
			expVal: 7777777,
			expErr: false,
		},
		{
			name:   "out of range: exceeds max limit",
			val:    "18446744073709551617",
			expVal: math.MaxUint64,
			expErr: true,
		},
		{
			name:   "out of range: negative value",
			val:    "-1",
			expVal: 0,
			expErr: true,
		},
		{
			name:   "invalid",
			val:    "invalid",
			expVal: 0,
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := StrToUint64(tc.val)
			if tc.expErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, res, tc.expVal)
		})
	}
}
