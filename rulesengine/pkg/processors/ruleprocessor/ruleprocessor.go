package ruleprocessor

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	config "github.com/kiranpt03/factorysphere/devicesphere/engines/rulesengine/pkg/config"
	nats "github.com/kiranpt03/factorysphere/devicesphere/engines/rulesengine/pkg/connectors/nats"
	postgres "github.com/kiranpt03/factorysphere/devicesphere/engines/rulesengine/pkg/databases/postgres"
	models "github.com/kiranpt03/factorysphere/devicesphere/engines/rulesengine/pkg/models"
	log "github.com/kiranpt03/factorysphere/devicesphere/engines/rulesengine/pkg/utils/loggers"
)

type RuleProcessor struct {
	PGRepository *postgres.PostgreSQLRepository
	NatsConsumer *nats.NatsConsumer
	NatsProducer *nats.NatsProducer
	Config       *config.Config
}

func NewRuleProcessor(config *config.Config) *RuleProcessor {
	// Create a new PostgreSQL repository instance
	repository, err := postgres.NewPostgreSQLRepository(config)
	if err != nil {
		panic(err)
	}

	natsConsumerConfig := nats.NatsConsumerConfig{
		Servers:  []string{config.Nats.Server},
		GroupID:  config.Nats.ConsumerGroup,
		Subjects: []string{config.Nats.InletSubject},
	}
	consumer := nats.NewNatsConsumer(natsConsumerConfig)

	natsProducerConfig := nats.NatsProducerConfig{
		Servers: []string{config.Nats.Server},
		Subject: config.Nats.OutletSubject,
	}

	producer, err := nats.NewNatsProducer(natsProducerConfig)
	if err != nil {
		panic(err)
	}

	return &RuleProcessor{
		PGRepository: repository,
		NatsConsumer: consumer,
		NatsProducer: producer,
		Config:       config,
	}
}

func (rp *RuleProcessor) getRuleDeviceDataFromDeviceIdandPropertyId(deviceId, propertyId string) (*models.Device, error) {
	device := &models.Device{
		Properties: []models.Property{},
	}

	// Query to fetch device details.
	deviceQuery := `
                SELECT id, reference_id, type, device_name, created_at, state, location, status, customer, site
                FROM devices
                WHERE id = $1
        `

	deviceResults, err := rp.PGRepository.ExecuteQuery(deviceQuery, deviceId)
	if err != nil {
		return nil, fmt.Errorf("error querying device: %w", err)
	}

	if len(deviceResults) == 0 {
		return nil, fmt.Errorf("device with ID %s not found", deviceId)
	}

	deviceMap := deviceResults[0]

	device.ID = fmt.Sprintf("%v", deviceMap["id"])
	device.ReferenceID = fmt.Sprintf("%v", deviceMap["reference_id"])
	device.Type = fmt.Sprintf("%v", deviceMap["type"])
	device.DeviceName = fmt.Sprintf("%v", deviceMap["device_name"])
	device.CreatedAt = fmt.Sprintf("%v", deviceMap["created_at"])
	device.State = fmt.Sprintf("%v", deviceMap["state"])
	device.Location = fmt.Sprintf("%v", deviceMap["location"])
	device.Status = fmt.Sprintf("%v", deviceMap["status"])
	device.Customer = fmt.Sprintf("%v", deviceMap["customer"])
	device.Site = fmt.Sprintf("%v", deviceMap["site"])

	// Query to fetch the specific property.
	propertyQuery := `
                SELECT id, reference_id, name, unit, state, status, data_type, value, threshold
                FROM properties
                WHERE device_id = $1 AND id = $2
        `

	propertyResults, err := rp.PGRepository.ExecuteQuery(propertyQuery, deviceId, propertyId)
	if err != nil {
		return nil, fmt.Errorf("error querying property: %w", err)
	}

	if len(propertyResults) == 0 {
		return nil, fmt.Errorf("property with ID %s not found for device %s", propertyId, deviceId)
	}

	propertyMap := propertyResults[0]

	property := models.Property{
		ID:          fmt.Sprintf("%v", propertyMap["id"]),
		ReferenceID: fmt.Sprintf("%v", propertyMap["reference_id"]),
		Name:        fmt.Sprintf("%v", propertyMap["name"]),
		Unit:        fmt.Sprintf("%v", propertyMap["unit"]),
		State:       fmt.Sprintf("%v", propertyMap["state"]),
		Status:      fmt.Sprintf("%v", propertyMap["status"]),
		DataType:    fmt.Sprintf("%v", propertyMap["data_type"]),
		Value:       fmt.Sprintf("%v", propertyMap["value"]),
		Threshold:   fmt.Sprintf("%v", propertyMap["threshold"]),
	}

	device.Properties = append(device.Properties, property)

	return device, nil
}

func (rp *RuleProcessor) getRuleByID(ruleID string) (*models.Rule, error) {
	query := `
                        SELECT
                                r.id AS rule_id,
                                r.name AS rule_name,
                                r.severity AS rule_severity,
                                r.status AS rule_status,
                                r.created_at AS rule_created_at,
                                r.updated_at AS rule_updated_at,
                                c.id AS condition_id,
                                c.position AS condition_position,
                                c.type AS condition_type,
                                c.device_id AS condition_device_id,
								c.device_name AS condition_device_name,
                                c.property_id AS condition_property_id,
								c.property_name AS condition_property_name,
                                c.operator_id AS condition_operator_id,
                                c.operator_symbol AS condition_operator_symbol,
                                c.value AS condition_value
                        FROM
                                rules r
                        LEFT JOIN
                                conditions c ON r.id = c.rule_id
                        WHERE
                                r.id = $1;
                `

	results, err := rp.PGRepository.ExecuteQuery(query, ruleID)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("rule not found")
	}

	rule := &models.Rule{}
	conditions := make([]models.Condition, 0)
	firstRow := true

	for _, row := range results {
		if firstRow {
			rule.ID, _ = row["rule_id"].(string)
			rule.Name, _ = row["rule_name"].(string)
			rule.Severity, _ = row["rule_severity"].(string)
			rule.Status, _ = row["rule_status"].(string)
			rule.CreatedAt, _ = row["rule_created_at"].(string)
			rule.UpdatedAt, _ = row["rule_updated_at"].(string)
			firstRow = false
		}

		conditionID, ok := row["condition_id"].(string)
		if ok && conditionID != "" {
			conditions = append(conditions, models.Condition{
				ID:             conditionID,
				Position:       row["condition_position"].(string),
				Type:           row["condition_type"].(string),
				DeviceId:       row["condition_device_id"].(string),
				DeviceName:     row["condition_device_name"].(string),
				PropertyId:     row["condition_property_id"].(string),
				PropertyName:   row["condition_property_name"].(string),
				OperatorId:     row["condition_operator_id"].(string),
				OperatorSymbol: row["condition_operator_symbol"].(string),
				Value:          row["condition_value"].(string),
			})
		}
	}

	rule.Conditions = conditions
	return rule, nil
}

func (rp *RuleProcessor) evaluateRule(rule models.Rule, deviceData map[string]models.Device) (bool, error) {
	conditions := make(map[int]models.Condition)
	positions := make([]int, 0)

	log.Debug("Device Data: %v", deviceData)

	for _, cond := range rule.Conditions {
		pos, err := strconv.Atoi(cond.Position)
		if err != nil {
			return false, fmt.Errorf("invalid condition position: %s", cond.Position)
		}
		conditions[pos] = cond
		positions = append(positions, pos)
	}

	sort.Ints(positions) // Sort positions to ensure correct order

	values := make([]interface{}, 0)
	operators := make([]string, 0)

	for _, pos := range positions {
		cond := conditions[pos]
		switch cond.Type {
		case "sensor":
			if device, ok := deviceData[cond.DeviceId]; ok {
				for _, prop := range device.Properties {
					if prop.ID == cond.PropertyId {
						if val, err := strconv.ParseFloat(prop.Value, 64); err == nil {
							values = append(values, val)
						} else {
							return false, fmt.Errorf("invalid sensor value: %s", prop.Value)
						}
						break
					}
				}
			} else {
				return false, fmt.Errorf("device not found: %s", cond.DeviceId)
			}
		case "input":
			if val, err := strconv.ParseFloat(cond.Value, 64); err == nil {
				values = append(values, val)
			} else {
				return false, fmt.Errorf("invalid input value: %s", cond.Value)
			}
		case "operator":
			operators = append(operators, cond.OperatorSymbol)
		}
	}

	if len(values) == 0 || len(operators) == 0 {
		return false, fmt.Errorf("invalid rule condition: not enough values or operators")
	}

	result, err := rp.evaluateExpression(values, operators)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (rp *RuleProcessor) evaluateExpression(values []interface{}, operators []string) (bool, error) {
	if len(values)-1 != len(operators) {
		return false, fmt.Errorf("number of operators does not match number of values")
	}

	if len(operators) == 1 {
		return rp.evaluateSimpleComparison(values[0], operators[0], values[1])
	}

	if len(operators) == 2 {
		v1, ok1 := values[0].(float64)
		v2, ok2 := values[1].(float64)
		v3, ok3 := values[2].(float64)

		if !ok1 || !ok2 || !ok3 {
			return false, fmt.Errorf("invalid type for composite rule")
		}

		if operators[0] == "+" {
			switch operators[1] {
			case ">":
				return v1+v2 > v3, nil
			case "<":
				return v1+v2 < v3, nil
			case "==":
				return v1+v2 == v3, nil
			case ">=":
				return v1+v2 >= v3, nil
			case "<=":
				return v1+v2 <= v3, nil
			case "!=":
				return v1+v2 != v3, nil
			default:
				return false, fmt.Errorf("invalid operator in composite rule: %s", operators[1])
			}
		}
		//Add other composite operators here if needed
	}

	return false, fmt.Errorf("composite rule evaluation not implemented")
}

func (rp *RuleProcessor) evaluateSimpleComparison(left, operator, right interface{}) (bool, error) {
	leftFloat, ok1 := left.(float64)
	rightFloat, ok2 := right.(float64)

	if !ok1 || !ok2 {
		return false, fmt.Errorf("invalid types for comparison")
	}

	switch operator {
	case ">":
		return leftFloat > rightFloat, nil
	case "<":
		return leftFloat < rightFloat, nil
	case "==":
		return leftFloat == rightFloat, nil
	case "!=":
		return leftFloat != rightFloat, nil
	case ">=":
		return leftFloat >= rightFloat, nil
	case "<=":
		return leftFloat <= rightFloat, nil
	default:
		return false, fmt.Errorf("invalid operator: %s", operator)
	}
}

func (rp *RuleProcessor) getRuleIds(deviceId, propertyId string) ([]string, error) {
	query := `SELECT rule_id FROM conditions WHERE device_id = $1 AND property_id = $2 AND position = $3` // Replace your_table_name

	results, err := rp.PGRepository.ExecuteQuery(query, deviceId, propertyId, "1") // position is string 1
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no rules found")
	}

	var ruleIDs []string
	for _, result := range results {
		ruleID, ok := result["rule_id"].(string)
		if !ok {
			return nil, fmt.Errorf("rule_id not found or not a string")
		}
		ruleIDs = append(ruleIDs, ruleID)
	}

	return ruleIDs, nil
}

func (rp *RuleProcessor) dataTransformer(data string) {
	log.Debug("Processing received data: %s", data)

	var deviceRec models.Device
	err := json.Unmarshal([]byte(data), &deviceRec)
	if err != nil {
		log.Error("Error unmarshalling JSON: %v", err)
		return
	}

	log.Debug("Device model %+v\n", deviceRec)

	deviceId := deviceRec.ID
	if len(deviceRec.Properties) == 0 {
		log.Error("Device has no properties to evaluate")
		return
	}
	propertyId := deviceRec.Properties[0].ID

	// All rule ID's associated with deviceID and PropertyID
	ruleIds, ruleErr := rp.getRuleIds(deviceId, propertyId)
	if ruleErr != nil {
		log.Error("Unable to parse rule data %v", ruleErr)
		return
	}
	log.Debug("Rule Id's %v", ruleIds)

	var rules []models.Rule
	rulesDeviceMap := make(map[string]map[string]models.Device)
	for _, ruleId := range ruleIds {
		log.Debug("Rule ID: %s", ruleId)
		ruleData, dataErr := rp.getRuleByID(ruleId)
		if dataErr != nil {
			log.Error("Unable to parse rule data %v", dataErr)
			return
		}

		deviceMap := make(map[string]models.Device)
		deviceMap[deviceId] = deviceRec
		fmt.Println(len(ruleData.Conditions))
		for _, ruleCondition := range ruleData.Conditions {
			log.Debug("Rule Condition: %v", ruleCondition)
			if ruleCondition.Type == "sensor" {
				condDeviceId := ruleCondition.DeviceId
				log.Debug("Condition device Id: %s", condDeviceId)
				condPropertyId := ruleCondition.PropertyId
				log.Debug("Condition property  Id: %s", condPropertyId)

				if condDeviceId != deviceId {
					device, devErr := rp.getRuleDeviceDataFromDeviceIdandPropertyId(condDeviceId, condPropertyId)
					if devErr != nil {
						log.Error("Error while fetching device from deviceId and propertyIf: %v", devErr)
					}
					deviceMap[condDeviceId] = *device
				}
			}
			rulesDeviceMap[ruleId] = deviceMap
		}

		rules = append(rules, *ruleData)
	}

	for _, rule := range rules {
		deviceData := rulesDeviceMap[rule.ID]
		result, err := rp.evaluateRule(rule, deviceData)
		if err != nil {
			fmt.Printf("Error evaluating rule %s: %v\n", rule.ID, err)
		} else {
			fmt.Printf("Rule %s evaluation result: %v\n", rule.ID, result)
			if result {
				ruleModelByte, marshalErr := json.Marshal(rule)
				if marshalErr != nil {
					log.Error("Error marshalling device model: %v", marshalErr)
					continue // Skip to the next property if marshalling fails
				}
				producer := rp.NatsProducer
				// defer producer.Close()
				err = producer.Produce(context.Background(), []byte(rule.ID), ruleModelByte)
				if err != nil {
					log.Error("Error producing message to Kafka: %v", err)
				}
			}
		}
	}

}

func (rp *RuleProcessor) ProcessData() {
	log.Info("Processing the data")
	consumer := rp.NatsConsumer
	defer consumer.Close()

	producer := rp.NatsProducer
	defer producer.Close()

	rawDataTopicChan := consumer.Consume(rp.Config.Nats.InletSubject)
	// processedDataTopicChan := consumer.Consume("other_subject")

	for {
		select {
		case msg := <-rawDataTopicChan:
			go rp.dataTransformer(string(msg.Data))
			// case msg := <-processedDataTopicChan:
			// 	log.Debug("Received message topic: processed_data: %s", msg.Data)
		}
	}
}
