package xacml

import (
	"github.com/clbanning/mxj"
)

type Response struct {
	response map[string]interface{}
}

//<Response>
//	<Result>
//		<Decision>
//		<Status>

const (
	ResponseDeny            = "Deny"
	ResponsePermit          = "Permit"
	ResponseNotApplicable   = "Not Applicable"
	ResponseIndeterminate   = "Indeterminate"
	ResponseIndeterminateD  = "Indeterminate{D}"
	ResponseIndeterminateP  = "Indeterminate{P}"
	ResponseIndeterminateDP = "Indeterminate{DP}"
)

func GetResultFromResponse(responseXML string) (string, error) {
	mxj.IncludeTagSeqNum(true)
	responseMap, err := mxj.NewMapXml([]byte(responseXML))
	if err != nil {
		return "", err
	}
	response := Response{responseMap}
	result, _ := response.GetResult()
	return result, nil
}

func (response Response) AddResult(result, status string) {
	if response.response == nil {
		response.response = make(map[string]interface{})
	}
	r := make(map[string]interface{})
	d := make(map[string]interface{})
	d["#text"] = result
	s := make(map[string]interface{})
	s["#text"] = status
	r["Decision"] = d
	r["Status"] = s

	responseTag, ok := response.response["Response"].(map[string]interface{})
	if !ok {
		response.response["Response"] = make(map[string]interface{})
		responseTag, _ = response.response["Response"].(map[string]interface{})
	}
	switch responseTag["Result"].(type) {
	case []interface{}:
		results, _ := responseTag["Result"].([]interface{})
		results = append(results, r)
	case map[string]interface{}:
		prevResult := responseTag["Result"]
		results := make([]interface{}, 0)
		results = append(results, prevResult)
		results = append(results, r)
		responseTag["Result"] = results
	default:
		responseTag["Result"] = r
	}
}

func (response Response) GetResult() (string, string) {
	responseTag, ok := response.response["Response"].(map[string]interface{})
	if !ok {
		return ResponseIndeterminate, "The response seems to be malformatted"
	}
	switch responseTag["Result"].(type) {
	case []interface{}:
		//TODO Use some algorithm to determine which one to return
		return ResponseIndeterminate, "At this time multiple result responses are not supported"
	case map[string]interface{}:
		resultTag, ok := responseTag["Result"].(map[string]interface{})
		if !ok {
			return ResponseIndeterminate, "The Result seems to be malformatted"
		}
		var decision, status string
		decisionTag, ok := resultTag["Decision"].(map[string]interface{})
		if !ok {
			decision, ok = resultTag["Decision"].(string)
			if !ok {
				return ResponseIndeterminate, "The Decision tag seems to be malformatted"
			}
		} else {
			decision, _ = decisionTag["#text"].(string)
		}
		statusTag, _ := resultTag["Status"].(map[string]interface{})
		if !ok {
			status, _ = resultTag["Status"].(string)
		} else {
			status, _ = statusTag["#text"].(string)
		}
		return decision, status
	default:
		return ResponseIndeterminate, "The Result seems to be malformatted"
	}

	return ResponseIndeterminate, "There was no applicable Result"
}

func (response Response) Xml() string {
	responseInterface, _ := response.response["Response"].(map[string]interface{})
	if responseInterface == nil {
		responseInterface = make(map[string]interface{})
		response.response["Response"] = responseInterface
	}
	responseInterface["-xmlns"] = "urn:oasis:names:tc:xacml:3.0:core:schema:wd-17"
	responseInterface["-xmlns:xsi"] = "http://www.w3.org/2001/XMLSchema-instance"
	responseInterface["-xsi:schemaLocation"] = "urn:oasis:names:tc:xacml:3.0:core:schema:wd-17 http://docs.oasis-open.org/xacml/3.0/xacml-core-v3-schema-wd-17.xsd"
	m := mxj.Map(response.response)
	ret, _ := m.Xml()
	return `<?xml version="1.0" encoding="UTF-8"?>` + string(ret)
}
