package xacml

import (
	"log"
)

const (
	ResultPermit        = "Permit"
	ResultDeny          = "Deny"
	ResultNotApplicable = "NotApplicable"
	ResultIndeterminate = "Indeterminate"
)

type Policy struct {
	policy map[string]interface{}
}

func evaluatePolicy(policy map[string]interface{}, request map[string]interface{}) Response {
	p := Policy{policy}
	r := Request{request}
	response := p.Evaluate(r)
	return response
}

func (policy Policy) Evaluate(request Request) Response {
	policyBody, _ := policy.policy["Policy"].(map[string]interface{})
	response := Response{make(map[string]interface{})}

	//Evaluate Target
	targetResponse := evaluateTargetBody(policyBody["Target"], request)
	//log.Println(targetResponse)
	if targetResponse != targetMatch {
		if targetResponse == targetNoMatch {
			response.AddResult(ResponseNotApplicable, "Policy Target returned No Match")
			return response
		}
		if targetResponse == targetIndeterminate {
			response.AddResult(ResponseIndeterminate, "Policy Target returned Indeterminate")
			return response
		}
	}

	rules := make([]evaluatable, 0)
	switch policyBody["Rule"].(type) {
	case map[string]interface{}:
		rules = append(rules, Rule{policyBody["Rule"].(map[string]interface{})})
	case []interface{}:
		for _, ruleInterface := range policyBody["Rule"].([]interface{}) {
			ruleMap, _ := ruleInterface.(map[string]interface{})
			rule := Rule{ruleMap}
			rules = append(rules, rule)
		}
	default:
		log.Println("Invalid Rule, was neither a tag or array")
		response.AddResult(ResponseIndeterminate, "Invalid Rule, was neither a tag or array")
		return response
	}

	ruleCombiningAlgorithm, _ := policyBody["-RuleCombiningAlgId"].(string)
	//Evaluate Rules
	response = combineAlgorithm(ruleCombiningAlgorithm, rules, request)
	//log.Println(response)
	return response
}
