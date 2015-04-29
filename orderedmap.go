package orderedmap

import (
	"fmt"
	"strings"
	"sync"
)

// OrderedMappable replicates the default net/url.Values interface
type OrderedMappable interface {
	Add(key, value string)
	Del(key string)
	Get(key string) string
	Set(key, value string)
	Len() int
	String() string
	GetAll(key string) []string
}

// NewOrderedMap takes a set of params, validates them and return the params to use
// Seriously, use this rather than constructing OrderedMap directly. It'll keep your sanity
func NewOrderedMap(keypairs [][]string) (*OrderedMap, error) {
	if keypairs == nil {
		return nil, fmt.Errorf("Params cannot be nil")
	}
	if !hasValidKeyPairs(keypairs) {
		return nil, fmt.Errorf("Invalid query keypairs `%s`", keypairs)
	}
	queryParams := OrderedMap{keypairs: keypairs}
	return &queryParams, nil
}

// OrderedMap handles url query string paramaters respecting order
// We replicate the interface of the default net/url.Values
type OrderedMap struct {
	keypairs [][]string
	sync.Mutex
}

// Add add the key-value to the params list
// This will not delete any existing pairs matching the key param
func (m *OrderedMap) Add(key, value string) {
	m.Lock()
	defer m.Unlock()
	m.add(key, value)
}

func (m *OrderedMap) add(key, value string) {
	m.keypairs = append(m.keypairs, []string{key, value})
}

// Remove removes the key value pair from the params list
func (m *OrderedMap) Remove(key, val string) {
	m.Lock()
	defer m.Unlock()
	m.remove(key, val)
}
func (m *OrderedMap) remove(key, val string) {
	newKeyPairs := [][]string{}
	// Add all keys that aren't the passed in key
	for _, pair := range m.keypairs {
		if !(stringEquals(pair[0], key) && stringEquals(pair[1], val)) {
			newKeyPairs = append(newKeyPairs, pair)
		}
	}
	// Replace the key pairs with our new, filtered keypair slice
	m.keypairs = newKeyPairs
}

// Del removes all values with the passed in key
func (m *OrderedMap) Del(key string) {
	m.Lock()
	defer m.Unlock()
	m.del(key)
}
func (m *OrderedMap) del(key string) {
	newKeyPairs := [][]string{}
	// Add all keys that aren't the passed in key
	for _, pair := range m.keypairs {
		if !stringEquals(pair[0], key) {
			newKeyPairs = append(newKeyPairs, pair)
		}
	}
	// Replace the existing keypairs obj with our new obj
	m.keypairs = newKeyPairs
}

// Keys returns all keys that exist in our OrderedMap
func (m *OrderedMap) Keys() []string {
	keys := []string{}
	for k, _ := range m.Map() {
		keys = append(keys, k)
	}
	return keys
}

// Map converts the our ordered list to a map
func (m *OrderedMap) Map() map[string][]string {
	mapOut := map[string][]string{}
	for _, keyPair := range m.keypairs {
		key := keyPair[0]
		val := keyPair[1]
		mapOut[key] = append(mapOut[key], val)
	}
	return mapOut
}

// Get gets the first value associated with key. If empty, returns an empty string
func (m *OrderedMap) Get(key string) string {
	all := m.GetAll(key)
	if len(all) > 0 {
		return all[0]
	}
	return ""
}

// Get a list of values based on the key
func (m *OrderedMap) GetAll(key string) []string {
	if key == "" {
		return nil
	}

	keyVals := []string{}
	for _, param := range m.keypairs {
		// Only build up our list of params for OUR key
		paramKey := param[0]
		if key == paramKey {
			keyVals = append(keyVals, param[1])
		}
	}
	return keyVals
}

// Set sets the key to val. All existing values are replaced
func (m *OrderedMap) Set(key, val string) {
	m.Lock()
	defer m.Unlock()
	m.set(key, val)
}
func (m *OrderedMap) set(key, val string) {
	m.del(key)
	m.add(key, val)
}

// Len get's the total number of entries within the map
func (m *OrderedMap) Len() int {
	return len(m.keypairs)
}

// String converts the query params into consumable query param strings
func (m *OrderedMap) String() string {
	// No paramaters are fine, just a waste to call this func
	numParams := m.Len()
	if numParams == 0 {
		return ""
	}

	// It's build time...
	queryStr := ""
	for i, queryParam := range m.keypairs {
		queryStr += queryParam[0] + "=" + queryParam[1]
		if numParams != 0 && i+1 != numParams {
			queryStr += "&"
		}
	}
	return queryStr
}

// Private helper methods

// isValidKeyPair verifies that we have a valid query param (a 2 pair of queryKey, queryVal)
func isValidKeyPair(param []string) bool {
	return len(param) == 2
}

//  hasValidKeyPairs ensures that all of the query parameters are correct
func hasValidKeyPairs(pairs [][]string) bool {
	for _, pair := range pairs {
		if !isValidKeyPair(pair) {
			return false
		}
	}
	return true
}

// stringEquals
func stringEquals(str1, str2 string) bool {
	return len(str1) == len(str2) && strings.Contains(str1, str2)
}
