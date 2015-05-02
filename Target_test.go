package xacml

import (
	"github.com/clbanning/mxj"
	"testing"
)

var requestBytes = []byte(`<?xml version="1.0" encoding="UTF-8"?>
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
var matchBytes = []byte(`<Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
						<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">Sean Murphy</AttributeValue>
						<AttributeDesignator MustBePresent="false" Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject" AttributeId="urn:oasis:names:tc:xacml:1.0:subject:subject-id" DataType="http://www.w3.org/2001/XMLSchema#string"/>
					</Match>`)
var targetBytes = []byte(`<Target>
         <AnyOf>
            <AllOf>
               <Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
                  <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">http://localhost:8280/services/Customers/getEmployees</AttributeValue>
                  <AttributeDesignator AttributeId="urn:oasis:names:tc:xacml:1.0:resource:resource-id" Category="urn:oasis:names:tc:xacml:3.0:attribute-category:resource" DataType="http://www.w3.org/2001/XMLSchema#string" MustBePresent="true"></AttributeDesignator>
               </Match>
               <Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
                  <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">GET</AttributeValue>
                  <AttributeDesignator AttributeId="urn:oasis:names:tc:xacml:1.0:action:action-id" Category="urn:oasis:names:tc:xacml:3.0:attribute-category:action" DataType="http://www.w3.org/2001/XMLSchema#string" MustBePresent="true"></AttributeDesignator>
               </Match>
            </AllOf>
			<AllOf>
               <Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
                  <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">http://localhost:8280/services/Customers/getEmployee</AttributeValue>
                  <AttributeDesignator AttributeId="urn:oasis:names:tc:xacml:1.0:resource:resource-id" Category="urn:oasis:names:tc:xacml:3.0:attribute-category:resource" DataType="http://www.w3.org/2001/XMLSchema#string" MustBePresent="true"></AttributeDesignator>
               </Match>
               <Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
                  <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">GET</AttributeValue>
                  <AttributeDesignator AttributeId="urn:oasis:names:tc:xacml:1.0:action:action-id" Category="urn:oasis:names:tc:xacml:3.0:attribute-category:action" DataType="http://www.w3.org/2001/XMLSchema#string" MustBePresent="true"></AttributeDesignator>
               </Match>
            </AllOf>
         </AnyOf>
		<AnyOf>
            <AllOf>
               <Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
                  <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">http://localhost:8280/services/Customers/getEmployees</AttributeValue>
                  <AttributeDesignator AttributeId="urn:oasis:names:tc:xacml:1.0:resource:resource-id" Category="urn:oasis:names:tc:xacml:3.0:attribute-category:resource" DataType="http://www.w3.org/2001/XMLSchema#string" MustBePresent="true"></AttributeDesignator>
               </Match>
               <Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
                  <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">GET</AttributeValue>
                  <AttributeDesignator AttributeId="urn:oasis:names:tc:xacml:1.0:action:action-id" Category="urn:oasis:names:tc:xacml:3.0:attribute-category:action" DataType="http://www.w3.org/2001/XMLSchema#string" MustBePresent="true"></AttributeDesignator>
               </Match>
            </AllOf>
		</AnyOf>
      </Target>`)

func TestEvaluateTarget(t *testing.T) {
	mxj.IncludeTagSeqNum(true)
	//Create request
	m, _ := mxj.NewMapXml(requestBytes)
	request := Request{m}

	//Create Target
	target, _ := mxj.NewMapXml(targetBytes)
	//log.Println(target)

	result := evaluateTarget(target, request)

	if result != "Match" {
		t.Errorf("Result: %v", result)
	}
}

func TestEvaluateMatch(t *testing.T) {
	mxj.IncludeTagSeqNum(true)
	//TODO Create request
	m, _ := mxj.NewMapXml(requestBytes)
	request := Request{m}

	//TODO Create Match
	match, _ := mxj.NewMapXml(matchBytes)

	result := evaluateMatch(match, request)

	if result != targetMatch {
		t.Errorf("Result: %v", result)
	}
}
