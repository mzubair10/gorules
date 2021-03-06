package gorules

import (
	"errors"
	"strconv"
)

// Value refers to anytype that can be evaluated to a concrete string value
type Value interface {
	Evaluate(interface{}) (string, error)
	String() string
}

// Constant is used to hold the string value
type Constant struct {
	value string
}

// Evaluate returns the string from the Constant
func (c Constant) Evaluate(_ interface{}) (string, error) {
	if startsWithSingleQuotes(c.value) {
		return StringBetweenSingleQuotes(c.value), nil
	}
	return "", errors.New("Not a Constant")
}

// String makes Constant implement Stringer
func (c Constant) String() string {
	return c.value
}

// NewConstant creates new Constant which is within single quotes.Creates an empty string if value has no quotes
func NewConstant(value string) Constant {
	return Constant{value: value}
}

// Path has the JSON Path. Needs data to be evaluated to the final string value
type Path struct {
	jsonPath string
}

// Evaluate returns the string from the Constant
func (p Path) Evaluate(data interface{}) (string, error) {
	// fmt.Println(p, data)
	return selectValue(data.(map[string]interface{}), p.jsonPath).(string), nil
}

// String makes Path implement Stringer
func (p Path) String() string {
	return p.jsonPath
}

// NewPath creates new JSON path which can be evaluated with supplied data
func NewPath(value string) Path {
	return Path{jsonPath: value}
}

// NewValue used to create any of the value type
func NewValue(value string) Value {
	if startsWithSingleQuotes(value) {
		return NewConstant(value)
	}
	return NewPath(value)
}

// MathExpression is used to evaluate mathematical expressions on json values
type MathExpression struct {
	operand1 Value
	operand2 Value
	operator MathOperator
}

//NewMathExpression is a wrapper around MathExpression
func NewMathExpression(expression string) MathExpression {
	parsedOperandsAndOperatorValue := parsedOperandsAndOperator(expression)
	multiplicationOperator, _ := toMathOperator(parsedOperandsAndOperatorValue[1])
	return MathExpression{operand1: NewValue(trim(parsedOperandsAndOperatorValue[0])), operand2: NewValue(trim(parsedOperandsAndOperatorValue[2])), operator: multiplicationOperator}
}

// Evaluate works out the expression and returns the result as a string
func (m MathExpression) Evaluate(_ interface{}) (string, error) {
	operand1, _ := m.operand1.Evaluate(make([]interface{}, 0))
	operand2, _ := m.operand2.Evaluate(make([]interface{}, 0))
	firstOperand, _ := strconv.Atoi(operand1)
	secondOperand, _ := strconv.Atoi(operand2)
	mathOperatorFunc := mathOperatorFuncList[m.operator]
	result, err := mathOperatorFunc(firstOperand, secondOperand)
	return strconv.Itoa(result), err
}

func (m MathExpression) String() string {
	dummyValue, _ := m.operand1.Evaluate(make([]interface{}, 0))
	return dummyValue
}
