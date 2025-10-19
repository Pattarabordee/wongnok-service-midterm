package assertion_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestBasicAssertions(t *testing.T) {
	// Equality assertions
	assert.Equal(t, 123, 123, "numbers should be equal")
	assert.NotEqual(t, 123, 456, "numbers should not be equal")

	// Boolean assertions
	assert.True(t, 1 < 2, "1 should be less than 2")
	assert.False(t, 1 > 2, "1 should not be greater than 2")

	// Nil assertions
	var nilPointer *int
	assert.Nil(t, nilPointer, "pointer should be nil")

	nonNilPointer := new(int)
	assert.NotNil(t, nonNilPointer, "pointer should not be nil")
}

func TestWithError(t *testing.T) {
	// Working with errors
	result, err := strconv.ParseInt("", 10, 64)
	assert.Error(t, err)
	assert.Empty(t, result)

	result, err = strconv.ParseInt("20", 10, 64)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

type UserTestSuite struct {
	suite.Suite // Extend suite
	name        string
}

// This will run once before all tests in the suite
func (s *UserTestSuite) SetupSuite() {
	fmt.Println("Before all")
}

// This will run once after all tests in the suite
func (s *UserTestSuite) TearDownSuite() {
	fmt.Println("After all")
}

// This will run before each test
func (s *UserTestSuite) SetupTest() {
	fmt.Println("Before each")
	s.name = "Peter"
}

// This will run after each test
func (s *UserTestSuite) TearDownTest() {
	fmt.Println("After each")
	s.name = ""
}

func (s *UserTestSuite) TestGetUser() {
	s.Equal("Peter", s.name)
}

func (s *UserTestSuite) TestCreateUser() {
	s.name = "Parker"
	s.Equal("Peter", s.name)
}

func TestUser(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
