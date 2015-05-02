package xacml

import (
	"github.com/clbanning/mxj"
	"testing"
)

var iamPolicy1 = `<?xml version="1.0" encoding="UTF-8"?>
<Policy xmlns="urn:oasis:names:tc:xacml:3.0:core:schema:wd-17" 
		xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
		xsi:schemaLocation="urn:oasis:names:tc:xacml:3.0:core:schema:wd-17 http://docs.oasis-open.org/xacml/3.0/xacml-core-v3-schema-wd-17.xsd" 
		PolicyId="httpbin-policy" 
		Version="1.0" 
		RuleCombiningAlgId="urn:oasis:names:tc:xacml:3.0:rule-combining-algorithm:deny-overrides">
	<Description>IAM Policy</Description>
	<Target>
		<AnyOf>
			<AllOf>
				<Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-regexp-match">
					<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">^https://api.emc.blue/auth/.*</AttributeValue>
					<AttributeDesignator
							MustBePresent="false"
							Category="urn:oasis:names:tc:xacml:3.0:attribute-category:resource"
							AttributeId="urn:oasis:names:tc:xacml:1.0:resource:resource-id" 
							DataType="http://www.w3.org/2001/XMLSchema#string"/>
				</Match>
			</AllOf>
		</AnyOf>
	</Target>
	<Rule RuleId="denySalesOps" Effect="Deny">
		<Description>This rule will look for the user to be in a specific tenant and will deny the tenant admin role</Description>
		<Target>
			<AnyOf>
				<AllOf>
					<Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
						<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">https://api.emc.blue/auth/role.tenant_admin</AttributeValue>
						<AttributeDesignator 
							MustBePresent="false"
							Category="urn:oasis:names:tc:xacml:3.0:attribute-category:resource"
							AttributeId="urn:oasis:names:tc:xacml:1.0:resource:resource-id"
							DataType="http://www.w3.org/2001/XMLSchema#string"/>
					</Match>
					<Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
						<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">0a56f466-0af7-4521-842b-1f3577a1f0de</AttributeValue>
						<AttributeDesignator 
							MustBePresent="false"
							Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject"
							AttributeId="tenant_id"
							DataType="http://www.w3.org/2001/XMLSchema#string"/>
					</Match>
				</AllOf>
			</AnyOf>
		</Target>
	</Rule>
	<Rule RuleId="permitTenantAdmin" Effect="Permit">
		<Description>This rule will look for the user to be in any tenant</Description>
		<Target>
			<AnyOf>
				<AllOf>
					<Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
						<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">https://api.emc.blue/auth/role.tenant_admin</AttributeValue>
						<AttributeDesignator 
							MustBePresent="false"
							Category="urn:oasis:names:tc:xacml:3.0:attribute-category:resource"
							AttributeId="urn:oasis:names:tc:xacml:1.0:resource:resource-id"
							DataType="http://www.w3.org/2001/XMLSchema#string"/>
					</Match>
					<Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-regexp-match">
						<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}</AttributeValue>
						<AttributeDesignator 
							MustBePresent="false"
							Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject"
							AttributeId="tenant_id"
							DataType="http://www.w3.org/2001/XMLSchema#string"/>
					</Match>
				</AllOf>
			</AnyOf>
		</Target>
	</Rule>
	<Rule RuleId="permitSalesOps" Effect="Permit">
		<Description>This rule will look for the user to be in a specific tenant</Description>
		<Target>
			<AnyOf>
				<AllOf>
					<Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
						<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">https://api.emc.blue/auth/role.salesops</AttributeValue>
						<AttributeDesignator 
							MustBePresent="false"
							Category="urn:oasis:names:tc:xacml:3.0:attribute-category:resource"
							AttributeId="urn:oasis:names:tc:xacml:1.0:resource:resource-id"
							DataType="http://www.w3.org/2001/XMLSchema#string"/>
					</Match>
					<Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
						<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">0a56f466-0af7-4521-842b-1f3577a1f0de</AttributeValue>
						<AttributeDesignator 
							MustBePresent="false"
							Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject"
							AttributeId="tenant_id"
							DataType="http://www.w3.org/2001/XMLSchema#string"/>
					</Match>
				</AllOf>
			</AnyOf>
		</Target>
	</Rule>
	<Rule RuleId="permitTrustedAgent" Effect="Permit">
		<Description>This rule will look for the user/agent to be in a specific tenant</Description>
		<Target>
			<AnyOf>
				<AllOf>
					<Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
						<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">https://api.emc.blue/auth/role.trusted_service</AttributeValue>
						<AttributeDesignator 
							MustBePresent="false"
							Category="urn:oasis:names:tc:xacml:3.0:attribute-category:resource"
							AttributeId="urn:oasis:names:tc:xacml:1.0:resource:resource-id"
							DataType="http://www.w3.org/2001/XMLSchema#string"/>
					</Match>
					<Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
						<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">85b281cf-c074-4d08-80d5-4fd98458640f</AttributeValue>
						<AttributeDesignator 
							MustBePresent="false"
							Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject"
							AttributeId="tenant_id"
							DataType="http://www.w3.org/2001/XMLSchema#string"/>
					</Match>
				</AllOf>
			</AnyOf>
		</Target>
	</Rule>
	<Rule RuleId="permitAPI" Effect="Permit">
		<Description>This rule will return a permit on any of the api family asks.</Description>
		<Target/>
		<Condition>
			<Apply FunctionId="urn:oasis:names:tc:xacml:3.0:function:string-is-in">
				<AttributeDesignator 
						MustBePresent="false"
						Category="urn:oasis:names:tc:xacml:3.0:attribute-category:resource"
						AttributeId="urn:oasis:names:tc:xacml:1.0:resource:resource-id"
						DataType="http://www.w3.org/2001/XMLSchema#string"/>
				<Apply FunctionId="urn:oasis:names:tc:xacml:3.0:function:string-bag">
					<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">https://api.emc.blue/auth/iam</AttributeValue>
					<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">https://api.emc.blue/auth/config</AttributeValue>
					<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">https://api.emc.blue/auth/agent</AttributeValue>
					<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">https://api.emc.blue/auth/storage</AttributeValue>
					<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">https://api.emc.blue/auth/event</AttributeValue>
					<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">https://api.emc.blue/auth/logging</AttributeValue>
					<AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">https://api.emc.blue/auth/provision</AttributeValue>
				</Apply>
			</Apply>
		</Condition>
	</Rule>
</Policy>`

func TestIAMTenantAdmin(t *testing.T) {
	request := NewRequest()
	request.AddAttribute(AttributeCategoryResource, IdentifierResourceId, DataTypeString, "https://api.emc.blue/auth/role.tenant_admin")
	request.AddAttribute(AttributeCategorySubjectAccessSubject, "tenant_id", DataTypeString, "0a56f466-0afa-4521-842b-1f3577a1f0da")

	mxj.IncludeTagSeqNum(true)
	responseXML, err := PolicyDecisionPoint(iamPolicy1, request.Xml())
	if err != nil {
		t.Errorf("Error evaluating policy", err)
	}
	r, err := GetResultFromResponse(responseXML)
	if r != ResponsePermit {
		t.Errorf("Result: %v", r)
	}
}

func TestIAMSalesOps(t *testing.T) {
	request := NewRequest()
	request.AddAttribute(AttributeCategoryResource, IdentifierResourceId, DataTypeString, "https://api.emc.blue/auth/role.salesops")
	request.AddAttribute(AttributeCategorySubjectAccessSubject, "tenant_id", DataTypeString, "0a56f466-0af7-4521-842b-1f3577a1f0de")

	mxj.IncludeTagSeqNum(true)
	responseXML, err := PolicyDecisionPoint(iamPolicy1, request.Xml())
	if err != nil {
		t.Errorf("Error evaluating policy", err)
	}
	r, err := GetResultFromResponse(responseXML)
	if r != ResponsePermit {
		t.Errorf("Result: %v", r)
	}
}

func TestIAMSalesOpsNoTenantAdmin(t *testing.T) {
	request := NewRequest()
	request.AddAttribute(AttributeCategoryResource, IdentifierResourceId, DataTypeString, "https://api.emc.blue/auth/role.tenant_admin")
	request.AddAttribute(AttributeCategorySubjectAccessSubject, "tenant_id", DataTypeString, "0a56f466-0af7-4521-842b-1f3577a1f0de")

	mxj.IncludeTagSeqNum(true)
	responseXML, err := PolicyDecisionPoint(iamPolicy1, request.Xml())
	if err != nil {
		t.Errorf("Error evaluating policy", err)
	}
	r, err := GetResultFromResponse(responseXML)
	if r != ResponseDeny {
		t.Errorf("Result: %v", r)
	}
}

func TestIAMTrustedAgent(t *testing.T) {
	request := NewRequest()
	request.AddAttribute(AttributeCategoryResource, IdentifierResourceId, DataTypeString, "https://api.emc.blue/auth/role.trusted_service")
	request.AddAttribute(AttributeCategorySubjectAccessSubject, "tenant_id", DataTypeString, "85b281cf-c074-4d08-80d5-4fd98458640f")

	mxj.IncludeTagSeqNum(true)
	responseXML, err := PolicyDecisionPoint(iamPolicy1, request.Xml())
	if err != nil {
		t.Errorf("Error evaluating policy", err)
	}
	r, err := GetResultFromResponse(responseXML)
	if r != ResponsePermit {
		t.Errorf("Result: %v", r)
	}
}

func TestIAMAPIFamily(t *testing.T) {
	request := NewRequest()
	request.AddAttribute(AttributeCategoryResource, IdentifierResourceId, DataTypeString, "https://api.emc.blue/auth/iam")

	mxj.IncludeTagSeqNum(true)
	responseXML, err := PolicyDecisionPoint(iamPolicy1, request.Xml())
	if err != nil {
		t.Errorf("Error evaluating policy", err)
	}
	r, err := GetResultFromResponse(responseXML)
	if r != ResponsePermit {
		t.Errorf("Result: %v", r)
	}
}

func TestIAMOther(t *testing.T) {
	request := NewRequest()
	request.AddAttribute(AttributeCategoryResource, IdentifierResourceId, DataTypeString, "https://api.emc.blue/auth/fina")

	mxj.IncludeTagSeqNum(true)
	responseXML, err := PolicyDecisionPoint(iamPolicy1, request.Xml())
	if err != nil {
		t.Errorf("Error evaluating policy", err)
	}
	r, err := GetResultFromResponse(responseXML)
	if r != ResponseNotApplicable {
		t.Errorf("Result: %v", r)
	}
}
