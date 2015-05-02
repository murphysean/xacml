package xacml

import (
	"github.com/clbanning/mxj"
	"testing"
)

var conditionReq = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<Request>
	<Attributes Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject">
		<Attribute AttributeId="urn:oasis:names:tc:xacml:1.0:subject:subject-id">
			<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">Sean Murphy</AttributeValue>
		</Attribute>
		<Attribute AttributeId="urn:oasis:names:tc:xacml:3.0:example:attribute:role">
			<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">superhero</AttributeValue>
		</Attribute>
	</Attributes>
	<Attributes Category="urn:oasis:names:tc:xacml:3.0:attribute-category:resource">
		<Attribute AttributeId="urn:oasis:names:tc:xacml:2.0:resource:target-namespace" >
			<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#anyURI">urn:ms:xacml:record</AttributeValue>
		</Attribute>
		<Attribute AttributeId="urn:oasis:names:tc:xacml:1.0:resource:resource-id">
			<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">http://localhost:8280/services/Customers/getEmployees</AttributeValue>
		</Attribute>
		<Attribute AttributeId="urn:oasis:names:tc:xacml:1.0:resource:resource-owner">
			<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">Sean Murphy</AttributeValue>
		</Attribute>
	</Attributes>
	<Attributes Category="urn:oasis:names:tc:xacml:3.0:attribute-category:action">
		<Attribute AttributeId="urn:oasis:names:tc:xacml:1.0:action:action-id">
			<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">GET</AttributeValue>
		</Attribute>
	</Attributes>
	<Attributes Category="urn:oasis:names:tc:xacml:3.0:attribute-category:environment">
		<Attribute AttributeId="urn:oasis:names:tc:xacml:1.0:environment:current-date" >
			<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#date">2010-01-11</AttributeValue>
		</Attribute>
	</Attributes>
</Request>`)
var goodSimpleCondition = []byte(`<Condition>
	<Apply FunctionId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
		<AttributeDesignator MustBePresent="false" 
			Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject" 
			AttributeId="urn:oasis:names:tc:xacml:1.0:subject:subject-id" 
			DataType="http://www.w3.org/2001/XMLSchema#string"/>
		<AttributeDesignator MustBePresent="false" 
			Category="urn:oasis:names:tc:xacml:3.0:attribute-category:resource" 
			AttributeId="urn:oasis:names:tc:xacml:1.0:resource:resource-owner" 
			DataType="http://www.w3.org/2001/XMLSchema#string"/>
	</Apply>
</Condition>`)
var badSimpleCondition = []byte(``)

func TestGoodSimpleCondition(t *testing.T) {
	mxj.IncludeTagSeqNum(true)
	//Create request
	m, _ := mxj.NewMapXml(conditionReq)
	request := Request{m}
	//b, _ := json.Marshal(request.request)
	//log.Println(string(b))

	//Create Condition
	condition, _ := mxj.NewMapXml(goodSimpleCondition)
	//log.Println(target)

	result := evaluateCondition(condition, request)

	if result != "True" {
		t.Errorf("Result: %v", result)
	}
}

func TestRequestCreation(t *testing.T) {
	mxj.IncludeTagSeqNum(true)
	request := Request{make(map[string]interface{})}
	request.AddAttribute(AttributeCategorySubjectAccessSubject, IdentifierSubjectId, DataTypeString, "11203456")
	request.AddAttribute(AttributeCategoryAction, IdentifierActionId, DataTypeString, "GET")
	xml := request.Xml()

	if xml == "" {
		t.Errorf("Empty Result")
	}
}
