package orderedmap

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

// Setup suite
type OrderedMapSuite struct {
	suite.Suite
}

func TestOrderedMapSuiteTest(t *testing.T) {
	suite.Run(t, new(OrderedMapSuite))
}

// NewOrderedMap tests
func (s *OrderedMapSuite) TestNewOrderedMapSinglVal() {
	p := [][]string{
		{"foo"},
	}

	_, err := NewOrderedMap(p)
	s.Error(err)
}

//
// OrderedMap.GetAll(key string) tests
//
func (s *OrderedMapSuite) TestEmptyOrderedMap() {
	q := OrderedMap{}
	s.Equal([]string{}, q.GetAll("foo"))
}

func (s *OrderedMapSuite) TestEmptyDuplicateParam() {
	p := [][]string{
		{"foo", "1"},
		{"bar", "2"},
		{"foo", "3"},
	}
	q, _ := NewOrderedMap(p)
	s.Equal([]string{"1", "3"}, q.GetAll("foo"))
}

//
// OrderedMap.Add tests
//
func (s *OrderedMapSuite) TestAddDoesntExist() {
	p := OrderedMap{}
	p.Add("key", "val")
	s.Equal("val", p.Get("key"))
}

func (s *OrderedMapSuite) TestAddExists() {
	key := "key"
	p := [][]string{
		{key, "val1"},
	}

	q, _ := NewOrderedMap(p)
	q.Add("key", "val2")
	s.Equal("val1", q.Get(key))
	s.Equal([]string{"val1", "val2"}, q.GetAll(key))
}

//
// OrderedMap.Remove tests
//

func (s *OrderedMapSuite) TestRemoveDoesntExist() {
	key := "key"
	p := [][]string{
		{"foo", "bar"},
	}
	q, _ := NewOrderedMap(p)
	q.Remove(key, "val")
	s.Equal(1, q.Len())
}

func (s *OrderedMapSuite) TestRemoveExists() {
	key := "key"
	val1 := "val1"
	p := [][]string{
		{key, "valfirst"},
		{"key2", "val1"},
		{key, val1},
	}

	q, _ := NewOrderedMap(p)
	q.Remove(key, val1)

	s.Equal(2, q.Len())
	s.Equal("valfirst", q.Get(key))
	s.Equal("val1", q.Get("key2"))
}

//
// OrderedMap.Del tests
//
func (s *OrderedMapSuite) TestDelDoesntExist() {
	p := OrderedMap{}
	p.Del("foo")
	s.Equal(0, p.Len())
}
func (s *OrderedMapSuite) TestDelSingleExists() {
	key := "key"
	p := [][]string{
		{key, "val1"},
	}
	q, err := NewOrderedMap(p)
	s.NoError(err)
	q.Del(key)
	s.Equal(0, q.Len())
}

func (s *OrderedMapSuite) TestDelDoubleExists() {
	key := "key"
	p := [][]string{
		{key, "val1"},
		{"something", "val2"},
		{key, "2"},
	}
	m, err := NewOrderedMap(p)
	s.NoError(err)
	m.Del(key)
	s.Equal(1, m.Len())
	s.Equal("val2", m.Get("something"))
}

//
// OrderedMap.Keys tests
//
func (s *OrderedMapSuite) TestKeysEmpty() {
	m, _ := NewOrderedMap([][]string{})
	s.Equal(0, len(m.Keys()))
}

func (s *OrderedMapSuite) TestKeysSingle() {
	m, _ := NewOrderedMap([][]string{
		{"foo", "bar"},
	})
	s.Equal(1, len(m.Keys()))
}

func (s *OrderedMapSuite) TestKeysDuplicateNext() {
	m, _ := NewOrderedMap([][]string{
		{"foo", "bar"},
		{"meh", "bar"},
		{"foo", "bar"},
	})
	s.Equal(2, len(m.Keys()))
}

//
// OrderedMap.Map tests
//
func (s *OrderedMapSuite) TestMapEmpty() {
	keyParams, err := NewOrderedMap([][]string{})
	s.NoError(err)
	s.Equal(map[string][]string{}, keyParams.Map())
}

func (s *OrderedMapSuite) TestMapSingle() {
	keyParams, err := NewOrderedMap([][]string{{"key", "val"}})
	s.NoError(err)
	expected := map[string][]string{
		"key": []string{"val"},
	}
	s.Equal(expected, keyParams.Map())
}

func (s *OrderedMapSuite) TestMapDuplicate() {
	key := "mykey"
	p := [][]string{
		{key, "9"},
		{"blah", "blash"},
		{key, "4"},
	}
	m, _ := NewOrderedMap(p)
	mapped := m.Map()
	s.Equal(mapped["blah"], []string{"blash"})
	// Duplicate keys should be combined into one result dict.
	// The first value declared by `key` ("9" in this case) should come first
	s.Equal(mapped[key], []string{"9", "4"})
}

//
// OrderedMap.Get tests
//
func (s *OrderedMapSuite) TestGetNoKey() {
	p := [][]string{}
	q, _ := NewOrderedMap(p)
	s.Equal("", q.Get("something"))
}
func (s *OrderedMapSuite) TestGetSingleKey() {
	p := [][]string{
		{"key", "val"},
		{"something", "val"},
	}
	q, _ := NewOrderedMap(p)
	s.Equal("val", q.Get("key"))
}
func (s *OrderedMapSuite) TestGetMultipleVals() {
	key := "key"
	p := [][]string{
		{key, "val1"},
		{"something", "val"},
		{key, "val2"},
	}
	q, _ := NewOrderedMap(p)
	s.Equal("val1", q.Get(key))
}

//
// OrderedMap.Set tests
//
func (s *OrderedMapSuite) TestSetNoKey() {
	p := [][]string{}
	q, _ := NewOrderedMap(p)
	q.Set("foo", "bar")
	s.Equal("bar", q.Get("foo"))
}

func (s *OrderedMapSuite) TestSetSingleKey() {
	key := "key"
	p := [][]string{
		{key, "val"},
	}
	q, _ := NewOrderedMap(p)

	newVal := "val2"
	q.Set(key, newVal)
	s.Equal(1, q.Len())
	s.Equal(newVal, q.Get(key))
}

func (s *OrderedMapSuite) TestSetMultipleKeys() {
	key := "key"
	val := "val1"
	p := [][]string{
		{key, val},
		{key, "somethingelse"},
	}
	q, _ := NewOrderedMap(p)
	newVal := "val2"
	q.Set(key, newVal)
	s.Equal(1, q.Len())
	s.Equal(newVal, q.Get(key))
}

//
// OrderedMap construction and String methods
//
func (s *OrderedMapSuite) TestNewOrderedMapNil() {
	r, err := NewOrderedMap(nil)
	s.Error(err)
	s.Nil(r)

}

// OrderedMap.String() tests
func (s *OrderedMapSuite) TestNewOrderedMapEmpty() {
	r, err := NewOrderedMap([][]string{})
	s.NoError(err)
	s.Equal("", r.String())
}

func (s *OrderedMapSuite) TestNewOrderedMapSingle() {
	t := [][]string{
		{"test", "param"},
	}
	r, err := NewOrderedMap(t)
	s.NoError(err)
	s.Equal("test=param", r.String())
}

func (s *OrderedMapSuite) TestNewOrderedMapMulti() {
	t := [][]string{
		{"test", "1"},
		{"test", "2"},
	}
	r, err := NewOrderedMap(t)
	s.NoError(err)
	s.Equal("test=1&test=2", r.String())
}

func (s *OrderedMapSuite) TestGetStrTwoParamsOneVal() {
	t := [][]string{
		{"test", "1"},
		{"foobar", "2"},
		{"test", "4"},
	}
	r, err := NewOrderedMap(t)
	s.NoError(err)
	s.Equal("test=1&foobar=2&test=4", r.String())
}

//
