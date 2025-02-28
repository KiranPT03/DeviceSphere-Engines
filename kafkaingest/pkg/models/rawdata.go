package models

type RawProperty struct {
    PropertyRefID  string `json:"referenceId" validate:"required"`
    Value        string `json:"value"`
}

type RawData struct {
	DeviceRefID  string 		`json:"deviceRefId"`
	Properties   []RawProperty 	`json:"properties"`
}