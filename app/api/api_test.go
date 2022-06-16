package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatAsDate(t *testing.T) {
	assert := assert.New(t)
	var unixTime int64 = 1637907727
	assert.Equal("26 November 2021", formatAsDate(unixTime))
}

func TestFormatAsPrice(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("0.00", formatAsPrice(0))
	assert.Equal("0.07", formatAsPrice(7))
	assert.Equal("0.42", formatAsPrice(42))
	assert.Equal("1.46", formatAsPrice(146))
	assert.Equal("1005.00", formatAsPrice(100500))
}

func TestFormatAsCurrency(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("$", formatAsCurrency("USD"))
	assert.Equal("Unknown", formatAsCurrency("Unknown"))
}
