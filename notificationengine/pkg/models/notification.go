package models

// Notification represents a generic notification structure.
type Notification struct {
	ID                       string                   `json:"id"`
	Type                     string                   `json:"type"`
	Severity                 string                   `json:"severity"`
	ReceivedAt               string                   `json:"receivedAt"`
	GeneratedAt              string                   `json:"generatedAt"`
	Source                   string                   `json:"source"`
	Name                     string                   `json:"name"`
	AlertID                  string                   `json:"alertId"`
	Description              string                   `json:"descripton"`
	DeviceData               DeviceData               `json:"deviceData"`
	CustomerData             CustomerData             `json:"customerData"`
	AckStateMetadata         AckStateMatadata         `json:"ackStateMetadata,omitempty"`
	BroadcastedStateMetadata BroadcastedStateMetadata `json:"broadcastedStateMetadata,omitempty"`
	ResolvedStateMetadata    ResolvedStateMetadata    `json:"resolvedStateMetadata,omitempty"`
	SnoozedStateMetadata     SnoozedStateMetadata     `json:"snoozedStateMetadata,omitempty"`
	Parent                   string                   `json:"parent,omitempty"`
	Priority                 string                   `json:"priority,omitempty"`
	Status                   string                   `json:"status"`
	Category                 string                   `json:"category"`
	RebroadcastedMetadata    RebroadcastedMetadata    `json:"rebroadcastedMetadata,omitempty"`
	StateChangedAt           string                   `json:"stateChangedAt,omitempty"`
	ReassignmentMetadata     ReassignmentMetadata     `json:"reassignmentMetadata,omitempty"`
}

// DeviceData represents device-specific information.
type DeviceData struct {
	DeviceName   string `json:"deviceName"`
	DeviceID     string `json:"deviceId"`
	PropertyID   string `json:"propertyId"`
	PropertyName string `json:"propertyName"`
	Location     string `json:"location,omitempty"`
}

// CustomerData represents customer-specific information.
type CustomerData struct {
	CustomerName string `json:"customerName,omitempty"`
	CustomerID   string `json:"customerId"`
	SiteName     string `json:"siteName,omitempty"`
	SiteID       string `json:"siteId"`
}

type AckStateMatadata struct {
	AcknowledgedBy        string `json:"acknowledgedBy"`
	AcknowledgerID        string `json:"acknowledgerId"`
	AcknowledgedAt        string `json:"acknowledgedAt"`
	InitialAcknowledgedAt string `json:"initialAckwnowledgedAt"`
}

type BroadcastedStateMetadata struct {
	BroadcastedBy string `json:"broadcastedBy"`
	BroadcastedAt string `json:"broadcastedAt"`
}

type ResolvedStateMetadata struct {
	ResolvedBy     string `json:"resolvedBy"`
	ResolverID     string `json:"resolverId"`
	ResolvedAt     string `json:"resolvedAt"`
	ResolutionNote string `json:"resolutionNote"`
	ResolutionType string `json:"resolutionType"`
}

type SnoozedStateMetadata struct {
	SnoozedBy   string `json:"snoozedBy"`
	SnoozerID   string `json:"snoozerId"`
	SnoozedAt   string `json:"snoozedAt"`
	SnoozedFrom string `json:"snoozedFrom"`
	SnoozedTill string `json:"snoozedTill"`
	SnoozeNote  string `json:"snoozeNote"`
}

type RebroadcastedMetadata struct {
	RebroadcastedBy    string `json:"rebroadcastedBy"`
	RebroadcasterID    string `json:"rebroadcasterId"`
	RebroadcastedAt    string `json:"rebroadcastedAt"`
	RebroadcastingNote string `json:"rebroadcastingNote"`
}

type ReassignmentMetadata struct {
	ReassignedBy   string `json:"reassignedBy"`
	ReassignerID   string `json:"reassignerId"`
	ReassignedAt   string `json:"reassignedAt"`
	ReassignNote   string `json:"reassignNote"`
	ReassignedTo   string `json:"reassignedTo"`
	ReassignedToID string `json:"reassignedUserId"`
}
