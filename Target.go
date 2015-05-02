package xacml

import (
	"log"
)

const (
	targetMatch         = "Match"
	targetNoMatch       = "No Match"
	targetIndeterminate = "Indeterminate"
)

func evaluateTarget(targetInterface map[string]interface{}, request Request) string {
	//log.Println("Target")
	//b, _ := json.Marshal(targetInterface)
	//log.Println("Target ", string(b))
	return evaluateTargetBody(targetInterface["Target"], request)
}

func evaluateTargetBody(targetInterface interface{}, request Request) string {
	//log.Println("Target Body")
	//b, _ := json.Marshal(targetInterface)
	//log.Println("TargetBody ", string(b))

	switch targetInterface.(type) {
	case map[string]interface{}:
		targetBody, _ := targetInterface.(map[string]interface{})
		switch targetBody["AnyOf"].(type) {
		case []interface{}:
			return evaluateAnyOfs(targetBody["AnyOf"].([]interface{}), request)
		case map[string]interface{}:
			return evaluateAnyOfBody(targetBody["AnyOf"].(map[string]interface{}), request)
		}
	case string:
		//log.Println("TargetBody Return: ", targetMatch)
		return targetMatch
	case nil:
		//log.Println("TargetBody Return: ", targetMatch)
		return targetMatch
	}
	//log.Println("TargetBody Return: ", targetIndeterminate)
	return targetIndeterminate
}

func evaluateAnyOfs(anyOfs []interface{}, request Request) string {
	//log.Println("AnyOfs")
	//b, _ := json.Marshal(anyOfs)
	//log.Println("AnyOfs ", string(b))
	allMatch := true

	for _, anyOfInterface := range anyOfs {
		result := targetIndeterminate
		switch anyOfInterface.(type) {
		case map[string]interface{}:
			result = evaluateAnyOfBody(anyOfInterface, request)
		}
		//If any is a "No Match" return "No Match"
		if result == targetNoMatch {
			//log.Println("AnyOfs Return: ", targetNoMatch)
			return targetNoMatch
		}
		if result != targetMatch {
			allMatch = false
		}
	}

	//If all match, target is match
	if allMatch {
		//log.Println("AnyOfs Return: ", targetMatch)
		return targetMatch
	}
	//Otherwise indeterminate
	//log.Println("AnyOfs Return: ", targetIndeterminate)
	return targetIndeterminate
}

func evaluateAnyOf(anyOfInterface map[string]interface{}, request Request) string {
	//log.Println("AnyOf")
	//b, _ := json.Marshal(anyOfInterface)
	//log.Println("AnyOf ", string(b))
	return evaluateAnyOfBody(anyOfInterface["AnyOf"], request)
}

func evaluateAnyOfBody(anyOfInterface interface{}, request Request) string {
	//log.Println("AnyOfBody")
	//b, _ := json.Marshal(anyOfInterface)
	//log.Println("AnyOfBody ", string(b))
	switch anyOfInterface.(type) {
	case map[string]interface{}:
		anyOfBody, _ := anyOfInterface.(map[string]interface{})
		switch anyOfBody["AllOf"].(type) {
		case []interface{}:
			return evaluateAllOfs(anyOfBody["AllOf"].([]interface{}), request)
		case map[string]interface{}:
			return evaluateAllOfBody(anyOfBody["AllOf"].(map[string]interface{}), request)
		}
	}
	//log.Println("AnyOfBody Return: ", targetIndeterminate)
	return targetIndeterminate
}

func evaluateAllOfs(allOfs []interface{}, request Request) string {
	//log.Println("AllOfs")
	//b, _ := json.Marshal(allOfs)
	//log.Println("AllOfs ", string(b))
	allNoMatch := true

	for _, allOfInterface := range allOfs {
		result := targetIndeterminate
		switch allOfInterface.(type) {
		case map[string]interface{}:
			result = evaluateAllOfBody(allOfInterface.(map[string]interface{}), request)
		}
		//At least one match, return match
		if result == targetMatch {
			//log.Println("AllOfs Return: ", targetMatch)
			return targetMatch
		}
		if result != targetNoMatch {
			allNoMatch = false
		}
	}

	//All no match, return no match
	if allNoMatch {
		//log.Println("AllOfs Return: ", targetNoMatch)
		return targetNoMatch
	}
	//Otherwise indeterminate
	//log.Println("AllOfs Return: ", targetIndeterminate)
	return targetIndeterminate
}

func evaluateAllOf(allOf map[string]interface{}, request Request) string {
	//log.Println("AllOf")
	//b, _ := json.Marshal(allOf)
	//log.Println("AllOf ", string(b))
	return evaluateAllOfBody(allOf["AllOf"], request)
}

func evaluateAllOfBody(allOfInterface interface{}, request Request) string {
	//log.Println("AllOfBody")
	//b, _ := json.Marshal(allOfInterface)
	//log.Println("AllOfBody ", string(b))
	switch allOfInterface.(type) {
	case map[string]interface{}:
		allOfBody, _ := allOfInterface.(map[string]interface{})
		switch allOfBody["Match"].(type) {
		case []interface{}:
			return evaluateMatchs(allOfBody["Match"].([]interface{}), request)
		case map[string]interface{}:
			return evaluateMatchBody(allOfBody["Match"].(map[string]interface{}), request)
		}
	}

	//log.Println("AllOfBody Return: ", targetIndeterminate)
	return targetIndeterminate
}

func evaluateMatchs(matchs []interface{}, request Request) string {
	//log.Println("Matchs")
	//b, _ := json.Marshal(matchs)
	//log.Println("Matchs ", string(b))
	allTrue := true

	for _, matchInterface := range matchs {
		result := targetIndeterminate
		switch matchInterface.(type) {
		case map[string]interface{}:
			result = evaluateMatchBody(matchInterface.(map[string]interface{}), request)
		}
		//At least one no match, return no match
		if result == targetNoMatch {
			//log.Println("Matchs Return: ", targetNoMatch)
			return targetNoMatch
		}
		if result != targetMatch {
			allTrue = false
		}
	}

	//All no match, return no match
	if allTrue {
		//log.Println("Matchs Return: ", targetMatch)
		return targetMatch
	}
	//Otherwise indeterminate
	//log.Println("Matchs Return: ", targetIndeterminate)
	return targetIndeterminate
}

func evaluateMatch(match map[string]interface{}, request Request) string {
	//log.Println("Match")
	//b, _ := json.Marshal(match)
	//log.Println("Match ", string(b))
	return evaluateMatchBody(match["Match"], request)
}

func evaluateMatchBody(matchInterface interface{}, request Request) string {
	//log.Println("MatchBody")
	//b, _ := json.Marshal(matchInterface)
	//log.Println("MatchBody ", string(b))
	matchBody, _ := matchInterface.(map[string]interface{})
	//Identify Function type
	functionId, ok := matchBody["-MatchId"].(string)
	if !ok {
		log.Println("functionId not ok")
		return targetIndeterminate
	}

	//Read in the attribute(bag) value
	attributeValueInterface, ok := matchBody["AttributeValue"].(map[string]interface{})
	if !ok {
		log.Println("attributeValueInterface not ok")
		return targetIndeterminate
	}
	var err error
	attributeValue, err := resolveAttributeValue(attributeValueInterface, request)
	if err != nil {
		log.Println(err)
		log.Println("attributeValue not ok")
		return targetIndeterminate
	}

	//Is attribute value 0?
	attributeSeq, _ := attributeValueInterface["_seq"].(string)

	//Read in the attribute designator, or selector (not supported for now)
	attributeDesignatorInterface, ok := matchBody["AttributeDesignator"].(map[string]interface{})
	if !ok {
		log.Println("attributeDesignatorInterface not ok")
		return targetIndeterminate
	}
	attributeDesignator, err := resolveAttributeDesignator(attributeDesignatorInterface, request)
	if err != nil {
		log.Println("attributeDesignator not ok")
		return targetIndeterminate
	}

	var result bool
	if attributeSeq == "0" {
		result, err = applyMatchFunction(functionId, attributeValue, attributeDesignator)
	} else {
		result, err = applyMatchFunction(functionId, attributeDesignator, attributeValue)
	}
	if err != nil {
		log.Println(err)
		log.Println("Function Failed to produce a result")
		return targetIndeterminate
	}
	if result {
		return targetMatch
	} else {
		return targetNoMatch
	}
}
