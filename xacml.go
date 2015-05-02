package xacml

import (
	"github.com/clbanning/mxj"
)

const (
	DataTypeString            = "http://www.w3.org/2001/XMLSchema#string"
	DataTypeBoolean           = "http://www.w3.org/2001/XMLSchema#boolean"
	DataTypeInteger           = "http://www.w3.org/2001/XMLSchema#integer"
	DataTypeDouble            = "http://www.w3.org/2001/XMLSchema#double"
	DataTypeTime              = "http://www.w3.org/2001/XMLSchema#time"
	DataTypeDate              = "http://www.w3.org/2001/XMLSchema#date"
	DataTypeDateTime          = "http://www.w3.org/2001/XMLSchema#dateTime"
	DataTypeDayTimeDuration   = "http://www.w3.org/2001/XMLSchema#dayTimeDuration"
	DataTypeYearMonthDuration = "http://www.w3.org/2001/XMLSchema#yearMonthDuration"
	DataTypeAnyURI            = "http://www.w3.org/2001/XMLSchema#anyURI"
	DataTypeHexBinary         = "http://www.w3.org/2001/XMLSchema#hexBinary"
	DataTypeBase64Binary      = "http://www.w3.org/2001/XMLSchema#base64Binary"
	DataTypeRFC822Name        = "urn:oasis:names:tc:xacml:1.0:data-type:rfc822Name"
	DataTypeX500Name          = "urn:oasis:names:tc:xacml:1.0:data-type:x500Name"
	DataTypeXPathExpression   = "urn:oasis:names:tc:xacml:3.0:data-type:xpathExpression"
	DataTypeIPAddress         = "urn:oasis:names:tc:xacml:2.0:data-type:ipAddress"
	DataTypeDNSName           = "urn:oasis:names:tc:xacml:2.0:data-type:dnsName"
)

type evaluatable interface {
	Evaluate(request Request) Response
}

func init(){
	mxj.IncludeTagSeqNum(true)
}

func PolicyDecisionPoint(policy string, request string) (string, error) {
	mxj.IncludeTagSeqNum(true)
	policyMap, err := mxj.NewMapXml([]byte(policy))
	//b, _ := json.Marshal(policyMap)
	//log.Println(string(b))
	if err != nil {
		return "", err
	}
	requestMap, err := mxj.NewMapXml([]byte(request))
	if err != nil {
		return "", err
	}
	response := policyDecisionPoint(policyMap, requestMap)
	xml := response.Xml()
	return xml, nil
}

func policyDecisionPoint(policy map[string]interface{}, request map[string]interface{}) Response {
	return evaluatePolicy(policy, request)
}
