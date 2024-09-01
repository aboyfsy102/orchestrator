package fleet

import (
	"encoding/json"
	"testing"
)

var fleetJsonSamples = []string{
	`{"requirements":{"total_target_capacity":1},"extra_tags":[{"key":"env","value":"dev"}]}`,
	`{"requirements":{"total_target_capacity":1, "vcpu_count":{"min":1, "max":2}},"extra_tags":[{"key":"env","value":"dev"}]}`,
	`{"requirements":{"total_target_capacity":1, "memory_gib":{"min":1, "max":2}},"extra_tags":[{"key":"env","value":"dev"}]}`,
	`{"requirements":{"total_target_capacity":1, "memory_gib_per_vcpu":{"min":1, "max":2}},"extra_tags":[{"key":"env","value":"dev"}]}`,
	`{"requirements":{"total_target_capacity":1, "excluded_instance_types":["t3.small"]},"extra_tags":[{"key":"env","value":"dev"}]}`,
	`{"requirements":{"total_target_capacity":1, "allowed_instance_types":["t3.small"]},"extra_tags":[{"key":"env","value":"dev"}]}`,
	`{"requirements":{"total_target_capacity":1, "excluded_instance_types":["t3.small"], "allowed_instance_types":["t3.medium"]},"extra_tags":[{"key":"env","value":"dev"}]}`,
}

func TestFleetRequest_Validate(t *testing.T) {
	for _, sample := range fleetJsonSamples {
		var fleet FleetRequest
		err := json.Unmarshal([]byte(sample), &fleet)
		if err != nil {
			t.Fatalf("Failed to unmarshal fleet request: %v", err)
		}
		if fleet.Requirements == nil {
			t.Fatalf("Failed to unmarshal fleet request: %v", err)
		}
	}
}
