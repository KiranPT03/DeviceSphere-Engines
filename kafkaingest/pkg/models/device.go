package models

type Device struct {
    ID           string `json:"id"`
    ReferenceID  string `json:"referenceId" validate:"required"`
    Type         string `json:"type"`
    DeviceName   string `json:"deviceName"`
    CreatedAt    string `json:"createdAt"`
    State        string `json:"state" validate:"required"`
    Location     string `json:"location"`
    Status       string `json:"status" validate:"required"`
    Customer     string `json:"customer"`
    Site         string `json:"site"`
    Properties   []Property `json:"properties" validate:"required,dive"`
}

type Property struct {
    ID           string `json:"id"`
    ReferenceID  string `json:"referenceId" validate:"required"`
    Name         string `json:"name" validate:"required"`
    Unit         string `json:"unit"`
    State        string `json:"state" validate:"required"`
    Status       string `json:"status" validate:"required"`
    DataType     string `json:"dataType"`
    Value        string `json:"value"`
    Threshold    string `json:"threshold" validate:"required"`
}