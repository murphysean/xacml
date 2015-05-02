package xacml

import ()

type Rule struct {
	rule map[string]interface{}
}

func (rule Rule) Evaluate(request Request) Response {
	ruleBody := rule.rule
	response := Response{make(map[string]interface{})}

	//Evaluate target
	targetResponse := evaluateTargetBody(ruleBody["Target"], request)
	if targetResponse == targetNoMatch {
		//log.Println("Rule Target didn't match")
		response.AddResult(ResponseNotApplicable, "Rule Target didn't match")
		return response
	}
	if targetResponse == targetIndeterminate {
		//log.Println("Error evaluating match")
		response.AddResult(ResponseIndeterminate, "Error evaluating target")
		return response
	}

	//Evaluate condition
	conditionResponse := evaluateConditionBody(ruleBody["Condition"], request)
	switch conditionResponse {
	case ConditionTrue:
		effect, _ := rule.rule["-Effect"].(string)
		response.AddResult(effect, "")
		return response
	case ConditionFalse:
		response.AddResult(ResponseNotApplicable, "Condition returned false")
		return response
	case ConditionIndeterminate:
		response.AddResult(ResponseIndeterminate, "Error evaluating Condition")
		return response
	}

	response.AddResult(ResponseIndeterminate, "No other case caught")
	return response
}
