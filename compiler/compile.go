package compiler

import "gitlab.com/kode4food/ale/data"

// IsEvaluable returns whether or not the provided value is subject
// to further evaluation
func IsEvaluable(v data.Value) bool {
	switch v.(type) {
	case data.String:
		return false
	case data.Sequence, data.Symbol:
		return true
	default:
		return false
	}
}
