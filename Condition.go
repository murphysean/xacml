package xacml

import (
	"errors"
	"strconv"
	"strings"
)

const (
	ConditionTrue          = "True"
	ConditionFalse         = "False"
	ConditionIndeterminate = "Indeterminate"
)

func evaluateCondition(condition map[string]interface{}, request Request) string {
	//b, _ := json.Marshal(condition)
	//log.Println("Condition ", string(b))
	return evaluateConditionBody(condition["Condition"], request)
}

func evaluateConditionBody(conditionBodyInterface interface{}, request Request) string {
	//b, _ := json.Marshal(conditionBodyInterface)
	//log.Println("ConditionBody ", string(b))
	//log.Println(reflect.TypeOf(conditionBodyInterface))
	var conditionBody map[string]interface{}
	switch conditionBodyInterface.(type) {
	case map[string]interface{}:
		conditionBody, _ = conditionBodyInterface.(map[string]interface{})
	default:
		return ConditionTrue
	}
	//Can be any of: Apply, AttributeSelector, AttributeValue, Function, VariableReference, AttributeDesignator
	//Not supporting AttributeSelector or Variable Reference or Function for now

	//Contains Apply?
	if applyBody, ok := conditionBody["Apply"].(map[string]interface{}); ok && applyBody != nil {
		ret, err := evaluateApplyBody(applyBody, request)
		if err != nil {
			return ConditionIndeterminate
		}
		switch ret.(type) {
		case bool:
			if ret.(bool) {
				return ConditionTrue
			} else {
				return ConditionFalse
			}
		default:
			return ConditionIndeterminate
		}
	}

	//attributeValue, _ := conditionBody["AttributeValue"].(map[string]interface{})
	//attributeDesignator, _ := conditionBody["AttributeDesignator"].(map[string]interface{})
	return ConditionIndeterminate
}

func evaluateApply(apply map[string]interface{}, request Request) (interface{}, error) {
	//b, _ := json.Marshal(applyInterface)
	//log.Println("Apply ", string(b))
	if applyBody, ok := apply["Apply"].(map[string]interface{}); ok && applyBody != nil {
		return evaluateApplyBody(applyBody, request)
	}
	return nil, errors.New("Apply Body not a map")
}

func evaluateApplyBody(applyBodyInterface map[string]interface{}, request Request) (interface{}, error) {
	//b, _ := json.Marshal(applyBodyInterface)
	//log.Println("ApplyBody ", string(b))

	functionId, _ := applyBodyInterface["-FunctionId"].(string)
	//log.Println(functionId)

	argCount := 0
	for k, v := range applyBodyInterface {
		if strings.HasSuffix(k, "-") || strings.HasSuffix(k, "_") {
			continue
		}

		switch v.(type) {
		case map[string]interface{}:
			argCount++
		case []interface{}:
			argCount += len(v.([]interface{}))
		}
	}

	args := make([]interface{}, argCount)

	for k, v := range applyBodyInterface {
		switch k {
		case "AttributeValue":
			scrapeAttributeValueArgs(args, v, request)
		case "AttributeDesignator":
			scrapeAttributeDesignatorArgs(args, v, request)
		case "Apply":
			innerApplyBody, _ := v.(map[string]interface{})
			argNum, _ := strconv.Atoi(innerApplyBody["_seq"].(string))
			applyResult, _ := evaluateApplyBody(innerApplyBody, request)
			args[argNum] = applyResult
		}
	}
	return applyFunction(functionId, args...)
}

func scrapeAttributeValueArgs(args []interface{}, attributeValueInterface interface{}, request Request) {
	//b, _ := json.Marshal(attributeValueInterface)
	//log.Println("scrapeAttributeValueArgs ", string(b))
	switch attributeValueInterface.(type) {
	case map[string]interface{}:
		aviMap := attributeValueInterface.(map[string]interface{})
		argNum, _ := strconv.Atoi(aviMap["_seq"].(string))
		arg, err := resolveAttributeValue(aviMap, request)
		if err != nil {
			args[argNum] = err
		} else {
			args[argNum] = arg
		}
	case []interface{}:
		for _, avi := range attributeValueInterface.([]interface{}) {
			aviMap, _ := avi.(map[string]interface{})
			argNum, _ := strconv.Atoi(aviMap["_seq"].(string))
			arg, err := resolveAttributeValue(aviMap, request)
			if err != nil {
				args[argNum] = err
			} else {
				args[argNum] = arg
			}
		}
	}
}

func scrapeAttributeDesignatorArgs(args []interface{}, attributeDesignatorInterface interface{}, request Request) {
	//b, _ := json.Marshal(attributeDesignatorInterface)
	//log.Println("scrapeAttributeDesignatorArgs ", string(b))
	switch attributeDesignatorInterface.(type) {
	case map[string]interface{}:
		adiMap := attributeDesignatorInterface.(map[string]interface{})
		argNum, _ := strconv.Atoi(adiMap["_seq"].(string))
		arg, err := resolveAttributeDesignator(attributeDesignatorInterface.(map[string]interface{}), request)
		if err != nil {
			args[argNum] = err
		} else {
			args[argNum] = arg
		}
	case []interface{}:
		for _, adi := range attributeDesignatorInterface.([]interface{}) {
			adiMap, _ := adi.(map[string]interface{})
			argNum, _ := strconv.Atoi(adiMap["_seq"].(string))
			arg, err := resolveAttributeDesignator(adiMap, request)
			if err != nil {
				args[argNum] = err
			} else {
				args[argNum] = arg
			}
		}
	}
}
