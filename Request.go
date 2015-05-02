package xacml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"github.com/clbanning/mxj"
)

const (
	AttributeCategorySubjectAccessSubject       = "urn:oasis:names:tc:xacml:1.0:subject-category:access-subject"
	AttributeCategorySubjectCodebase            = "urn:oasis:names:tc:xacml:1.0:subject-category:codebase"
	AttributeCategorySubjectIntermediarySubject = "urn:oasis:names:tc:xacml:1.0:subject-category:intermediary-subject"
	AttributeCategorySubjectRecipientSubject    = "urn:oasis:names:tc:xacml:1.0:subject-category:recipient-subject"
	AttributeCategorySubjectRequestingMachine   = "urn:oasis:names:tc:xacml:1.0:subject-category:requesting-machine"
	AttributeCategoryResource                   = "urn:oasis:names:tc:xacml:3.0:attribute-category:resource"
	AttributeCategoryAction                     = "urn:oasis:names:tc:xacml:3.0:attribute-category:action"
	AttributeCategoryEnvironment                = "urn:oasis:names:tc:xacml:3.0:attribute-category:environment"
	IdentifierSubjectAuthNLocalityDNSName       = "urn:oasis:names:tc:xacml:1.0:subject:authn-locality:dns-name"
	IdentifierSubjectAuthNLocalityIPAddress     = "urn:oasis:names:tc:xacml:1.0:subject:authn-locality:ip-address"
	IdentifierSubjectAuthenticationMethod       = "urn:oasis:names:tc:xacml:1.0:subject:authentication-method"
	IdentifierSubjectAuthenticationTime         = "urn:oasis:names:tc:xacml:1.0:subject:authentication-time"
	IdentifierSubjectKeyInfo                    = "urn:oasis:names:tc:xacml:1.0:subject:key-info"
	IdentifierSubjectRequestTime                = "urn:oasis:names:tc:xacml:1.0:subject:request-time"
	IdentifierSubjectSessionStartTime           = "urn:oasis:names:tc:xacml:1.0:subject:session-start-time"
	IdentifierSubjectId                         = "urn:oasis:names:tc:xacml:1.0:subject:subject-id"
	IdentifierSubjectIdQualifier                = "urn:oasis:names:tc:xacml:1.0:subject:subject-id-qualifier"
	IdentifierResourceLocation                  = "urn:oasis:names:tc:xacml:1.0:resource:resource-location"
	IdentifierResourceId                        = "urn:oasis:names:tc:xacml:1.0:resource:resource-id"
	IdentifierResourceSimpleFileName            = "urn:oasis:names:tc:xacml:1.0:resource:simple-file-name"
	IdentifierActionId                          = "urn:oasis:names:tc:xacml:1.0:action:action-id"
	IdentifierActionImpliedAction               = "urn:oasis:names:tc:xacml:1.0:action:implied-action"
	IdentifierEnvironmentCurrentTime            = "urn:oasis:names:tc:xacml:1.0:environment:current-time"
	IdentifierEnvironmentCurrentDate            = "urn:oasis:names:tc:xacml:1.0:environment:current-date"
	IdentifierEnvironmentCurrentDateTime        = "urn:oasis:names:tc:xacml:1.0:environment:current-dateTime"
)

func NewRequest() Request {
	return Request{make(map[string]interface{})}
}

type Request struct {
	request map[string]interface{}
}

func (request Request) AddAttribute(category, attributeId, dataType, value string) {
	var buff bytes.Buffer
	xml.EscapeText(&buff, []byte(value))
	value = buff.String()

	attributeValue := make(map[string]interface{})
	attributeValue["-DataType"] = dataType
	attributeValue["#text"] = value
	attribute := make(map[string]interface{})
	attribute["-AttributeId"] = attributeId
	attribute["AttributeValue"] = attributeValue

	requestInterface, _ := request.request["Request"].(map[string]interface{})
	if requestInterface == nil {
		requestInterface = make(map[string]interface{})
		request.request["Request"] = requestInterface
	}

	categoryInterface := request.getCategory(category)
	if categoryInterface == nil {
		categoryInterface := make(map[string]interface{})
		categoryInterface["-Category"] = category
		categoryInterface["Attribute"] = attribute
		//TODO Add Category to Attributes
		attributesInterface := requestInterface["Attributes"]
		switch attributesInterface.(type) {
		case nil:
			requestInterface["Attributes"] = categoryInterface
		case []interface{}:
			attributesArray := attributesInterface.([]interface{})
			attributesArray = append(attributesArray, categoryInterface)
			requestInterface["Attributes"] = attributesArray
		case map[string]interface{}:
			attributesArray := make([]interface{}, 0)
			attributesArray = append(attributesArray, attributesInterface)
			attributesArray = append(attributesArray, categoryInterface)
			requestInterface["Attributes"] = attributesArray
		}
	} else {
		//Add attribute to attribute
		switch categoryInterface["Attribute"].(type) {
		case []interface{}:
			//append attribute to attribute array
			attributeArray, _ := categoryInterface["Attribute"].([]interface{})
			categoryInterface["Attribute"] = append(attributeArray, attribute)
		case map[string]interface{}:
			attributeArray := make([]interface{}, 0)
			attributeArray = append(attributeArray, categoryInterface["Attribute"])
			attributeArray = append(attributeArray, attribute)
			categoryInterface["Attribute"] = attributeArray
		}
	}
}

func (request Request) AddAttributeBag(category, attributeId, dataType string, values []string) {
	attributeValues := make([]interface{}, 0)
	for _, value := range values {
		attributeValue := make(map[string]interface{})
		attributeValue["-DataType"] = dataType
		attributeValue["#text"] = value
		attributeValues = append(attributeValues, attributeValue)
	}

	attribute := make(map[string]interface{})
	attribute["-AttributeId"] = attributeId
	attribute["AttributeValue"] = attributeValues

	requestInterface, _ := request.request["Request"].(map[string]interface{})
	if requestInterface == nil {
		requestInterface = make(map[string]interface{})
		request.request["Request"] = requestInterface
	}

	categoryInterface := request.getCategory(category)
	if categoryInterface == nil {
		categoryInterface := make(map[string]interface{})
		categoryInterface["-Category"] = category
		categoryInterface["Attribute"] = attribute
		//TODO Add Category to Attributes
		attributesInterface := requestInterface["Attributes"]
		switch attributesInterface.(type) {
		case nil:
			requestInterface["Attributes"] = categoryInterface
		case []interface{}:
			attributesArray := attributesInterface.([]interface{})
			attributesArray = append(attributesArray, categoryInterface)
			requestInterface["Attributes"] = attributesArray
		case map[string]interface{}:
			attributesArray := make([]interface{}, 0)
			attributesArray = append(attributesArray, attributesInterface)
			attributesArray = append(attributesArray, categoryInterface)
			requestInterface["Attributes"] = attributesArray
		}
	} else {
		//Add attribute to attribute
		switch categoryInterface["Attribute"].(type) {
		case []interface{}:
			//append attribute to attribute array
			attributeArray, _ := categoryInterface["Attribute"].([]interface{})
			categoryInterface["Attribute"] = append(attributeArray, attribute)
		case map[string]interface{}:
			attributeArray := make([]interface{}, 0)
			attributeArray = append(attributeArray, categoryInterface["Attribute"])
			attributeArray = append(attributeArray, attribute)
			categoryInterface["Attribute"] = attributeArray
		}
	}
}

func (request Request) getCategory(category string) map[string]interface{} {
	if request.request == nil {
		return nil
	}
	requestInterface, _ := request.request["Request"].(map[string]interface{})
	attributesInterface := requestInterface["Attributes"]

	switch attributesInterface.(type) {
	case []interface{}:
		attributesArray := attributesInterface.([]interface{})
		for _, attributesObjInterface := range attributesArray {
			if attributesObj, ok := attributesObjInterface.(map[string]interface{}); ok {
				if cat, ok := attributesObj["-Category"].(string); ok && cat == category {
					return attributesObj
				}
			}
		}
	case map[string]interface{}:
		attributesObj := attributesInterface.(map[string]interface{})
		if cat, ok := attributesObj["-Category"].(string); ok && cat == category {
			return attributesObj
		}
	}
	return nil
}

func getAttributeFromAttributeMap(attributeMap map[string]interface{}, attributeId, attributeDataType string, mustBePresent bool) (interface{}, error) {
	if attrId, ok := attributeMap["-AttributeId"].(string); ok {
		if attrId == attributeId {
			if attributeValueObj, ok := attributeMap["AttributeValue"].(map[string]interface{}); ok {
				attributeValue, ok := attributeValueObj["#text"].(string)
				if !ok {
					return nil, errors.New("Error parsing Attribute Value")
				}
				dataType, ok := attributeValueObj["-DataType"].(string)
				if !ok {
					return nil, errors.New("Error parsing Attribute DataType")
				}
				if dataType != attributeDataType {
					return nil, errors.New("Expected DataType and found DataType do not match")
				}
				return convertAttributeToDataType(attributeValue, dataType)
			} else if attributeValueArray, ok := attributeMap["AttributeValue"].([]interface{}); ok {
				ret := make([]interface{}, 0)
				for _, attributeValueInterface := range attributeValueArray {
					if attributeValueObj, ok := attributeValueInterface.(map[string]interface{}); ok {
						attributeValue, ok := attributeValueObj["#text"].(string)
						if !ok {
							continue
						}
						dataType, ok := attributeValueObj["-DataType"].(string)
						if !ok {
							continue
						}
						if dataType != attributeDataType {
							continue
						}
						toAdd, _ := convertAttributeToDataType(attributeValue, dataType)
						ret = append(ret, toAdd)
					}
				}
				return ret, nil
			} else {
				return nil, errors.New("Error parsing Attribute Value Object")
			}
		} else {
			return defaultValueForDataType(attributeDataType), errors.New("Attribute Ids do not match")
		}
	}
	return nil, errors.New("Attribute Not Found in attribute map")
}

func (request Request) GetAttribute(category, attributeId, attributeDataType string, mustBePresent bool) (interface{}, error) {
	categoryObj := request.getCategory(category)
	if categoryObj == nil {
		if mustBePresent {
			return nil, errors.New("Category not included in request")
		} else {
			return defaultValueForDataType(attributeDataType), nil
		}
	}

	attributeArrayInterface := categoryObj["Attribute"]
	switch attributeArrayInterface.(type) {
	case map[string]interface{}:
		return getAttributeFromAttributeMap(attributeArrayInterface.(map[string]interface{}), attributeId, attributeDataType, mustBePresent)
	case []interface{}:
		for _, attributeObjInterface := range attributeArrayInterface.([]interface{}) {
			if attributeObj, ok := attributeObjInterface.(map[string]interface{}); ok {
				response, err := getAttributeFromAttributeMap(attributeObj, attributeId, attributeDataType, mustBePresent)
				if err == nil {
					return response, err
				}
			}
		}
	}
	if mustBePresent {
		return nil, errors.New("Attribute not found")
	} else {
		return defaultValueForDataType(attributeDataType), nil
	}
}

func (request Request) Xml() string {
	//TODO Make sure the xacml3 headers are on the request
	requestInterface, _ := request.request["Request"].(map[string]interface{})
	if requestInterface == nil {
		requestInterface = make(map[string]interface{})
		request.request["Request"] = requestInterface
	}
	requestInterface["-xmlns"] = "urn:oasis:names:tc:xacml:3.0:core:schema:wd-17"
	requestInterface["-xmlns:xsi"] = "http://www.w3.org/2001/XMLSchema-instance"
	requestInterface["-xsi:schemaLocation"] = "urn:oasis:names:tc:xacml:3.0:core:schema:wd-17 http://docs.oasis-open.org/xacml/3.0/xacml-core-v3-schema-wd-17.xsd"
	m := mxj.Map(request.request)
	ret, _ := m.Xml()
	return `<?xml version="1.0" encoding="UTF-8"?>` + string(ret)
}
