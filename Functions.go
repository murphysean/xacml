package xacml

import (
	"errors"
	"log"
	"regexp"
)

const (
	functionStringEqual            = "urn:oasis:names:tc:xacml:1.0:function:string-equal"
	functionStringEqualIgnoreCase  = "urn:oasis:names:tc:xacml:1.0:function:string-equal-ignore-case"
	functionBooleanEqual           = "urn:oasis:names:tc:xacml:1.0:function:boolean-equal"
	functionIntegerEqual           = "urn:oasis:names:tc:xacml:1.0:function:integer-equal"
	functionDoubleEqual            = "urn:oasis:names:tc:xacml:1.0:function:double-equal"
	functionDateEqual              = "urn:oasis:names:tc:xacml:1.0:function:date-equal"
	functionTimeEqual              = "urn:oasis:names:tc:xacml:1.0:function:time-equal"
	functionDateTimeEqual          = "urn:oasis:names:tc:xacml:1.0:function:dateTime-equal"
	functionDayTimeDurationEqual   = "urn:oasis:names:tc:xacml:1.0:function:dayTimeDuration-equal"
	functionYearMonthDurationEqual = "urn:oasis:names:tc:xacml:1.0:function:yearMonthDuration-equal"
	functionAnyURIEqual            = "urn:oasis:names:tc:xacml:1.0:function:anyURI-equal"
	functionX500NameEqual          = "urn:oasis:names:tc:xacml:1.0:function:x500Name-equal"
	functionRFC822NameEqual        = "urn:oasis:names:tc:xacml:1.0:function:rfc822Name-equal"
	functionHexBinaryEqual         = "urn:oasis:names:tc:xacml:1.0:function:hexBinary-equal"
	functionBase64BinaryEqual      = "urn:oasis:names:tc:xacml:1.0:function:base64Binary-equal"

	functionIntegerAdd      = "urn:oasis:names:tc:xacml:1.0:function:integer-add"
	functionDoubleAdd       = "urn:oasis:names:tc:xacml:1.0:function:double-add"
	functionIntegerSubtract = "urn:oasis:names:tc:xacml:1.0:function:integer-subtract"
	functionDoubleSubtract  = "urn:oasis:names:tc:xacml:1.0:function:double-subtract"
	functionIntegerMultiply = "urn:oasis:names:tc:xacml:1.0:function:integer-multiply"
	functionDoubleMultiply  = "urn:oasis:names:tc:xacml:1.0:function:double-multiply"
	functionIntegerDivide   = "urn:oasis:names:tc:xacml:1.0:function:integer-divide"
	functionDoubleDivide    = "urn:oasis:names:tc:xacml:1.0:function:double-divide"
	functionIntegerMod      = "urn:oasis:names:tc:xacml:1.0:function:integer-mod"
	functionIntegerAbs      = "urn:oasis:names:tc:xacml:1.0:function:integer-abs"
	functionDoubleAbs       = "urn:oasis:names:tc:xacml:1.0:function:double-abs"
	functionRound           = "urn:oasis:names:tc:xacml:1.0:function:round"
	functionFloor           = "urn:oasis:names:tc:xacml:1.0:function:floor"

	functionStringNormalizeSpace       = "urn:oasis:names:tc:xacml:1.0:function:string-normalize-space"
	functionStringNormalizeToLowerCase = "urn:oasis:names:tc:xacml:1.0:function:string-normalize-to-lower-case"

	functionDoubleToInteger = "urn:oasis:names:tc:xacml:1.0:function:double-to-integer"
	functionIntegerToDouble = "urn:oasis:names:tc:xacml:1.0:function:integer-to-double"

	functionOr  = "urn:oasis:names:tc:xacml:1.0:function:or"
	functionAnd = "urn:oasis:names:tc:xacml:1.0:function:and"
	functionNOf = "urn:oasis:names:tc:xacml:1.0:function:n-of"
	functionNot = "urn:oasis:names:tc:xacml:1.0:function:not"

	functionIntegerGreaterThan        = "urn:oasis:names:tc:xacml:1.0:function:integer-greater-than"
	functionIntegerGreaterThanOrEqual = "urn:oasis:names:tc:xacml:1.0:function:integer-greater-than-or-equal"
	functionIntegerLessThan           = "urn:oasis:names:tc:xacml:1.0:function:integer-less-than"
	functionIntegerLessThanOrEqual    = "urn:oasis:names:tc:xacml:1.0:function:integer-less-than-or-equal"
	functionDoubleGreaterThan         = "urn:oasis:names:tc:xacml:1.0:function:double-greater-than"
	functionDoubleGreaterThanOrEqual  = "urn:oasis:names:tc:xacml:1.0:function:double-greater-than-or-equal"
	functionDoubleLessThan            = "urn:oasis:names:tc:xacml:1.0:function:double-less-than"
	functionDoubleLessThanOrEqual     = "urn:oasis:names:tc:xacml:1.0:function:double-less-than-or-equal"

	functionStringConcatenate = "urn:oasis:names:tc:xacml:2.0:function:string-concatenate"
	functionBooleanFromString = "urn:oasis:names:tc:xacml:3.0:function:boolean-from-string"
	functionStringFromBoolean = "urn:oasis:names:tc:xacml:3.0:function:string-from-boolean"
	functionIntegerFromString = "urn:oasis:names:tc:xacml:3.0:function:integer-from-string"
	functionStringFromInteger = "urn:oasis:names:tc:xacml:3.0:function:string-from-integer"
	functionDoubleFromString  = "urn:oasis:names:tc:xacml:3.0:function:double-from-string"
	functionStringFromDouble  = "urn:oasis:names:tc:xacml:3.0:function:string-from-double"

	functionStringStartsWith = "urn:oasis:names:tc:xacml:3.0:function:string-starts-with"
	functionStringEndsWith   = "urn:oasis:names:tc:xacml:3.0:function:string-ends-with"
	functionStringContains   = "urn:oasis:names:tc:xacml:3.0:function:string-contains"
	functionStringSubString  = "urn:oasis:names:tc:xacml:3.0:function:string-substring"

	functionStringOneAndOnly = "urn:oasis:names:tc:xacml:3.0:function:string-one-and-only"
	functionStringBagSize    = "urn:oasis:names:tc:xacml:3.0:function:string-bag-size"
	functionStringIsIn       = "urn:oasis:names:tc:xacml:3.0:function:string-is-in"
	functionStringBag        = "urn:oasis:names:tc:xacml:3.0:function:string-bag"

	functionStringIntersection       = "urn:oasis:names:tc:xacml:3.0:function:string-intersection"
	functionStringAtLeastOneMemberOf = "urn:oasis:names:tc:xacml:3.0:function:string-at-least-one-member-of"
	functionStringUnion              = "urn:oasis:names:tc:xacml:3.0:function:string-union"
	functionStringSubset             = "urn:oasis:names:tc:xacml:3.0:function:string-subset"
	functionStringSetEquals          = "urn:oasis:names:tc:xacml:3.0:function:string-set-equals"

	functionStringRegExpMatch    = "urn:oasis:names:tc:xacml:1.0:function:string-regexp-match"
	functionAnyURIRegExpMatch    = "urn:oasis:names:tc:xacml:1.0:function:anyURI-regexp-match"
	functionIpAddressRegExpMatch = "urn:oasis:names:tc:xacml:1.0:function:ipAddress-regexp-match"
	functionDNSNameRegExpMatch   = "urn:oasis:names:tc:xacml:1.0:function:dnsName-regexp-match"
)

func applyFunction(functionId string, args ...interface{}) (interface{}, error) {
	switch functionId {
	case functionStringEqual:
		return stringEqual(args[0], args[1])
	case functionBooleanEqual:
		return booleanEqual(args[0], args[1])
	case functionIntegerEqual:
		return integerEqual(args[0], args[1])
	case functionStringOneAndOnly:
		return stringOneAndOnly(args[0])
	case functionStringBagSize:
		return stringBagSize(args[0])
	case functionStringIsIn:
		return stringIsIn(args[0], args[1])
	case functionStringBag:
		return stringBag(args...)
	case functionStringRegExpMatch:
		return stringRegExpMatch(args[0], args[1])
	default:
		log.Println("Unknown or Unsupported Function: ", functionId)
		return false, errors.New("Unknown or Unsupported Function")
	}
}

func applyMatchFunction(functionId string, attributeOne, attributeTwo interface{}) (bool, error) {
	ret, err := applyFunction(functionId, attributeOne, attributeTwo)
	if err != nil {
		return false, err
	}
	retb, _ := ret.(bool)
	return retb, err
}

func stringEqual(attributeOne, attributeTwo interface{}) (bool, error) {
	a1, ok := attributeOne.(string)
	if !ok {
		return false, errors.New("Invalid attributeOne type")
	}
	a2, ok := attributeTwo.(string)
	if !ok {
		return false, errors.New("Invalid attributeTwo type")
	}
	return a1 == a2, nil
}

func stringOneAndOnly(attributeOne interface{}) (string, error) {
	attributeArray, ok := attributeOne.([]interface{})
	if !ok {
		return "", errors.New("Argument must be an array")
	}
	if len(attributeArray) > 1 {
		return "", errors.New("Argument must be an array with only one value")
	}
	stringVal, ok := attributeArray[0].(string)
	if !ok {
		return "", errors.New("Each member of the array must be a string")
	}
	return stringVal, nil
}

func stringBagSize(attributeOne interface{}) (int64, error) {
	attributeArray, ok := attributeOne.([]interface{})
	if !ok {
		return 0, errors.New("Argument must be an array")
	}
	return int64(len(attributeArray)), nil
}

func stringIsIn(attributeOne, attributeTwo interface{}) (bool, error) {
	//log.Println(attributeOne)
	//log.Println(attributeTwo)
	a1, ok := attributeOne.(string)
	if !ok {
		return false, errors.New("First argument must be a string")
	}
	attributeArray, ok := attributeTwo.([]interface{})
	if !ok {
		attributeValue, ok := attributeTwo.(string)
		if !ok {
			return false, errors.New("Second argument must be an array, or a string")
		}
		attributeArray = append(attributeArray, attributeValue)
	}
	for _, s := range attributeArray {
		ret, err := stringEqual(a1, s)
		if err != nil {
			continue
		}
		if ret {
			return true, nil
		}
	}
	return false, nil
}

func stringBag(attributes ...interface{}) ([]interface{}, error) {
	ret := make([]interface{}, 0)
	for _, attr := range attributes {
		if str, ok := attr.(string); ok {
			ret = append(ret, str)
		}
	}
	return ret, nil
}

func stringRegExpMatch(attributeOne, attributeTwo interface{}) (bool, error) {
	a1, ok := attributeOne.(string)
	if !ok {
		return false, errors.New("Invalid attributeOne type")
	}
	a2, ok := attributeTwo.(string)
	if !ok {
		return false, errors.New("Invalid attributeTwo type")
	}
	return regexp.MatchString(a1, a2)
}

func booleanEqual(attributeOne, attributeTwo interface{}) (bool, error) {
	a1, ok := attributeOne.(bool)
	if !ok {
		return false, errors.New("Invalid attributeOne type")
	}
	a2, ok := attributeTwo.(bool)
	if !ok {
		return false, errors.New("Invalid attributeTwo type")
	}
	return a1 == a2, nil
}

func integerEqual(attributeOne, attributeTwo interface{}) (bool, error) {
	a1, ok := attributeOne.(int64)
	if !ok {
		return false, errors.New("Invalid attributeOne type")
	}
	a2, ok := attributeTwo.(int64)
	if !ok {
		return false, errors.New("Invalid attributeTwo type")
	}
	return a1 == a2, nil
}
