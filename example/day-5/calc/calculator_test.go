package calc_test

import (
	"testing"
	"wongnok/example/day-5/calc"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	result := calc.Add(2, 3)
	expected := 5
	if result != expected {
		t.Errorf(
			"Actual %d; expected %d",
			result,
			expected,
		)
	}
}

func TestAddWithTestify(t *testing.T) {
	result := calc.Add(2, 3)
	expected := 5

	assert.Equal(t, expected, result,
		"Add(2, 3) should equal 5",
	)
}
