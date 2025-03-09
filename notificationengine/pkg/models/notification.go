package models

type ElasticStruct struct {
	Action string `json:"action"`
	Notification Notification `json:"notification"`
}

// Notification represents a generic notification structure.
type Notification struct {
	ID                    string `json:"id"`
	Type                  string `json:"type"`
	Severity              string `json:"severity"`
	GeneratedAt           string `json:"generatedAt"`
	Source                string `json:"source"`
	Name                  string `json:"name"`
	AlertID               string `json:"alertId"`
	Description           string `json:"descripton"`
	DeviceData            DeviceData `json:"deviceData"`
	CustomerData          CustomerData `json:"customerData"`
	AcknowledgedBy        string `json:"acknwoledgedBy,omitempty"`
	AcknowledgerID        string `json:"acknowledgerId,omitempty"`
	ResolvedBy            string `json:"resolvedBy,omitempty"`
	ResolverID            string `json:"resolverId,omitempty"`
	SnoozedBy             string `json:"snoozedBy,omitempty"`
	SnoozerID             string `json:"snoozerId,omitempty"`
	Parent                string `json:"parent,omitempty"`
	ResolutionNote        string `json:"resolutionNote,omitempty"`
	SnoozeNote            string `json:"snoozeNote,omitempty"`
	SnoozeFrom            string `json:"snoozeFrom,omitempty"` 
	SnoozeTill            string `json:"snoozeTill,omitempty"` 
	ResolutionType        string `json:"resolutionType,omitempty"`
	Priority              string `json:"priority,omitempty"`
	Status                string `json:"status"`
	Category              string `json:"category"`
	AcknowledgedAt        string `json:"acknowledgedAt,omitempty"` 
	SnoozedAt             string `json:"snoozedAt,omitempty"` 
	ResolvedAt            string `json:"resolvedAt,omitempty"` 
	BroadcastedAt         string `json:"broadcastedAt,omitempty"`
	RebroadcastedAt       string `json:"rebroadcstedAt,omitempty"`
	RebroadcastedBy       string `json:"rebroadcastedBy,omitempty"`
	RebroadcasterID       string `json:"rebroadcasterId,omitempty"`
	StateChangeAt         string `json:"stateChangeAt,omitempty"` 
	ReassignedAt          string `json:"reaasignedAt,omitempty"`
	InitialAcknowledgedAt string `json:"initialAcknowledgedAt,omitempty"` 
}

// DeviceData represents device-specific information.
type DeviceData struct {
	DeviceName 		string `json:"deviceName"`
	DeviceID   		string `json:"deviceId"`
	PropertyID 		string `json:"propertyId"`
	PropertyName 	string `json:"propertyName"`
	Location   		string `json:"location,omitempty"`
}

// CustomerData represents customer-specific information.
type CustomerData struct {
	CustomerName 	string `json:"customerName,omitempty"`
	CustomerID 		string `json:"customerId"`
	SiteName     	string `json:"siteName,omitempty"`
	SiteID 			string `json:"siteId"`
}