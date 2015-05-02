package xacml

import (
	"github.com/clbanning/mxj"
	"testing"
)

var requestString = `<?xml version="1.0" encoding="UTF-8"?>
<Request xmlns="urn:oasis:names:tc:xacml:3.0:core:schema:wd-17" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:oasis:names:tc:xacml:3.0:core:schema:wd-17 http://docs.oasis-open.org/xacml/3.0/xacml-core-v3-schema-wd-17.xsd">
<Attributes Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject">
    <Attribute AttributeId="urn:oasis:names:tc:xacml:1.0:subject:subject-id">
        <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">Sean Murphy</AttributeValue>
    </Attribute>
	<Attribute AttributeId="group">
		<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">group1</AttributeValue>
		<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">group2</AttributeValue>
		<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">group3</AttributeValue>
	</Attribute>
    <Attribute AttributeId="ip-address">
        <AttributeValue DataType="urn:oasis:names:tc:xacml:2.0:data-type:ipAddress">127.0.0.1:50563</AttributeValue>
    </Attribute>
    <Attribute AttributeId="user-agent">
        <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:33.0) Gecko/20100101 Firefox/33.0</AttributeValue>
    </Attribute>
    <Attribute AttributeId="referer">
        <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string"/>
    </Attribute>
</Attributes>
<Attributes Category="urn:oasis:names:tc:xacml:1.0:subject-category:intermediary-subject">
    <Attribute AttributeId="urn:oasis:names:tc:xacml:1.0:subject:subject-id">
        <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">localhost:8443</AttributeValue>
    </Attribute>
</Attributes>
<Attributes Category="urn:oasis:names:tc:xacml:3.0:attribute-category:resource">
    <Attribute AttributeId="urn:oasis:names:tc:xacml:1.0:resource:resource-id">
        <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">/httpbin/</AttributeValue>
    </Attribute>
    <Attribute AttributeId="urn:oasis:names:tc:xacml:1.0:resource:resource-location">
        <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#anyURI">/httpbin/</AttributeValue>
    </Attribute>
    <Attribute AttributeId="urn:oasis:names:tc:xacml:1.0:resource:simple-file-name">
        <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">/httpbin/</AttributeValue>
    </Attribute>
</Attributes>
<Attributes Category="urn:oasis:names:tc:xacml:3.0:attribute-category:action">
    <Attribute AttributeId="urn:oasis:names:tc:xacml:1.0:action:action-id">
        <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">GET</AttributeValue>
    </Attribute>
    <Attribute AttributeId="protocol">
        <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">https</AttributeValue>
    </Attribute>
    <Attribute AttributeId="secure">
        <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#boolean">true</AttributeValue>
    </Attribute>
    <Attribute AttributeId="cookie">
        <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#boolean">true</AttributeValue>
    </Attribute>
    <Attribute AttributeId="basic">
        <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#boolean">false</AttributeValue>
    </Attribute>
    <Attribute AttributeId="bearer">
        <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#boolean">false</AttributeValue>
    </Attribute>
</Attributes>
</Request>`

var policyString = `<?xml version="1.0" encoding="UTF-8"?>
<Policy xmlns="urn:oasis:names:tc:xacml:3.0:core:schema:wd-17" 
		xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
		xsi:schemaLocation="urn:oasis:names:tc:xacml:3.0:core:schema:wd-17 http://docs.oasis-open.org/xacml/3.0/xacml-core-v3-schema-wd-17.xsd" 
		PolicyId="httpbin-policy" 
		Version="1.0" 
		RuleCombiningAlgId="urn:oasis:names:tc:xacml:3.0:rule-combining-algorithm:deny-overrides">
	<Description>HTTP Bin Policy</Description>
	<Target>
		<AnyOf>
			<AllOf>
				<Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-regexp-match">
					<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">^/httpbin/</AttributeValue>
					<AttributeDesignator
							MustBePresent="false"
							Category="urn:oasis:names:tc:xacml:3.0:attribute-category:resource"
							AttributeId="urn:oasis:names:tc:xacml:1.0:resource:resource-id" 
							DataType="http://www.w3.org/2001/XMLSchema#string"/>
				</Match>
			</AllOf>
		</AnyOf>
	</Target>
	<Rule RuleId="denynonsecure" Effect="Deny">
		<Description></Description>
		<Target>
			<AnyOf>
				<AllOf>
					<Match MatchId="urn:oasis:names:tc:xacml:1.0:function:boolean-equal">
						<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#boolean">false</AttributeValue>
						<AttributeDesignator 
							MustBePresent="false"
							Category="urn:oasis:names:tc:xacml:3.0:attribute-category:action"
							AttributeId="secure"
							DataType="http://www.w3.org/2001/XMLSchema#boolean"/>
					</Match>
				</AllOf>
			</AnyOf>
		</Target>
	</Rule>
</Policy>`

var policyString2 = `<?xml version="1.0" encoding="UTF-8"?>
<Policy xmlns="urn:oasis:names:tc:xacml:3.0:core:schema:wd-17" 
		xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
		xsi:schemaLocation="urn:oasis:names:tc:xacml:3.0:core:schema:wd-17 http://docs.oasis-open.org/xacml/3.0/xacml-core-v3-schema-wd-17.xsd" 
		PolicyId="httpbin-policy" 
		Version="1.0" 
		RuleCombiningAlgId="urn:oasis:names:tc:xacml:3.0:rule-combining-algorithm:deny-overrides">
	<Description>HTTP Bin Policy</Description>
	<Target>
		<AnyOf>
			<AllOf>
				<Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-regexp-match">
					<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">^/httpbin/</AttributeValue>
					<AttributeDesignator
							MustBePresent="false"
							Category="urn:oasis:names:tc:xacml:3.0:attribute-category:resource"
							AttributeId="urn:oasis:names:tc:xacml:1.0:resource:resource-id" 
							DataType="http://www.w3.org/2001/XMLSchema#string"/>
				</Match>
			</AllOf>
		</AnyOf>
	</Target>
	<Rule RuleId="denynonsecure" Effect="Deny">
		<Description></Description>
		<Target>
			<AnyOf>
				<AllOf>
					<Match MatchId="urn:oasis:names:tc:xacml:1.0:function:boolean-equal">
						<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#boolean">false</AttributeValue>
						<AttributeDesignator 
							MustBePresent="false"
							Category="urn:oasis:names:tc:xacml:3.0:attribute-category:action"
							AttributeId="secure"
							DataType="http://www.w3.org/2001/XMLSchema#boolean"/>
					</Match>
				</AllOf>
			</AnyOf>
		</Target>
	</Rule>
	<Rule RuleId="denygroup1" Effect="Deny">
		<Description></Description>
		<Target>
			<AnyOf>
				<AllOf>
					<Match MatchId="urn:oasis:names:tc:xacml:3.0:function:string-is-in">
						<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">group1</AttributeValue>
						<AttributeDesignator
								MustBePresent="false"
								Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject"
								AttributeId="group" 
								DataType="http://www.w3.org/2001/XMLSchema#string"/>
					</Match>
				</AllOf>
			</AnyOf>
		</Target>
	</Rule>
	<Rule RuleId="permitnamedpersons" Effect="Permit">
		<Description></Description>
		<Target/>
		<Condition>
			<Apply FunctionId="urn:oasis:names:tc:xacml:3.0:function:string-is-in">
				<AttributeDesignator 
						MustBePresent="false"
						Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject"
						AttributeId="urn:oasis:names:tc:xacml:1.0:subject:subject-id"
						DataType="http://www.w3.org/2001/XMLSchema#string"/>
				<Apply FunctionId="urn:oasis:names:tc:xacml:3.0:function:string-bag">
					<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">Sean Murphy</AttributeValue>
					<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">Eric Murphy</AttributeValue>
					<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">Ryan Murphy</AttributeValue>
				</Apply>
			</Apply>
		</Condition>
	</Rule>
</Policy>`

func TestPolicyDecisionPoint(t *testing.T) {
	mxj.IncludeTagSeqNum(true)
	responseXML, err := PolicyDecisionPoint(policyString2, requestString)
	if err != nil {
		t.Errorf("Error evaluating policy", err)
	}
	r, err := GetResultFromResponse(responseXML)
	if r != ResponseDeny {
		t.Errorf("Result: %v", r)
	}
}
