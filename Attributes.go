package xacml

import (
	"errors"
	"strconv"
)

func resolveAttributeValue(attributeValueMap map[string]interface{}, request Request) (interface{}, error) {
	attributeValueString, ok := attributeValueMap["#text"].(string)
	if !ok {
		return nil, errors.New("attributeValueString not ok")
	}
	attributeValueDataType, ok := attributeValueMap["-DataType"].(string)
	if !ok {
		return nil, errors.New("attributeValueDataType not ok")
	}
	return convertAttributeToDataType(attributeValueString, attributeValueDataType)
}

func resolveAttributeDesignator(attributeDesignatorMap map[string]interface{}, request Request) (interface{}, error) {
	attributeDesignatorMustBePresentString, ok := attributeDesignatorMap["-MustBePresent"].(string)
	if !ok {
		attributeDesignatorMustBePresentString = "false"
	}
	attributeDesignatorMustBePresent, err := strconv.ParseBool(attributeDesignatorMustBePresentString)
	if err != nil {
		return nil, err
	}
	attributeDesignatorAttributeId, ok := attributeDesignatorMap["-AttributeId"].(string)
	if !ok {
		return nil, errors.New("attributeDesignatorAttributeId not ok")
	}
	attributeDesignatorCategory, ok := attributeDesignatorMap["-Category"].(string)
	if !ok {
		return nil, errors.New("attributeDesignatorCategory not ok")
	}
	attributeDesignatorDataType, ok := attributeDesignatorMap["-DataType"].(string)
	if err != nil {
		return nil, errors.New("attributeDesignatorDataType not ok")
	}
	return request.GetAttribute(attributeDesignatorCategory, attributeDesignatorAttributeId, attributeDesignatorDataType, attributeDesignatorMustBePresent)
}

func convertAttributeToDataType(attributeValue, dataType string) (interface{}, error) {
	switch dataType {
	case DataTypeString:
		return attributeValue, nil
	case DataTypeBoolean:
		return strconv.ParseBool(attributeValue)
	case DataTypeInteger:
		return strconv.ParseInt(attributeValue, 10, 64)
	case DataTypeDouble:
		return strconv.ParseFloat(attributeValue, 64)
	default:
		return attributeValue, errors.New("Unknown or Unsupported DataType")
	}
}

func defaultValueForDataType(dataType string) interface{} {
	switch dataType {
	case DataTypeString:
		return ""
	case DataTypeBoolean:
		return false
	case DataTypeInteger:
		return int64(0)
	case DataTypeDouble:
		return float64(0)
	default:
		return ""
	}
}
