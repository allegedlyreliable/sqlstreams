package cli

import "testing"

func TestMetricFiltersRejectConflictingAttributeValues(t *testing.T) {
	// setup
	attributes := []string{"stream=orders", "stream=payments"}

	// test
	_, err := parseAttributePairs(attributes)

	// verify
	if err == nil || exitCode(err) != 2 {
		t.Errorf("parseAttributePairs(%v) = %v, want usage error with exit 2", attributes, err)
	}
}

func TestMetricFiltersRequireEveryAttributeIncludingEmptyValues(t *testing.T) {
	// setup
	filter, err := parseAttributePairs([]string{"stream=orders", "stream=orders", "region="})
	if err != nil {
		t.Fatal(err)
	}

	// test
	missing := attributesMatch(map[string]string{"stream": "orders"}, filter)
	matching := attributesMatch(map[string]string{"stream": "orders", "region": ""}, filter)
	different := attributesMatch(map[string]string{"stream": "payments", "region": ""}, filter)

	// verify
	if missing {
		t.Error("attributesMatch(stream=orders, region absent) = true, want false")
	}
	if !matching {
		t.Error("attributesMatch(stream=orders, region empty) = false, want true")
	}
	if different {
		t.Error("attributesMatch(stream=payments, region empty) = true, want false")
	}
}
