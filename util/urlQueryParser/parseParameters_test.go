package urlQueryParser

import (
	"math"
	"net/url"
	"strconv"
	"testing"
)

func TestUrlQueryParser_GetInt(t *testing.T) {

	count8 := url.Values{}
	count8.Add("count", "8")

	page8 := url.Values{}
	page8.Add("page", "8")

	pageUnset := url.Values{}

	sizeIntegerOverflow := url.Values{}
	sizeIntegerOverflow.Add("size", strconv.Itoa(math.MaxInt64)+"0")

	tests := []struct {
		values       url.Values
		key          string
		defaultValue int
		expected     int
	}{
		{count8, "count", 1, 8},
		{page8, "page", 1, 8},
		{pageUnset, "page", 1, 1},
		{sizeIntegerOverflow, "size", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {

			parser := New(tt.values)
			result := parser.GetInt(tt.key, tt.defaultValue)

			if result != tt.expected {
				t.Errorf("got %d, expected %d", result, tt.expected)
			}
		})
	}
}

func TestUrlQueryParser_GetString(t *testing.T) {

	categorySSD := url.Values{}
	categorySSD.Add("category", "SSD")

	nameAsus := url.Values{}
	nameAsus.Add("name", "Asus")

	categoryUnset := url.Values{}

	tests := []struct {
		values       url.Values
		key          string
		defaultValue string
		expected     string
	}{
		{categorySSD, "category", "", "SSD"},
		{nameAsus, "name", "", "Asus"},
		{categoryUnset, "category", "Prozessoren", "Prozessoren"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {

			parser := New(tt.values)
			result := parser.GetString(tt.key, tt.defaultValue)

			if result != tt.expected {
				t.Errorf("got %s, expected %s", result, tt.expected)
			}
		})
	}
}
