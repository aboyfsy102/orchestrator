package fleet

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type FleetRequest struct {
	Requirements *struct {
		VCpuCount *struct {
			Min int64 `json:"min,omitempty"`
			Max int64 `json:"max,omitempty"`
		} `json:"vcpu_count"`
		MemoryGib *struct {
			Min int64 `json:"min,omitempty"`
			Max int64 `json:"max,omitempty"`
		} `json:"memory_gib"`
		MemoryGibPerVCpu *struct {
			Min int64 `json:"min,omitempty"`
			Max int64 `json:"max,omitempty"`
		} `json:"memory_gib_per_vcpu"`
		ExcludedInstanceTypes []string `json:"excluded_instance_types,omitempty"`
		AllowedInstanceTypes  []string `json:"allowed_instance_types,omitempty"`
		TotalTargetCapacity   int64    `json:"total_target_capacity"`
	} `json:"requirements"`
	ExtraTags []*struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"extra_tags,omitempty"`
	Alarms []*struct {
		Threshold         int64  `json:"threshold"`
		Metric            string `json:"metric"`
		DatapointsToAlarm int64  `json:"datapoints_to_alarm"`
		EvaluationPeriods int64  `json:"evaluation_periods"`
		Period            int64  `json:"period"`
	} `json:"alarms,omitempty"`
	Userdata string `json:"userdata,omitempty"`
}

type FleetOrder struct {
	FleetRequest
	OrderId   string           `json:"order_id"`
	Status    FleetOrderStatus `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	DeletedAt time.Time        `json:"deleted_at"`
}

type FleetOrderStatus string

const (
	FleetOrderStatusPending      FleetOrderStatus = "pending"
	FleetOrderStatusProvisioning FleetOrderStatus = "provisioning"
	FleetOrderStatusProvisioned  FleetOrderStatus = "provisioned"
	FleetOrderStatusCancelled    FleetOrderStatus = "cancelled"
	FleetOrderStatusCompleted    FleetOrderStatus = "completed"
)

// NewFleetOrder creates a new fleet order from a string in json format
func NewFleetOrder(jsonString string) (*FleetOrder, error) {
	var order FleetOrder
	err := json.Unmarshal([]byte(jsonString), &order)
	if err != nil {
		return nil, err
	}
	order.OrderId = uuid.New().String()
	order.CreatedAt = time.Now()
	return &order, nil
}

// fleetorder save
func (order *FleetOrder) Save(status FleetOrderStatus) error {
	order.Status = status
	return nil
}

func (order *FleetOrder) Delete() error {
	order.DeletedAt = time.Now()
	return nil
}
