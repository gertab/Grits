package process

import (
	"testing"
)

func compareOutput(t *testing.T, got []string, expected []string) {
	if len(got) != len(expected) {
		t.Errorf("len of got %d, does not match len of expected %d\n", len(got), len(expected))
		return
	}

	for index := range got {
		if got[index] != expected[index] {
			// t.Errorf("got %s, expected %s\n", "sa", "de")
			t.Errorf("got %s, expected %s\n", got[index], expected[index])
		}
	}
}

func assertEqualNames(t *testing.T, got []Name, expected []Name) {
	if len(got) != len(expected) {
		t.Errorf("len of got %d, does not match len of expected %d\n", len(got), len(expected))
		return
	}

	for index := range got {
		if !got[index].Equal(expected[index]) {
			// t.Errorf("got %s, expected %s\n", "sa", "de")
			t.Errorf("got %s, expected %s\n", got[index].String(), expected[index].String())
		}
	}
}

func assertEqual(t *testing.T, i1, i2 Form) {
	if !EqualForm(i1, i2) {
		t.Errorf("got %s, expected %s\n", i1.String(), i2.String())
	}
}

func assertNotEqual(t *testing.T, i1, i2 Form) {
	if EqualForm(i1, i2) {
		t.Errorf("got %s, expected %s\n", i1.String(), i2.String())
	}
}

func compareOutputProgram(t *testing.T, got []Form, expected []Form) {
	if len(got) != len(expected) {
		t.Errorf("len of got %d, does not match len of expected %d\n", len(got), len(expected))
		return
	}

	for index := range got {
		if !EqualForm(got[index], expected[index]) {
			t.Errorf("[%d] got %s, expected %s\n", index, got[index].String(), expected[index].String())
		}
	}
}

func TestBasicTokens(t *testing.T) {
	expected, output := []string{}, []string{}
	to_c := Name{Ident: "to_c"}
	pay_c := Name{Ident: "pay_c"}
	cont_c := Name{Ident: "cont_c"}
	from_c := Name{Ident: "from_c"}
	end := NewClose(Name{Ident: "self"})

	// Send
	input := NewSend(to_c, pay_c, cont_c)
	output = append(output, input.String())
	expected = append(expected, "send to_c<pay_c,cont_c>")

	// Receive
	input2 := NewReceive(pay_c, cont_c, from_c, end)
	output = append(output, input2.String())
	expected = append(expected, "<pay_c,cont_c> <- recv from_c; close self")

	// Case
	input4 := NewCase(from_c, []*BranchForm{NewBranch(Label{L: "label1"}, pay_c, end)})
	output = append(output, input4.String())
	expected = append(expected, "case from_c (label1<pay_c> => close self)")

	// Case (multiple cases)
	input5 := NewCase(from_c, []*BranchForm{NewBranch(Label{L: "label1"}, pay_c, end), NewBranch(Label{L: "label2"}, pay_c, end)})
	output = append(output, input5.String())
	expected = append(expected, "case from_c (label1<pay_c> => close self | label2<pay_c> => close self)")

	// Select
	input6 := NewSelect(to_c, Label{L: "label1"}, cont_c)
	output = append(output, input6.String())
	expected = append(expected, "to_c.label1<cont_c>")

	// New
	input7 := NewNew(cont_c, end, end)
	output = append(output, input7.String())
	expected = append(expected, "cont_c <- new (close self); close self")

	// Close
	input8 := NewClose(from_c)
	output = append(output, input8.String())
	expected = append(expected, "close from_c")

	// Forward
	input9 := NewForward(to_c, from_c)
	output = append(output, input9.String())
	expected = append(expected, "fwd to_c from_c")

	// Split
	input10 := NewSplit(pay_c, cont_c, from_c, end)
	output = append(output, input10.String())
	expected = append(expected, "<pay_c,cont_c> <- split from_c; close self")

	// Wait
	input11 := NewWait(to_c, end)
	output = append(output, input11.String())
	expected = append(expected, "wait to_c; close self")

	// Cast
	input12 := NewCast(to_c, cont_c)
	output = append(output, input12.String())
	expected = append(expected, "cast to_c<cont_c>")

	// Receive
	input13 := NewShift(cont_c, from_c, end)
	output = append(output, input13.String())
	expected = append(expected, "cont_c <- shift from_c; close self")

	// Wait
	input14 := NewDrop(to_c, end)
	output = append(output, input14.String())
	expected = append(expected, "drop to_c; close self")

	compareOutput(t, output, expected)
}

func TestSubstitutions(t *testing.T) {
	expected, output := []Form{}, []Form{}
	to_c := Name{Ident: "to_c", IsSelf: false}
	new_to_c := Name{Ident: "new_to_c", IsSelf: false}
	pay_c := Name{Ident: "pay_c", IsSelf: false}
	new_pay_c := Name{Ident: "new_pay_c", IsSelf: false}
	cont_c := Name{Ident: "cont_c", IsSelf: false}
	new_cont_c := Name{Ident: "new_cont_c", IsSelf: false}
	from_c := Name{Ident: "from_c", IsSelf: false}
	new_from_c := Name{Ident: "new_from_c", IsSelf: false}
	self := Name{IsSelf: true}
	new_self := Name{IsSelf: true}
	end := NewClose(self)
	new_end := NewClose(new_self)

	// Send
	input := NewSend(to_c, pay_c, cont_c)
	input.Substitute(to_c, new_to_c)
	input.Substitute(pay_c, new_pay_c)
	input.Substitute(cont_c, new_cont_c)
	output = append(output, input)
	result := NewSend(new_to_c, new_pay_c, new_cont_c)
	expected = append(expected, result)

	// Receive
	input2 := NewReceive(pay_c, cont_c, from_c, end)
	input2.Substitute(cont_c, new_cont_c)
	input2.Substitute(pay_c, new_pay_c)
	input2.Substitute(from_c, new_from_c)
	input2.Substitute(self, new_self)
	result2 := NewReceive(pay_c, cont_c, new_from_c, new_end)
	output = append(output, input2)
	expected = append(expected, result2)

	// Receive
	input2other := NewReceive(pay_c, cont_c, pay_c, end)
	input2other.Substitute(cont_c, new_cont_c)
	input2other.Substitute(pay_c, new_pay_c)
	input2other.Substitute(self, new_self)
	result2other := NewReceive(pay_c, cont_c, new_pay_c, new_end)
	output = append(output, input2other)
	expected = append(expected, result2other)

	// Case
	input3 := NewCase(from_c, []*BranchForm{NewBranch(Label{L: "from_c"}, pay_c, end)})
	input3.Substitute(from_c, new_from_c)
	input3.Substitute(pay_c, new_pay_c)
	input3.Substitute(cont_c, new_cont_c)
	input3.Substitute(self, new_self)
	result3 := NewCase(new_from_c, []*BranchForm{NewBranch(Label{L: "from_c"}, pay_c, new_end)})
	output = append(output, input3)
	expected = append(expected, result3)

	// Select
	input4 := NewSelect(to_c, Label{L: "label1"}, cont_c)
	input4.Substitute(to_c, new_to_c)
	input4.Substitute(cont_c, new_cont_c)
	result4 := NewSelect(new_to_c, Label{L: "label1"}, new_cont_c)
	output = append(output, input4)
	expected = append(expected, result4)

	// New
	input5 := NewNew(cont_c, end, end)
	input5.Substitute(cont_c, new_cont_c)
	input5.Substitute(self, new_self)
	result5 := NewNew(cont_c, end, end)
	output = append(output, input5)
	expected = append(expected, result5)

	// Close
	input6 := NewClose(from_c)
	input6.Substitute(from_c, new_from_c)
	input6.Substitute(self, new_self)
	result6 := NewClose(new_from_c)
	output = append(output, input6)
	expected = append(expected, result6)

	// Forward
	input7 := NewForward(to_c, from_c)
	input7.Substitute(from_c, new_from_c)
	input7.Substitute(to_c, new_to_c)
	result7 := NewForward(new_to_c, new_from_c)
	output = append(output, input7)
	expected = append(expected, result7)

	// Split
	input8 := NewSplit(pay_c, cont_c, from_c, end)
	input8.Substitute(pay_c, new_pay_c)
	input8.Substitute(cont_c, new_cont_c)
	input8.Substitute(from_c, new_from_c)
	input8.Substitute(self, new_self)
	result8 := NewSplit(pay_c, cont_c, new_from_c, new_end)
	output = append(output, input8)
	expected = append(expected, result8)

	// Call
	input9 := NewCall("func_name", []Name{from_c, to_c})
	input9.Substitute(pay_c, new_pay_c)
	input9.Substitute(cont_c, new_cont_c)
	input9.Substitute(from_c, new_from_c)
	input9.Substitute(self, new_self)
	result9 := NewCall("func_name", []Name{new_from_c, to_c})
	output = append(output, input9)
	expected = append(expected, result9)

	// Wait
	input10 := NewWait(to_c, end)
	input10.Substitute(to_c, new_to_c)
	input10.Substitute(self, new_self)
	result11 := NewWait(new_to_c, end)
	output = append(output, input10)
	expected = append(expected, result11)

	// Cast
	input12 := NewCast(pay_c, cont_c)
	input12.Substitute(pay_c, new_pay_c)
	input12.Substitute(cont_c, new_cont_c)
	output = append(output, input12)
	result12 := NewCast(new_pay_c, new_cont_c)
	expected = append(expected, result12)

	// Shift
	input13 := NewShift(cont_c, from_c, end)
	input13.Substitute(cont_c, new_cont_c)
	input13.Substitute(from_c, new_from_c)
	input13.Substitute(self, new_self)
	result13 := NewShift(cont_c, new_from_c, new_end)
	output = append(output, input13)
	expected = append(expected, result13)

	compareOutputProgram(t, output, expected)
}

func TestCopy(t *testing.T) {
	to_c := Name{Ident: "to_c", IsSelf: false}
	// new_to_c := Name{Ident: "new_to_c", IsSelf: false}
	pay_c := Name{Ident: "pay_c", IsSelf: false}
	// new_pay_c := Name{Ident: "new_pay_c", IsSelf: false}
	cont_c := Name{Ident: "cont_c", IsSelf: false}
	// new_cont_c := Name{Ident: "new_cont_c", IsSelf: false}
	from_c := Name{Ident: "from_c", IsSelf: false}
	// new_from_c := Name{Ident: "new_from_c", IsSelf: false}
	self := Name{IsSelf: true}
	// new_self := Name{IsSelf: true}
	end := NewClose(self)
	// new_end := NewClose(new_self)

	// Send
	input := NewSend(Name{Ident: "to_c", IsSelf: false}, pay_c, cont_c)
	copy := CopyForm(input)
	copyWithType := copy.(*SendForm)
	copyWithType.to_c.Ident = "to_c"
	assertEqual(t, input, copy)
	copyWithType.to_c.Ident = "to_c_edited"
	assertNotEqual(t, input, copyWithType)
	c1 := make(chan Message)
	c2 := make(chan Message)
	// c2 := make(chan Message)
	input.to_c.Channel = c1
	copyWithType.to_c.Ident = "to_c"
	copyWithType.to_c.Channel = c1
	assertEqual(t, input, copyWithType)
	copyWithType.to_c.Channel = c2
	assertNotEqual(t, input, copyWithType)
	if input.to_c.Channel == copyWithType.to_c.Channel {
		t.Errorf("Channels shouldn't be equal: got %s and %s\n", input.to_c.String(), copyWithType.to_c.String())
	}

	// Receive
	input2 := NewReceive(pay_c, cont_c, from_c, end)
	copy = CopyForm(input2)
	copyWithType2 := copy.(*ReceiveForm)
	copyWithType2.payload_c.Ident = "pay_c"
	assertEqual(t, input2, copy)
	copyWithType2.payload_c.Ident = "pay_c_edited"
	assertNotEqual(t, input2, copy)

	// Case
	input3 := NewCase(from_c, []*BranchForm{NewBranch(Label{L: "from_c"}, pay_c, end)})
	copy3 := CopyForm(input3)
	copyWithType3 := copy3.(*CaseForm)
	copyWithType3.from_c.Ident = "from_c"
	assertEqual(t, input3, copy3)
	copyWithType3.from_c.Ident = "from_c_edited"
	assertNotEqual(t, input3, copy3)

	// Select
	input4 := NewSelect(to_c, Label{L: "label1"}, cont_c)
	copy4 := CopyForm(input4)
	copyWithType4 := copy4.(*SelectForm)
	copyWithType4.to_c.Ident = "to_c"
	assertEqual(t, input4, copy4)
	copyWithType4.to_c.Ident = "to_c_edited"
	assertNotEqual(t, input4, copy4)

	// New
	input5 := NewNew(cont_c, end, end)
	copy5 := CopyForm(input5)
	copyWithType5 := copy5.(*NewForm)
	copyWithType5.continuation_c.Ident = "cont_c"
	assertEqual(t, input5, copy5)
	copyWithType5.continuation_c.Ident = "cont_c_edited"
	assertNotEqual(t, input5, copy5)

	// Close
	input6 := NewClose(from_c)
	copy6 := CopyForm(input6)
	copyWithType6 := copy6.(*CloseForm)
	copyWithType6.from_c.Ident = "from_c"
	assertEqual(t, input6, copy6)
	copyWithType6.from_c.Ident = "from_c_edited"
	assertNotEqual(t, input6, copy6)

	// Forward
	input7 := NewForward(to_c, from_c)
	copy7 := CopyForm(input7)
	copyWithType7 := copy7.(*ForwardForm)
	copyWithType7.from_c.Ident = "from_c"
	assertEqual(t, input7, copy7)
	copyWithType7.from_c.Ident = "from_c_edited"
	assertNotEqual(t, input7, copy7)

	// Split
	input8 := NewSplit(pay_c, cont_c, from_c, end)
	copy8 := CopyForm(input8)
	copyWithType8 := copy8.(*SplitForm)
	copyWithType8.from_c.Ident = "from_c"
	assertEqual(t, input8, copy8)
	copyWithType8.from_c.Ident = "from_c_edited"
	assertNotEqual(t, input8, copy8)

	// Call
	input9 := NewCall("func_name", []Name{from_c})
	copy9 := CopyForm(input9)
	copyWithType9 := copy9.(*CallForm)
	copyWithType9.parameters[0].Ident = "from_c"
	assertEqual(t, input9, copy9)
	copyWithType9.functionName = "changed_function_name"
	assertNotEqual(t, input9, copy9)
	copyWithType9.functionName = "func_name"
	assertEqual(t, input9, copy9)
	copyWithType9.parameters[0].Ident = "from_c2"
	assertNotEqual(t, input9, copy9)

}

func TestFreeNames(t *testing.T) {
	to_c := Name{Ident: "to_c", IsSelf: false}
	pay_c := Name{Ident: "pay_c", IsSelf: false}
	cont_c := Name{Ident: "cont_c", IsSelf: false}
	from_c := Name{Ident: "from_c", IsSelf: false}
	new_from_c := Name{Ident: "new_from_c", IsSelf: false}
	other_c := Name{Ident: "other_c", IsSelf: false}
	self := Name{IsSelf: true}
	end := NewClose(self)

	// Send
	input := NewSend(to_c, pay_c, cont_c)
	assertEqualNames(t, input.FreeNames(), []Name{to_c, pay_c, cont_c})

	// Receive
	input2 := NewReceive(pay_c, cont_c, from_c, input)
	assertEqualNames(t, input2.FreeNames(), []Name{from_c, to_c})

	// Case
	input3 := NewCase(from_c, []*BranchForm{NewBranch(Label{L: "labell"}, pay_c, input)})
	assertEqualNames(t, input3.FreeNames(), []Name{from_c, to_c, cont_c})

	input3other := NewCase(from_c, []*BranchForm{NewBranch(Label{L: "labell"}, pay_c, end)})
	assertEqualNames(t, input3other.FreeNames(), []Name{from_c})

	otherName := Name{Ident: "other_var", IsSelf: false}
	input3branches := NewCase(from_c, []*BranchForm{NewBranch(Label{L: "labell"}, pay_c, end), NewBranch(Label{L: "otherbranch"}, pay_c, NewClose(otherName))})
	assertEqualNames(t, input3branches.FreeNames(), []Name{from_c, otherName})

	input3copy := CopyForm(input3branches)
	input3copy.Substitute(from_c, new_from_c)
	input3copy.Substitute(new_from_c, other_c)
	assertEqualNames(t, input3copy.FreeNames(), []Name{other_c, otherName})

	// Select
	input4 := NewSelect(to_c, Label{L: "label1"}, cont_c)
	assertEqualNames(t, input4.FreeNames(), []Name{to_c, cont_c})

	// New
	input5 := NewNew(cont_c, end, end)
	assertEqualNames(t, input5.FreeNames(), []Name{})

	input5other := NewNew(cont_c, input3, end)
	assertEqualNames(t, input5other.FreeNames(), []Name{from_c, to_c})

	// Close
	input6 := NewClose(from_c)
	assertEqualNames(t, input6.FreeNames(), []Name{from_c})

	// Forward
	input7 := NewForward(to_c, from_c)
	assertEqualNames(t, input7.FreeNames(), []Name{to_c, from_c})

	input7other := NewForward(self, from_c)
	assertEqualNames(t, input7other.FreeNames(), []Name{from_c})

	// Split
	input8 := NewSplit(pay_c, cont_c, from_c, end)
	assertEqualNames(t, input8.FreeNames(), []Name{from_c})

	// Split
	input9 := NewCall("func_name", []Name{from_c})
	assertEqualNames(t, input9.FreeNames(), []Name{from_c})
}

func TestFormHasContinuation(t *testing.T) {
	expectedFalse, expectedTrue := []bool{}, []bool{}
	to_c := Name{Ident: "to_c"}
	pay_c := Name{Ident: "pay_c"}
	cont_c := Name{Ident: "cont_c"}
	from_c := Name{Ident: "from_c"}
	end := NewClose(Name{Ident: "self"})

	// Send
	input := NewSend(to_c, pay_c, cont_c)
	expectedFalse = append(expectedFalse, FormHasContinuation(input))

	// Select
	input2 := NewSelect(to_c, Label{L: "label1"}, cont_c)
	expectedFalse = append(expectedFalse, FormHasContinuation(input2))

	// Close
	input3 := NewClose(from_c)
	expectedFalse = append(expectedFalse, FormHasContinuation(input3))

	// Forward
	input4 := NewForward(to_c, from_c)
	expectedFalse = append(expectedFalse, FormHasContinuation(input4))

	// Call
	input5 := NewCall("f", []Name{to_c})
	expectedFalse = append(expectedFalse, FormHasContinuation(input5))

	// Cast
	input6 := NewCast(to_c, cont_c)
	expectedFalse = append(expectedFalse, FormHasContinuation(input6))

	// case *DropForm:

	// Receive
	input7 := NewReceive(pay_c, cont_c, from_c, end)
	expectedTrue = append(expectedTrue, FormHasContinuation(input7))

	// Case
	input8 := NewCase(from_c, []*BranchForm{NewBranch(Label{L: "label1"}, pay_c, end)})
	expectedTrue = append(expectedTrue, FormHasContinuation(input8))

	// New
	input9 := NewNew(cont_c, end, end)
	expectedTrue = append(expectedTrue, FormHasContinuation(input9))

	// Split
	input10 := NewSplit(pay_c, cont_c, from_c, end)
	expectedTrue = append(expectedTrue, FormHasContinuation(input10))

	// Wait
	input11 := NewWait(to_c, end)
	expectedTrue = append(expectedTrue, FormHasContinuation(input11))

	// Drop
	input12 := NewDrop(to_c, end)
	expectedTrue = append(expectedTrue, FormHasContinuation(input12))

	for i := range expectedFalse {
		if expectedFalse[i] {
			t.Errorf("expected FormHasContinuation to return false for case %d but found true", i)
		}
	}

	for i := range expectedTrue {
		if !expectedTrue[i] {
			t.Errorf("expected FormHasContinuation to return true for case %d but found false", i)
		}
	}
}
