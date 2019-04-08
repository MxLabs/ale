// Code generated by "stringer -type=Opcode -linecomment"; DO NOT EDIT.

package isa

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Add-0]
	_ = x[Arg-1]
	_ = x[ArgLen-2]
	_ = x[Bind-3]
	_ = x[Call-4]
	_ = x[Call0-5]
	_ = x[Call1-6]
	_ = x[Closure-7]
	_ = x[CondJump-8]
	_ = x[Const-9]
	_ = x[Declare-10]
	_ = x[Div-11]
	_ = x[Dup-12]
	_ = x[Eq-13]
	_ = x[False-14]
	_ = x[Gt-15]
	_ = x[Gte-16]
	_ = x[Jump-17]
	_ = x[Load-18]
	_ = x[Lt-19]
	_ = x[Lte-20]
	_ = x[MakeCall-21]
	_ = x[MakeTruthy-22]
	_ = x[Mod-23]
	_ = x[Mul-24]
	_ = x[Neg-25]
	_ = x[NegInfinity-26]
	_ = x[NegOne-27]
	_ = x[Neq-28]
	_ = x[Nil-29]
	_ = x[NoOp-30]
	_ = x[Not-31]
	_ = x[One-32]
	_ = x[Panic-33]
	_ = x[Pop-34]
	_ = x[PosInfinity-35]
	_ = x[Resolve-36]
	_ = x[RestArg-37]
	_ = x[Return-38]
	_ = x[ReturnFalse-39]
	_ = x[ReturnNil-40]
	_ = x[ReturnTrue-41]
	_ = x[Self-42]
	_ = x[Store-43]
	_ = x[Sub-44]
	_ = x[True-45]
	_ = x[Two-46]
	_ = x[Zero-47]
}

const _Opcode_name = "AdditionRetrieve Argument ValueRetrieve Argument CountBind GlobalCall FunctionZero-Arg CallOne-Arg CallRetrieve Closure ValueConditional JumpRetrieve ConstantDeclare GlobalDivisionDuplicate ValueNumeric EqualityPush FalseGreater ThanGreater or Equal ToAbsolute JumpRetrieve Local ValueLess Than ComparisonLess or Equal ToMake Value CallableMake Value BooleanRemainderMultiplicationNegationPush Negative InfinityPush Negative OneNumeric InequalityPush NilNon-OperatorBoolean NegationPush OneAbnormally HaltDiscard ValuePositive InfinityResolve Global SymbolRetrieve Arguments TailReturn ValueReturn FalseReturn NilReturn TruePush Current FunctionStore LocalSubtractionPush TruePush TwoPush Zero"

var _Opcode_index = [...]uint16{0, 8, 31, 54, 65, 78, 91, 103, 125, 141, 158, 172, 180, 195, 211, 221, 233, 252, 265, 285, 305, 321, 340, 358, 367, 381, 389, 411, 428, 446, 454, 466, 482, 490, 505, 518, 535, 556, 579, 591, 603, 613, 624, 645, 656, 667, 676, 684, 693}

func (i Opcode) String() string {
	if i >= Opcode(len(_Opcode_index)-1) {
		return "Opcode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Opcode_name[_Opcode_index[i]:_Opcode_index[i+1]]
}
