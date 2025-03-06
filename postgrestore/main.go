package main

import (
	config "github.com/kiranpt03/factorysphere/devicesphere/engines/postgrestore/pkg/config"
	postgreprocessor "github.com/kiranpt03/factorysphere/devicesphere/engines/postgrestore/pkg/processors/postgreprocessors"
	log "github.com/kiranpt03/factorysphere/devicesphere/engines/postgrestore/pkg/utils/loggers"
)

func main() {
	log.Info("Strating application...")
	cfg, cfgErr := config.GetConfig()
	if cfgErr != nil {
		log.Error("Unable to read config: %v", cfgErr)
	}

	postgreProcessor := postgreprocessor.NewPostgreProcessor(cfg)
	postgreProcessor.ProcessData()
}

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"sort"
// 	"strconv"
// )

// // Rule represents a rule from your JSON data.
// type Rule struct {
// 	ID          string      `json:"id"`
// 	Name        string      `json:"name"`
// 	Severity    string      `json:"severity"`
// 	Status      string      `json:"status"`
// 	Type        string      `json:"type"`
// 	Description string      `json:"description"`
// 	CreatedAt   string      `json:"createdAt"`
// 	UpdatedAt   string      `json:"updatedAt"`
// 	Conditions  []Condition `json:"condition"`
// }

// // Condition represents a condition within a rule.
// type Condition struct {
// 	ID             string `json:"id"`
// 	Position       string `json:"position"`
// 	Type           string `json:"type"`
// 	DeviceID       string `json:"deviceId,omitempty"`
// 	PropertyID     string `json:"propertyId,omitempty"`
// 	OperatorID     string `json:"operatorId,omitempty"`
// 	OperatorSymbol string `json:"operatorSymbol,omitempty"`
// 	Value          string `json:"value,omitempty"`
// }

// // Device represents your device model.
// type Device struct {
// 	ID          string     `json:"id"`
// 	ReferenceID string     `json:"referenceId"`
// 	Type        string     `json:"type"`
// 	DeviceName  string     `json:"deviceName"`
// 	CreatedAt   string     `json:"createdAt"`
// 	State       string     `json:"state"`
// 	Location    string     `json:"location"`
// 	Status      string     `json:"status"`
// 	Customer    string     `json:"customer"`
// 	Site        string     `json:"site"`
// 	Properties  []Property `json:"properties"`
// }

// // Property represents a property of a device.
// type Property struct {
// 	ID          string `json:"id"`
// 	ReferenceID string `json:"referenceId"`
// 	Name        string `json:"name"`
// 	Unit        string `json:"unit"`
// 	State       string `json:"state"`
// 	Status      string `json:"status"`
// 	DataType    string `json:"dataType"`
// 	Value       string `json:"value"`
// 	Threshold   string `json:"threshold"`
// }

// // Sample device data (replace with your actual device data retrieval).
// var deviceData = map[string]Device{
// 	"0773341b-ab16-4324-b970-ba17a397fbce": {
// 		ID:          "0773341b-ab16-4324-b970-ba17a397fbce",
// 		ReferenceID: "ref-awtbwrc4w",
// 		Type:        "sensor",
// 		DeviceName:  "PLC_001",
// 		CreatedAt:   "2025-03-01T09:36:02.325Z",
// 		State:       "active",
// 		Location:    "Sangli, Maharashtra",
// 		Status:      "online",
// 		Customer:    "SphereFactory",
// 		Site:        "Sangli MIDC",
// 		Properties: []Property{
// 			{
// 				ID:          "18801602-15d3-44b6-a15a-ade3fa2a00a8",
// 				ReferenceID: "ref-awtbwrc4w",
// 				Name:        "Temperature",
// 				Unit:        "Celsius",
// 				State:       "active",
// 				Status:      "online",
// 				DataType:    "float",
// 				Value:       "65",
// 				Threshold:   "90",
// 			},
// 		},
// 	},
// 	"10c8fb9b-0704-47e3-88df-49344e865eeb": {
// 		ID:          "10c8fb9b-0704-47e3-88df-49344e865eeb",
// 		ReferenceID: "ref-other",
// 		Type:        "sensor",
// 		DeviceName:  "PLC_002",
// 		CreatedAt:   "2025-03-01T09:37:02.325Z",
// 		State:       "active",
// 		Location:    "Mumbai, Maharashtra",
// 		Status:      "online",
// 		Customer:    "SphereFactory",
// 		Site:        "Mumbai MIDC",
// 		Properties: []Property{
// 			{
// 				ID:          "1aedde00-5d38-4fe6-bfe2-ed3cf15f7a65",
// 				ReferenceID: "ref-other",
// 				Name:        "Temperature",
// 				Unit:        "Celsius",
// 				State:       "active",
// 				Status:      "online",
// 				DataType:    "float",
// 				Value:       "30",
// 				Threshold:   "90",
// 			},
// 		},
// 	},
// }

// func evaluateRule(rule Rule) (bool, error) {
// 	conditions := make(map[int]Condition)
// 	positions := make([]int, 0)

// 	for _, cond := range rule.Conditions {
// 		pos, err := strconv.Atoi(cond.Position)
// 		if err != nil {
// 			return false, fmt.Errorf("invalid condition position: %s", cond.Position)
// 		}
// 		conditions[pos] = cond
// 		positions = append(positions, pos)
// 	}

// 	sort.Ints(positions) // Sort positions to ensure correct order

// 	values := make([]interface{}, 0)
// 	operators := make([]string, 0)

// 	for _, pos := range positions {
// 		cond := conditions[pos]
// 		switch cond.Type {
// 		case "sensor":
// 			if device, ok := deviceData[cond.DeviceID]; ok {
// 				for _, prop := range device.Properties {
// 					if prop.ID == cond.PropertyID {
// 						if val, err := strconv.ParseFloat(prop.Value, 64); err == nil {
// 							values = append(values, val)
// 						} else {
// 							return false, fmt.Errorf("invalid sensor value: %s", prop.Value)
// 						}
// 						break
// 					}
// 				}
// 			} else {
// 				return false, fmt.Errorf("device not found: %s", cond.DeviceID)
// 			}
// 		case "input":
// 			if val, err := strconv.ParseFloat(cond.Value, 64); err == nil {
// 				values = append(values, val)
// 			} else {
// 				return false, fmt.Errorf("invalid input value: %s", cond.Value)
// 			}
// 		case "operator":
// 			operators = append(operators, cond.OperatorSymbol)
// 		}
// 	}

// 	if len(values) == 0 || len(operators) == 0 {
// 		return false, fmt.Errorf("invalid rule condition: not enough values or operators")
// 	}

// 	result, err := evaluateExpression(values, operators)
// 	if err != nil {
// 		return false, err
// 	}
// 	return result, nil
// }

// func evaluateExpression(values []interface{}, operators []string) (bool, error) {
// 	if len(values)-1 != len(operators) {
// 		return false, fmt.Errorf("number of operators does not match number of values")
// 	}

// 	if len(operators) == 1 {
// 		return evaluateSimpleComparison(values[0], operators[0], values[1])
// 	}

// 	if len(operators) == 2 {
// 		v1, ok1 := values[0].(float64)
// 		v2, ok2 := values[1].(float64)
// 		v3, ok3 := values[2].(float64)

// 		if !ok1 || !ok2 || !ok3 {
// 			return false, fmt.Errorf("invalid type for composite rule")
// 		}

// 		if operators[0] == "+" {
// 			switch operators[1] {
// 			case ">":
// 				return v1+v2 > v3, nil
// 			case "<":
// 				return v1+v2 < v3, nil
// 			case "==":
// 				return v1+v2 == v3, nil
// 			case ">=":
// 				return v1+v2 >= v3, nil
// 			case "<=":
// 				return v1+v2 <= v3, nil
// 			case "!=":
// 				return v1+v2 != v3, nil
// 			default:
// 				return false, fmt.Errorf("invalid operator in composite rule: %s", operators[1])
// 			}
// 		}
// 		//Add other composite operators here if needed
// 	}

// 	return false, fmt.Errorf("composite rule evaluation not implemented")
// }

// func evaluateSimpleComparison(left, operator, right interface{}) (bool, error) {
// 	leftFloat, ok1 := left.(float64)
// 	rightFloat, ok2 := right.(float64)

// 	if !ok1 || !ok2 {
// 		return false, fmt.Errorf("invalid types for comparison")
// 	}

// 	switch operator {
// 	case ">":
// 		return leftFloat > rightFloat, nil
// 	case "<":
// 		return leftFloat < rightFloat, nil
// 	case "==":
// 		return leftFloat == rightFloat, nil
// 	case "!=":
// 		return leftFloat != rightFloat, nil
// 	case ">=":
// 		return leftFloat >= rightFloat, nil
// 	case "<=":
// 		return leftFloat <= rightFloat, nil
// 	default:
// 		return false, fmt.Errorf("invalid operator: %s", operator)
// 	}
// }

// func main() {
// 	jsonData := `[
//                 {
//                         "id": "b8989173-49ad-46b0-8b69-c65527b47c75",
//                         "name": "Rule_001_Simple_Rule",
//                         "severity": "warning",
//                         "status": "active",
//                         "type": "sensor",
//                         "description": "Temperature more that 60",
//                         "createdAt": "2025-03-01T09:49:21.307Z",
//                         "updatedAt": "2025-03-01T09:49:21.307Z",
//                         "condition": [
//                                 {
//                                         "id": "cabcb962-6dd0-4ae4-a6ef-cb10d7096093",
//                                         "position": "1",
//                                         "type": "sensor",
//                                         "deviceId": "0773341b-ab16-4324-b970-ba17a397fbce",
//                                         "propertyId": "18801602-15d3-44b6-a15a-ade3fa2a00a8"
//                                 },
//                                 {
//                                         "id": "b22b8028-791d-4f97-9ec1-82e10efabfc5",
//                                         "position": "2",
//                                         "type": "operator",
//                                         "operatorId": "operator-1",
//                                         "operatorSymbol": "=="
//                                 },
//                                 {
//                                         "id": "285003bd-fd00-4291-8504-93660dd6b725",
//                                         "position": "3",
//                                         "type": "input",
//                                         "value": "60"
//                                 }
//                         ]
//                 },
//                 {
//                         "id": "10c36b8a-3ca6-4392-af66-81e83d1389ad",
//                         "name": "Rule_002_Simple_Rule_two_sensors",
//                         "severity": "critical",
//                         "status": "active",
//                         "type": "sensor",
//                         "description": "Compare temperature senor from PLC_001 to temperature sensor PLC_002",
//                         "createdAt": "2025-03-01T09:50:37.791Z",
//                         "updatedAt": "2025-03-01T09:50:37.791Z",
//                         "condition": [
//                                 {
//                                         "id": "275de4a3-5aff-46bd-858f-e093b7e4a1c6",
//                                         "position": "1",
//                                         "type": "sensor",
//                                         "deviceId": "0773341b-ab16-4324-b970-ba17a397fbce",
//                                         "propertyId": "18801602-15d3-44b6-a15a-ade3fa2a00a8"
//                                 },
//                                 {
//                                         "id": "80d758f6-0698-4947-b617-e95599624050",
//                                         "position": "2",
//                                         "type": "operator",
//                                         "operatorId": "operator-1",
//                                         "operatorSymbol": ">"
//                                 },
//                                 {
//                                         "id": "98e4e9fb-dd3a-4537-a6dc-2cadfdd6d7e3",
//                                         "position": "3",
//                                         "type": "sensor",
//                                         "deviceId": "10c8fb9b-0704-47e3-88df-49344e865eeb",
//                                         "propertyId": "1aedde00-5d38-4fe6-bfe2-ed3cf15f7a65"
//                                 }
//                         ]
//                 },
//                 {
//                         "id": "86efc0a6-c655-4e5e-a96b-0066f9b4a8e2",
//                         "name": "Rule_003_Composite_Rule",
//                         "severity": "warning",
//                         "status": "active",
//                         "type": "sensor",
//                         "description": "Failing the composite rule",
//                         "createdAt": "2025-03-01T09:51:34.007Z",
//                         "updatedAt": "2025-03-01T09:51:34.007Z",
//                         "condition": [
//                                 {
//                                         "id": "fcf3b9b7-b5e7-4702-8874-d8e96ae2b2af",
//                                         "position": "1",
//                                         "type": "sensor",
//                                         "deviceId": "0773341b-ab16-4324-b970-ba17a397fbce",
//                                         "propertyId": "18801602-15d3-44b6-a15a-ade3fa2a00a8"
//                                 },
//                                 {
//                                         "id": "fb331496-5d28-4ecb-b45f-28a9002524ed",
//                                         "position": "2",
//                                         "type": "operator",
//                                         "operatorId": "operator-1",
//                                         "operatorSymbol": "+"
//                                 },
//                                 {
//                                         "id": "7f4ce0dc-dff6-4b61-a373-9fcdeb2d9b7e",
//                                         "position": "3",
//                                         "type": "sensor",
//                                         "deviceId": "10c8fb9b-0704-47e3-88df-49344e865eeb",
//                                         "propertyId": "1aedde00-5d38-4fe6-bfe2-ed3cf15f7a65"
//                                 },
//                                 {
//                                         "id": "3c339712-69d2-4bd4-ba66-159acad3ccca",
//                                         "position": "4",
//                                         "type": "operator",
//                                         "operatorId": "operator-3",
//                                         "operatorSymbol": ">"
//                                 },
//                                 {
//                                         "id": "4b46e349-1d1e-4a5f-b8ac-1bc4b1197f5b",
//                                         "position": "5",
//                                         "type": "input",
//                                         "value": "90"
//                                 }
//                         ]
//                 },
//                 {
//                         "id": "86efc0a6-c655-4e5e-a96b-0066f9b4a8e3",
//                         "name": "Rule_004_Composite_Rule_Less",
//                         "severity": "warning",
//                         "status": "active",
//                         "type": "sensor",
//                        "description": "Failing the composite rule with less than",
//                                                                 "createdAt": "2025-03-01T09:51:34.007Z",
//                                                                 "updatedAt": "2025-03-01T09:51:34.007Z",
//                                                                 "condition": [
//                                                                         {
//                                                                                 "id": "fcf3b9b7-b5e7-4702-8874-d8e96ae2b2a1",
//                                                                                 "position": "1",
//                                                                                 "type": "sensor",
//                                                                                 "deviceId": "0773341b-ab16-4324-b970-ba17a397fbce",
//                                                                                 "propertyId": "18801602-15d3-44b6-a15a-ade3fa2a00a8"
//                                                                         },
//                                                                         {
//                                                                                 "id": "fb331496-5d28-4ecb-b45f-28a9002524e1",
//                                                                                 "position": "2",
//                                                                                 "type": "operator",
//                                                                                 "operatorId": "operator-1",
//                                                                                 "operatorSymbol": "+"
//                                                                         },
//                                                                         {
//                                                                                 "id": "7f4ce0dc-dff6-4b61-a373-9fcdeb2d9b71",
//                                                                                 "position": "3",
//                                                                                 "type": "sensor",
//                                                                                 "deviceId": "10c8fb9b-0704-47e3-88df-49344e865eeb",
//                                                                                 "propertyId": "1aedde00-5d38-4fe6-bfe2-ed3cf15f7a65"
//                                                                         },
//                                                                         {
//                                                                                 "id": "3c339712-69d2-4bd4-ba66-159acad3ccca1",
//                                                                                 "position": "4",
//                                                                                 "type": "operator",
//                                                                                 "operatorId": "operator-3",
//                                                                                 "operatorSymbol": "<"
//                                                                         },
//                                                                         {
//                                                                                 "id": "4b46e349-1d1e-4a5f-b8ac-1bc4b1197f5b1",
//                                                                                 "position": "5",
//                                                                                 "type": "input",
//                                                                                 "value": "100"
//                                                                         }
//                                                                 ]
//                                                         },
//                                                         {
//                                                                 "id": "86efc0a6-c655-4e5e-a96b-0066f9b4a8e4",
//                                                                 "name": "Rule_005_Simple_Equal",
//                                                                 "severity": "warning",
//                                                                 "status": "active",
//                                                                 "type": "sensor",
//                                                                 "description": "Simple equal rule",
//                                                                 "createdAt": "2025-03-01T09:51:34.007Z",
//                                                                 "updatedAt": "2025-03-01T09:51:34.007Z",
//                                                                 "condition": [
//                                                                         {
//                                                                                 "id": "cabcb962-6dd0-4ae4-a6ef-cb10d7096094",
//                                                                                 "position": "1",
//                                                                                 "type": "sensor",
//                                                                                 "deviceId": "0773341b-ab16-4324-b970-ba17a397fbce",
//                                                                                 "propertyId": "18801602-15d3-44b6-a15a-ade3fa2a00a8"
//                                                                         },
//                                                                         {
//                                                                                 "id": "b22b8028-791d-4f97-9ec1-82e10efabfc6",
//                                                                                 "position": "2",
//                                                                                 "type": "operator",
//                                                                                 "operatorId": "operator-1",
//                                                                                 "operatorSymbol": "=="
//                                                                         },
//                                                                         {
//                                                                                 "id": "285003bd-fd00-4291-8504-93660dd6b726",
//                                                                                 "position": "3",
//                                                                                 "type": "input",
//                                                                                 "value": "65"
//                                                                         }
//                                                                 ]
//                                                         },
//                                                         {
//                                                                 "id": "86efc0a6-c655-4e5e-a96b-0066f9b4a8e5",
//                                                                 "name": "Rule_006_Simple_Less",
//                                                                 "severity": "warning",
//                                                                 "status": "active",
//                                                                 "type": "sensor",
//                                                                 "description": "Simple less rule",
//                                                                 "createdAt": "2025-03-01T09:51:34.007Z",
//                                                                 "updatedAt": "2025-03-01T09:51:34.007Z",
//                                                                 "condition": [
//                                                                         {
//                                                                                 "id": "cabcb962-6dd0-4ae4-a6ef-cb10d7096095",
//                                                                                 "position": "1",
//                                                                                 "type": "sensor",
//                                                                                 "deviceId": "0773341b-ab16-4324-b970-ba17a397fbce",
//                                                                                 "propertyId": "18801602-15d3-44b6-a15a-ade3fa2a00a8"
//                                                                         },
//                                                                         {
//                                                                                 "id": "b22b8028-791d-4f97-9ec1-82e10efabfc7",
//                                                                                 "position": "2",
//                                                                                 "type": "operator",
//                                                                                 "operatorId": "operator-1",
//                                                                                 "operatorSymbol": "<"
//                                                                         },
//                                                                         {
//                                                                                 "id": "285003bd-fd00-4291-8504-93660dd6b727",
//                                                                                 "position": "3",
//                                                                                 "type": "input",
//                                                                                 "value": "70"
//                                                                         }
//                                                                 ]
//                                                         }
//                                 ]`

// 	var rules []Rule
// 	err := json.Unmarshal([]byte(jsonData), &rules)
// 	if err != nil {
// 		fmt.Println("Error unmarshaling JSON:", err)
// 		return
// 	}

// 	for _, rule := range rules {
// 		result, err := evaluateRule(rule)
// 		if err != nil {
// 			fmt.Printf("Error evaluating rule %s: %v\n", rule.ID, err)
// 		} else {
// 			fmt.Printf("Rule %s evaluation result: %v\n", rule.ID, result)
// 		}
// 	}
// }
