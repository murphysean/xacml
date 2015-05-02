package xacml

import (
	"log"
)

const (
	PolicyCombiningAlgorithmDenyOverrides          = "urn:oasis:names:tc:xacml:3.0:policy-combining-algorithm:deny-overrides"
	PolicyCombiningAlgorithmPermitOverrides        = "urn:oasis:names:tc:xacml:3.0:policy-combining-algorithm:permit-overrides"
	PolicyCombiningAlgorithmOrderedDenyOverrides   = "urn:oasis:names:tc:xacml:3.0:policy-combining-algorithm:ordered-deny-overrides"
	PolicyCombiningAlgorithmOrderedPermitOverrides = "urn:oasis:names:tc:xacml:3.0:policy-combining-algorithm:ordered-permit-overrides"
	PolicyCombiningAlgorithmDenyUnlessPermit       = "urn:oasis:names:tc:xacml:3.0:policy-combining-algorithm:deny-unless-permit"
	PolicyCombiningAlgorithmPermitUnlessDeny       = "urn:oasis:names:tc:xacml:3.0:policy-combining-algorithm:permit-unless-deny"

	PolicyCombiningAlgorithmFirstApplicable   = "urn:oasis:names:tc:xacml:1.0:policy-combining-algorithm:first-applicable"
	PolicyCombiningAlgorithmOnlyOneApplicable = "urn:oasis:names:tc:xacml:1.0:policy-combining-algorithm:only-one-applicable"

	RuleCombiningAlgorithmDenyOverrides          = "urn:oasis:names:tc:xacml:3.0:rule-combining-algorithm:deny-overrides"
	RuleCombiningAlgorithmPermitOverrides        = "urn:oasis:names:tc:xacml:3.0:rule-combining-algorithm:permit-overrides"
	RuleCombiningAlgorithmOrderedDenyOverrides   = "urn:oasis:names:tc:xacml:3.0:rule-combining-algorithm:ordered-deny-overrides"
	RuleCombiningAlgorithmOrderedPermitOverrides = "urn:oasis:names:tc:xacml:3.0:rule-combining-algorithm:ordered-permit-overrides"
	RuleCombiningAlgorithmDenyUnlessPermit       = "urn:oasis:names:tc:xacml:3.0:rule-combining-algorithm:deny-unless-permit"
	RuleCombiningAlgorithmPermitUnlessDeny       = "urn:oasis:names:tc:xacml:3.0:rule-combining-algorithm:permit-unless-deny"

	RuleCombiningAlgorithmFirstApplicable = "urn:oasis:names:tc:xacml:1.0:rule-combining-algorithm:first-applicable"
)

func combineAlgorithm(combiningAlgorithm string, evaluatables []evaluatable, request Request) Response {
	ret := Response{make(map[string]interface{})}
	switch combiningAlgorithm {
	case PolicyCombiningAlgorithmDenyOverrides:
		fallthrough
	case PolicyCombiningAlgorithmOrderedDenyOverrides:
		fallthrough
	case RuleCombiningAlgorithmDenyOverrides:
		fallthrough
	case RuleCombiningAlgorithmOrderedDenyOverrides:
		ret = denyOverridesCombiningAlgorithm(evaluatables, request)

	case PolicyCombiningAlgorithmPermitOverrides:
		fallthrough
	case PolicyCombiningAlgorithmOrderedPermitOverrides:
		fallthrough
	case RuleCombiningAlgorithmPermitOverrides:
		fallthrough
	case RuleCombiningAlgorithmOrderedPermitOverrides:
		ret = permitOverridesCombiningAlgorithm(evaluatables, request)

	case PolicyCombiningAlgorithmDenyUnlessPermit:
		fallthrough
	case RuleCombiningAlgorithmDenyUnlessPermit:
		ret = denyUnlessPermitCombiningAlgorithm(evaluatables, request)

	case PolicyCombiningAlgorithmPermitUnlessDeny:
		fallthrough
	case RuleCombiningAlgorithmPermitUnlessDeny:
		ret = permitUnlessDenyCombiningAlgorithm(evaluatables, request)

	default:
		log.Println("Unsupported Combining Algorithm")
		ret.AddResult(ResponseIndeterminate, "Unsupported Combining Algorithm")
	}
	return ret
}

func denyOverridesCombiningAlgorithm(evaluatables []evaluatable, request Request) Response {
	ret := Response{make(map[string]interface{})}
	atLeastOneErrorD := false
	atLeastOneErrorP := false
	atLeastOneErrorDP := false
	atLeastOnePermit := false
	for _, e := range evaluatables {
		response := e.Evaluate(request)
		decision, status := response.GetResult()
		if decision == ResponseDeny {
			ret.AddResult(decision, status)
			return ret
		}
		if decision == ResponsePermit {
			atLeastOnePermit = true
			continue
		}
		if decision == ResponseNotApplicable {
			continue
		}
		if decision == ResponseIndeterminateD {
			atLeastOneErrorD = true
			continue
		}
		if decision == ResponseIndeterminateP {
			atLeastOneErrorP = true
			continue
		}
		if decision == ResponseIndeterminateDP {
			atLeastOneErrorDP = true
			continue
		}
	}
	if atLeastOneErrorDP {
		ret.AddResult(ResponseIndeterminateDP, "At Least One Error DP")
		return ret
	}
	if atLeastOneErrorD && (atLeastOneErrorP || atLeastOnePermit) {
		ret.AddResult(ResponseIndeterminateDP, "At Least One Error D & either At Least One Error P or At Least One Permit")
		return ret
	}
	if atLeastOneErrorD {
		ret.AddResult(ResponseIndeterminateD, "At Least One Error D")
		return ret
	}
	if atLeastOnePermit {
		ret.AddResult(ResponsePermit, "At Least One Permit")
		return ret
	}
	if atLeastOneErrorP {
		ret.AddResult(ResponseIndeterminateP, "At Least One Error P")
		return ret
	}

	ret.AddResult(ResponseNotApplicable, "No Errors or Permits")
	return ret
}

func permitOverridesCombiningAlgorithm(evaluatables []evaluatable, request Request) Response {
	ret := Response{make(map[string]interface{})}
	atLeastOneErrorD := false
	atLeastOneErrorP := false
	atLeastOneErrorDP := false
	atLeastOneDeny := false
	for _, e := range evaluatables {
		response := e.Evaluate(request)
		decision, status := response.GetResult()
		if decision == ResponseDeny {
			atLeastOneDeny = true
			continue
		}
		if decision == ResponsePermit {
			ret.AddResult(ResponsePermit, status)
			return ret
		}
		if decision == ResponseNotApplicable {
			continue
		}
		if decision == ResponseIndeterminateD {
			atLeastOneErrorD = true
			continue
		}
		if decision == ResponseIndeterminateP {
			atLeastOneErrorP = true
			continue
		}
		if decision == ResponseIndeterminateDP {
			atLeastOneErrorDP = true
			continue
		}
	}
	if atLeastOneErrorDP {
		ret.AddResult(ResponseIndeterminateDP, "At Least One Error DP")
		return ret
	}
	if atLeastOneErrorP && (atLeastOneErrorD || atLeastOneDeny) {
		ret.AddResult(ResponseIndeterminateDP, "At Least One Error P & either At Least One Error D or At Least One Deny")
		return ret
	}
	if atLeastOneErrorP {
		ret.AddResult(ResponseIndeterminateP, "At Least One Error P")
		return ret
	}
	if atLeastOneDeny {
		ret.AddResult(ResponseDeny, "At Least One Deny")
		return ret
	}
	if atLeastOneErrorD {
		ret.AddResult(ResponseIndeterminateD, "At Least One Error D")
		return ret
	}
	ret.AddResult(ResponseNotApplicable, "No Errors Or Denys")
	return ret
}

func denyUnlessPermitCombiningAlgorithm(evaluatables []evaluatable, request Request) Response {
	ret := Response{make(map[string]interface{})}
	for _, e := range evaluatables {
		response := e.Evaluate(request)
		decision, status := response.GetResult()
		if decision == ResponsePermit {
			ret.AddResult(ResponsePermit, status)
			return ret
		}
	}
	ret.AddResult(ResponseDeny, "No Permits, So Deny")
	return ret
}

func permitUnlessDenyCombiningAlgorithm(evaluatables []evaluatable, request Request) Response {
	ret := Response{make(map[string]interface{})}
	for _, e := range evaluatables {
		response := e.Evaluate(request)
		decision, status := response.GetResult()
		if decision == ResponseDeny {
			ret.AddResult(ResponseDeny, status)
			return ret
		}
	}
	ret.AddResult(ResponsePermit, "No Denys, So Permit")
	return ret
}

func firstApplicableEffectCombiningAlgorithm(evaluatables []evaluatable, request Request) Response {
	ret := Response{make(map[string]interface{})}
	for _, e := range evaluatables {
		response := e.Evaluate(request)
		decision, status := response.GetResult()
		if decision == ResponseDeny {
			ret.AddResult(ResponseDeny, status)
			return ret
		}
		if decision == ResponsePermit {
			ret.AddResult(ResponsePermit, status)
			return ret
		}
		if decision == ResponseNotApplicable {
			continue
		}
		if decision == ResponseIndeterminate {
			ret.AddResult(ResponseIndeterminate, status)
			return ret
		}
	}
	ret.AddResult(ResponseNotApplicable, "No Applicable decisions")
	return ret
}
