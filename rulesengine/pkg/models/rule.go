package models

// Condition represents a condition in the rule
type Condition struct {
	ID             string `json:"id"`
	Position       string `json:"position"`
	Type           string `json:"type"`
	DeviceId       string `json:"deviceId,omitempty"`
	PropertyId     string `json:"propertyId,omitempty"`
	OperatorId     string `json:"operatorId,omitempty"`
	OperatorSymbol string `json:"operatorSymbol,omitempty"`
	Value          string `json:"value,omitempty"`
}

// Rule represents a rule
type Rule struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Severity   string      `json:"severity"`
	Status     string      `json:"status"`
	CreatedAt  string      `json:"createdAt"`
	UpdatedAt  string      `json:"updatedAt"`
	Conditions []Condition `json:"condition"`
}
