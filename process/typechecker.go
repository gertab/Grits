package process

import (
	"bytes"
	"fmt"
	"phi/types"
)

func Typecheck(processes []*Process, processesFreeNames [][]Name, globalEnv *GlobalEnvironment) error {
	errorChan := make(chan error)
	doneChan := make(chan bool)

	fmt.Println("Initiating typechecking")

	// Running in a separate process allows us to break the typechecking part as soon as the first
	// error is found
	go typecheckFunctionsAndProcesses(processes, processesFreeNames, globalEnv, errorChan, doneChan)

	select {
	case err := <-errorChan:
		return err
	case <-doneChan:
		fmt.Println("Typecheck successful")
	}

	return nil
}

func typecheckFunctionsAndProcesses(processes []*Process, processesFreeNames [][]Name, globalEnv *GlobalEnvironment, errorChan chan error, doneChan chan bool) {
	defer func() {
		// todo replace with WG
		doneChan <- true
	}()

	// Start with some preliminary check on the types
	err := preliminaryChecks(processes, processesFreeNames, globalEnv)
	if err != nil {
		errorChan <- err
	}

	// We can initiate the more heavyweight typechecking
	// 1) Typecheck function definitions
	err = typecheckFunctionDefinitions(globalEnv)
	if err != nil {
		errorChan <- err
	}
	// 2) todo Process definitions

	// todo add checks to make sure that all types are there

}

// Perform some preliminary checks about the type definitions
// Ensures that types only referred to existing labelled types (i.e. recursion is used correctly). Also, ensures that there are no missing types and that types are well formed
func preliminaryChecks(processes []*Process, processesFreeNames [][]Name, globalEnv *GlobalEnvironment) error {

	// First analyse the labelled types (i.e. type A = ...)
	err := types.SanityChecksTypeDefinitions(*globalEnv.Types)

	if err != nil {
		return err
	}

	// Analyse the function declarations types (i.e. let f(x : B) : A = ...)
	// todo make sure that each function has a unique name
	var all_types []types.SessionType
	for _, f := range *globalEnv.FunctionDefinitions {
		if f.Type != nil {
			all_types = append(all_types, f.Type)
		} else {
			return fmt.Errorf("function %s has a missing type of provider", f.String())
		}

		for _, p := range f.Parameters {
			if p.Type != nil {
				all_types = append(all_types, p.Type)
			} else {
				return fmt.Errorf("in function definition %s, parameter %s has a missing type", f.String(), p.String())
			}
		}
	}

	// Check the types for the processes (i.e. prc[a] : A = ...)
	for i := range processes {
		if processes[i].Type != nil {
			all_types = append(all_types, processes[i].Type)
		} else {
			return fmt.Errorf("process %s has a missing type of provider", processes[i].OutlineString())
		}

		// Check also that the types annotated for each process are correct
		// -> prc[a] : A = ... % b : A1, c : A2, ...
		actualFreeNames := processes[i].Body.FreeNames()
		if len(actualFreeNames) != len(processesFreeNames[i]) {
			return fmt.Errorf("Process %s has %d free names (i.e. %s), but there were %d found (%s)", processes[i].String(), len(actualFreeNames), NamesToString(actualFreeNames), len(processesFreeNames[i]), NamesToString(processesFreeNames[i]))
		}

		// todo check that processesFreeNames[i] are all unique
		// todo check that all names in processesFreeNames[i] match with the ones in actualFreeNames

	}

	err = types.SanityChecksType(all_types, *globalEnv.Types)
	if err != nil {
		return err
	}

	return nil
}

func typecheckFunctionDefinitions(globalEnv *GlobalEnvironment) error {
	labelledTypesEnv := types.ProduceLabelledSessionTypeEnvironment(*globalEnv.Types)
	functionDefinitionsEnv := produceFunctionDefinitionsEnvironment(*globalEnv.FunctionDefinitions)

	for _, funcDef := range *globalEnv.FunctionDefinitions {
		gammaNameTypesCtx := produceNameTypesCtx(funcDef.Parameters)
		providerType := funcDef.Type

		err := funcDef.Body.typecheckForm(gammaNameTypesCtx, nil, providerType, labelledTypesEnv, functionDefinitionsEnv)
		if err != nil {
			return err
		}
	}

	return nil
}

// type LabelledTypesEnv types.LabelledTypesEnv
type FunctionTypesEnv map[string]FunctionType /* represented as sigma in the type system */
type NamesTypesCtx map[string]NamesType       /* maps names to their types */

////////////////////////////////////////////////////////////////
///////////////// Syntax directed typechecking /////////////////
////////////////////////////////////////////////////////////////

// Each form has a dedicated typechecking function

// typecheckForm uses these parameters:
// -> gammaNameTypesCtx NamesTypesCtx   		<- names in context to be used (in case of linearity, ...)
// -> providerShadowName *Name          		<- name of the process providing on (nil when name 'self' is used instead)
// -> providerType types.SessionType    		<- the type of the provider (i.e. type of provider name 'self')
// -> labelledTypesEnv types.LabelledTypesEnv 	<- [read-only] keeps the mapping of pre-defined types (type A = ...)
// -> sigma FunctionTypesEnv           	 		<- [read-only] keeps the mapping of pre-defined function definitions (let f() : A = ...)

// */-*: send w<u, v>
func (p *SendForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	if isProvider(p.to_c, providerShadowName) {
		// MulR: *
		logRule("rule MulR")

		providerType = types.Unfold(providerType, labelledTypesEnv)
		// The type of the provider must be SendType
		providerSendType, sendTypeOk := providerType.(*types.SendType)

		if !sendTypeOk {
			// wrong type: A * B
			return fmt.Errorf("expected '%s' to have a send type (A * B), but found type '%s' instead", p.String(), providerType.String())
		}

		expectedLeftType := providerSendType.Left
		expectedRightType := providerSendType.Right
		foundLeftType, errorLeft := consumeName(p.payload_c, gammaNameTypesCtx)
		foundLeftType = types.Unfold(foundLeftType, labelledTypesEnv)
		foundRightType, errorRight := consumeName(p.continuation_c, gammaNameTypesCtx)
		foundRightType = types.Unfold(foundRightType, labelledTypesEnv)

		if errorLeft != nil {
			return errorLeft
		}

		if errorRight != nil {
			return errorRight
		}

		// The expected and found types must match
		if !types.EqualType(expectedLeftType, foundLeftType, labelledTypesEnv) {
			return fmt.Errorf("expected type of '%s' to be '%s', but found type '%s' instead", p.payload_c.String(), expectedLeftType.String(), foundLeftType.String())
		}

		if !types.EqualType(expectedRightType, foundRightType, labelledTypesEnv) {
			return fmt.Errorf("expected type of '%s' to be '%s', but found type '%s' instead", p.continuation_c.String(), expectedRightType.String(), foundRightType.String())
		}

		// Set the types for the names
		p.to_c.Type = providerSendType
		p.payload_c.Type = foundLeftType
		p.continuation_c.Type = foundRightType

	} else if isProvider(p.continuation_c, providerShadowName) {
		// ImpL: -*
		logRule("rule ImpL")

		clientType, errorClient := consumeName(p.to_c, gammaNameTypesCtx)
		if errorClient != nil {
			return errorClient
		}

		clientType = types.Unfold(clientType, labelledTypesEnv)
		// The type of the client must be ReceiveType
		clientReceiveType, clientTypeOk := clientType.(*types.ReceiveType)

		if !clientTypeOk {
			// wrong type: A -* B
			return fmt.Errorf("expected '%s' to have a receive type (A -* B), but found type '%s' instead", p.to_c.String(), clientType.String())
		}

		expectedLeftType := clientReceiveType.Left
		expectedRightType := clientReceiveType.Right
		foundLeftType, errorLeft := consumeName(p.payload_c, gammaNameTypesCtx)
		foundRightType, errorRight := consumeNameMaybeSelf(p.continuation_c, providerShadowName, gammaNameTypesCtx, providerType)

		expectedLeftType = types.Unfold(expectedLeftType, labelledTypesEnv)
		expectedRightType = types.Unfold(expectedRightType, labelledTypesEnv)
		foundLeftType = types.Unfold(foundLeftType, labelledTypesEnv)
		foundRightType = types.Unfold(foundRightType, labelledTypesEnv)

		if errorLeft != nil {
			return errorLeft
		}

		if errorRight != nil {
			return errorRight
		}

		// The expected and found types must match
		if !types.EqualType(expectedLeftType, foundLeftType, labelledTypesEnv) {
			return fmt.Errorf("expected type of '%s' to be '%s', but found type '%s' instead", p.payload_c.String(), expectedLeftType.String(), foundLeftType.String())
		}

		if !types.EqualType(expectedRightType, foundRightType, labelledTypesEnv) {
			return fmt.Errorf("expected type of '%s' to be '%s', but found type '%s' instead", p.continuation_c.String(), expectedRightType.String(), foundRightType.String())
		}

		// Set the types for the names
		p.to_c.Type = clientReceiveType
		p.payload_c.Type = foundLeftType
		p.continuation_c.Type = foundRightType
	} else if isProvider(p.payload_c, providerShadowName) {
		return fmt.Errorf("the send construct requires that you use the self name or send self as a continuation. In '%s', self was used as a payload", p.String())
	} else {
		return fmt.Errorf("the send construct requires that you use the self name or send self as a continuation. In '%s', self was not used appropriately", p.String())
	}

	// make sure that no variables are left in gamma
	err := linearGammaContext(gammaNameTypesCtx)

	return err
}

// */-*: <x, y> <- recv w; P
func (p *ReceiveForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	if isProvider(p.from_c, providerShadowName) {
		// ImpR: -*
		logRule("rule ImpR")

		providerType = types.Unfold(providerType, labelledTypesEnv)
		// The type of the provider must be ReceiveType
		providerReceiveType, receiveTypeOk := providerType.(*types.ReceiveType)

		if !receiveTypeOk {
			// wrong type: A -* B
			return fmt.Errorf("expected '%s' to have a receive type (A -* B), but found type '%s' instead", p.String(), providerType.String())
		}

		newLeftType := providerReceiveType.Left
		newLeftType = types.Unfold(newLeftType, labelledTypesEnv)
		newRightType := providerReceiveType.Right
		newRightType = types.Unfold(newRightType, labelledTypesEnv)

		if nameTypeExists(gammaNameTypesCtx, p.payload_c.Ident) ||
			nameTypeExists(gammaNameTypesCtx, p.continuation_c.Ident) {
			// Names are not fresh [todo check if needed]
			return fmt.Errorf("variable names <%s, %s> already defined. Use unique names in %s", p.payload_c.String(), p.continuation_c.String(), p.String())
		}

		// todo maybe remove check
		if isProvider(p.payload_c, providerShadowName) ||
			isProvider(p.continuation_c, providerShadowName) {
			// Unwanted reference to self
			return fmt.Errorf("variable names <%s, %s> should not refer to self", p.payload_c.String(), p.continuation_c.String())
		}

		gammaNameTypesCtx[p.payload_c.Ident] = NamesType{Type: newLeftType}

		p.from_c.Type = providerReceiveType
		p.payload_c.Type = newLeftType
		p.continuation_c.Type = newRightType

		continuationError := p.continuation_e.typecheckForm(gammaNameTypesCtx, &p.continuation_c, newRightType, labelledTypesEnv, sigma)

		return continuationError
	} else if isProvider(p.payload_c, providerShadowName) || isProvider(p.continuation_c, providerShadowName) {
		// providerType = types.Unfold(providerType, labelledTypesEnv)
		// _, receiveTypeOk := providerType.(*types.ReceiveType)

		// if receiveTypeOk {
		// 	} else {
		// todo check further type info
		// 	}
		return fmt.Errorf("you cannot assign self to a new channel (%s)", p.String())
	} else {
		// MulL: *
		logRule("rule MulL")

		clientType, errorClient := consumeName(p.from_c, gammaNameTypesCtx)
		if errorClient != nil {
			return errorClient
		}

		clientType = types.Unfold(clientType, labelledTypesEnv)
		// The type of the client must be SendType
		clientSendType, clientTypeOk := clientType.(*types.SendType)

		if !clientTypeOk {
			// wrong type, expected A * B
			return fmt.Errorf("expected '%s' to have a send type (A * B), but found type '%s' instead", p.from_c.String(), clientType.String())
		}

		newLeftType := clientSendType.Left
		newRightType := clientSendType.Right

		if nameTypeExists(gammaNameTypesCtx, p.payload_c.Ident) ||
			nameTypeExists(gammaNameTypesCtx, p.continuation_c.Ident) {
			// Names are not fresh [todo check if needed]
			return fmt.Errorf("variable names <%s, %s> already defined. Use unique names", p.payload_c.String(), p.continuation_c.String())
		}

		if isProvider(p.payload_c, providerShadowName) ||
			isProvider(p.continuation_c, providerShadowName) {
			// Unwanted reference to self
			return fmt.Errorf("variable names <%s, %s> should not refer to self", p.payload_c.String(), p.continuation_c.String())
		}

		gammaNameTypesCtx[p.payload_c.Ident] = NamesType{Type: newLeftType}
		gammaNameTypesCtx[p.continuation_c.Ident] = NamesType{Type: newRightType}

		p.from_c.Type = clientSendType
		p.payload_c.Type = newLeftType
		p.continuation_c.Type = newRightType

		continuationError := p.continuation_e.typecheckForm(gammaNameTypesCtx, providerShadowName, providerType, labelledTypesEnv, sigma)

		return continuationError
	}
}

// Internal/External Choice: w.l<u>
func (p *SelectForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	if isProvider(p.to_c, providerShadowName) {
		// IChoiceR: +{label1: T1, ...}
		logRule("rule IChoiceR")

		providerType = types.Unfold(providerType, labelledTypesEnv)
		// The type of the provider must be SelectLabelType
		providerSelectLabelType, selectLabelTypeOk := providerType.(*types.SelectLabelType)

		if !selectLabelTypeOk {
			// wrong type, expected +{...}
			return fmt.Errorf("expected '%s' to have a select type (+{...}), but found type '%s' instead", p.String(), providerType.String())
		}

		// Match branch by label
		continuationType, continuationTypeFound := types.FetchSelectBranch(providerSelectLabelType.Branches, p.label.L)

		if continuationTypeFound {
			foundContinuationType, errorContinuationType := consumeName(p.continuation_c, gammaNameTypesCtx)

			if errorContinuationType != nil {
				return errorContinuationType
			}

			if !types.EqualType(continuationType, foundContinuationType, labelledTypesEnv) {
				return fmt.Errorf("type of '%s' is '%s'. Expected type to be '%s'", p.continuation_c.String(), foundContinuationType.String(), continuationType.String())
			}

			p.to_c.Type = providerSelectLabelType
			p.continuation_c.Type = continuationType
			// return nil
		} else {
			return fmt.Errorf("could not match label '%s' (from '%s') with the labels from the type '%s'", p.label.String(), p.String(), providerSelectLabelType.String())
		}
	} else if isProvider(p.continuation_c, providerShadowName) {
		// EChoiceL: &{label1: T1, ...}
		logRule("rule EChoiceL")

		clientType, errorClient := consumeName(p.to_c, gammaNameTypesCtx)
		if errorClient != nil {
			return errorClient
		}

		clientType = types.Unfold(clientType, labelledTypesEnv)
		// The type of the client must be BranchCaseType
		clientBranchCaseType, clientTypeOk := clientType.(*types.BranchCaseType)

		if !clientTypeOk {
			// wrong type, expected &{...}
			return fmt.Errorf("expected '%s' to have a branching type (&{...}), but found type '%s' instead", p.String(), clientType.String())
		}

		// Match branch by label
		continuationType, continuationTypeFound := types.FetchSelectBranch(clientBranchCaseType.Branches, p.label.L)

		if continuationTypeFound {
			foundContinuationType, errorContinuationType := consumeNameMaybeSelf(p.continuation_c, providerShadowName, gammaNameTypesCtx, providerType)

			if errorContinuationType != nil {
				return errorContinuationType
			}

			if !types.EqualType(continuationType, foundContinuationType, labelledTypesEnv) {
				return fmt.Errorf("type of '%s' is '%s'. Expected type to be '%s'", p.continuation_c.String(), foundContinuationType.String(), continuationType.String())
			}

			// Type ok

			// Set types
			p.to_c.Type = clientBranchCaseType
			p.continuation_c.Type = continuationType
		} else {
			return fmt.Errorf("could not match label '%s' (from '%s') with the labels from the type '%s'", p.label.String(), p.String(), clientBranchCaseType.String())
		}
	} else {
		return fmt.Errorf("expected '%s' to either receive or send label on 'self', e.g. self.%s<%s> or %s.%s<self>", p.String(), p.label.String(), p.to_c.String(), p.continuation_c.String(), p.label.String())
	}

	// make sure that no variables are left in gamma
	err := linearGammaContext(gammaNameTypesCtx)

	return err
}

// Case: case from_c ( branches )
func (p *CaseForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	if isProvider(p.from_c, providerShadowName) {
		// EChoiceR: &{label1: T1, ...}
		logRule("rule EChoiceR")

		providerType = types.Unfold(providerType, labelledTypesEnv)
		// The type of the provider must be BranchCaseType
		providerBranchCaseType, branchCaseTypeOk := providerType.(*types.BranchCaseType)

		if !branchCaseTypeOk {
			// wrong type, expected &{...}
			return fmt.Errorf("expected '%s' to have a branching type (&{...}), but found type '%s' instead", p.String(), providerType.String())
		}

		labelsChecked := []string{}

		// Check each branch individually
		for _, curBranchForm := range p.branches {
			// Pick each branch ast and match it with its type using the label
			expectedBranchType, typeFound := types.LookupBranchByLabel(providerBranchCaseType.Branches, curBranchForm.label.L)

			labelsChecked = append(labelsChecked, curBranchForm.label.L)

			if !typeFound {
				return fmt.Errorf("branch labelled '%s' does not match the branches of type '%s'", curBranchForm.String(), providerBranchCaseType.String())
			}

			// Set type
			curBranchForm.payload_c.Type = expectedBranchType.Session_type

			continuationError := curBranchForm.continuation_e.typecheckForm(gammaNameTypesCtx, &curBranchForm.payload_c, expectedBranchType.Session_type, labelledTypesEnv, sigma)

			if continuationError != nil {
				return continuationError
			}
		}

		if len(labelsChecked) < len(providerBranchCaseType.Branches) {
			labels := extractUnusedLabels(providerBranchCaseType.Branches, labelsChecked)

			return fmt.Errorf("some labels (i.e. %s) from the type '%s' are not pattern matched in the case construct: %s", labels, providerBranchCaseType.String(), p.String())
		}

		// Set type of case
		p.from_c.Type = providerBranchCaseType
	} else {
		// IChoiceL: +{label1: T1, ...}
		logRule("rule IChoiceL")

		clientType, errorClient := consumeName(p.from_c, gammaNameTypesCtx)
		if errorClient != nil {
			return errorClient
		}

		clientType = types.Unfold(clientType, labelledTypesEnv)
		// The type of the client must be SelectLabelType
		clientSelectLabelType, clientTypeOk := clientType.(*types.SelectLabelType)

		if !clientTypeOk {
			// wrong type, expected +{...}
			return fmt.Errorf("expected '%s' to have a select type (+{...}), but found type '%s' instead", p.String(), clientType.String())
		}

		labelsChecked := []string{}

		// Check each branch individually
		for _, curBranchForm := range p.branches {
			// Pick each branch ast and match it with its type using the label
			expectedBranchType, typeFound := types.LookupBranchByLabel(clientSelectLabelType.Branches, curBranchForm.label.L)

			labelsChecked = append(labelsChecked, curBranchForm.label.L)

			if !typeFound {
				return fmt.Errorf("case labelled '%s' does not match the branches of type '%s'", curBranchForm.String(), clientSelectLabelType.String())
			}

			// curBranchForm.payload_c cannot exist in gammaNameTypesCtx
			gammaNameTypesCtx[curBranchForm.payload_c.Ident] = NamesType{Type: expectedBranchType.Session_type}

			// Set type
			curBranchForm.payload_c.Type = expectedBranchType.Session_type

			continuationError := curBranchForm.continuation_e.typecheckForm(gammaNameTypesCtx, providerShadowName, providerType, labelledTypesEnv, sigma)

			if continuationError != nil {
				return continuationError
			}
		}

		if len(labelsChecked) < len(clientSelectLabelType.Branches) {
			labels := extractUnusedLabels(clientSelectLabelType.Branches, labelsChecked)

			return fmt.Errorf("some labels (i.e. %s) from the type '%s' are not pattern matched in the case construct: %s", labels, clientSelectLabelType.String(), p.String())
		}

		// Set type of case
		p.from_c.Type = clientSelectLabelType
	}

	return nil
}

// Branch: label<payload_c> => continuation_e
func (p *BranchForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	panic("this should never be called directly")
	// return nil
}

func (p *NewForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	return nil
}

// 1 : close w
func (p *CloseForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	// EndR: 1
	logRule("rule EndR")
	providerType = types.Unfold(providerType, labelledTypesEnv)

	if isProvider(p.from_c, providerShadowName) {
		providerUnitType, unitTypeOk := providerType.(*types.UnitType)

		if !unitTypeOk {
			return fmt.Errorf("expected '%s' to have a unit type (1), but found type '%s' instead", p.String(), providerType.String())
		}

		p.from_c.Type = providerUnitType
	} else {
		// Closing on the wrong name
		_, unitTypeOk := providerType.(*types.UnitType)

		if unitTypeOk {
			return fmt.Errorf("expected '%s' to have a send on 'self' instead", p.String())
		} else {
			return fmt.Errorf("expected '%s' to have a unit type (1), but found type '%s' instead", p.String(), providerType.String())
		}
	}

	// make sure that no variables are left in gamma
	err := linearGammaContext(gammaNameTypesCtx)

	return err
}

// 1 : wait w; ...
func (p *WaitForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	// EndL: 1
	logRule("rule EndL")

	// Can only wait for a client (not self)
	if !isProvider(p.to_c, providerShadowName) {
		clientType, errorClient := consumeName(p.to_c, gammaNameTypesCtx)
		if errorClient != nil {
			return errorClient
		}

		clientType = types.Unfold(clientType, labelledTypesEnv)
		// The type of the client must be UnitType
		clientUnitType, clientTypeOk := clientType.(*types.UnitType)

		if clientTypeOk {
			// Set type
			p.to_c.Type = clientUnitType

			// Continue checking the remaining process
			continuationError := p.continuation_e.typecheckForm(gammaNameTypesCtx, providerShadowName, providerType, labelledTypesEnv, sigma)

			return continuationError
		} else {
			return fmt.Errorf("expected '%s' to have a unit type (1), but found type '%s' instead", p.to_c.String(), clientUnitType.String())
		}
	} else {
		// Waiting on the wrong name
		providerType = types.Unfold(providerType, labelledTypesEnv)
		_, unitTypeOk := providerType.(*types.UnitType)

		if unitTypeOk {
			return fmt.Errorf("expected '%s' to have a wait on a 'non-self' channel instead (%s is acting as self)", p.String(), p.to_c.String())
		} else {
			return fmt.Errorf("expected '%s' to have a unit type (1), but found type '%s' instead", p.String(), providerType.String())
		}
	}
}

// fwd w u
func (p *ForwardForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	// ID: 1
	logRule("rule ID")

	if isProvider(p.from_c, providerShadowName) {
		return fmt.Errorf("forwarding to self (%s) is not allowed. Use 'fwd %s %s' instead)", p.String(), p.from_c.String(), p.to_c.String())
	}

	if !isProvider(p.to_c, providerShadowName) {
		return fmt.Errorf("not forwarding on self (%s). Expected forward to refer to self (fwd %s %s)", p.String(), p.from_c.String(), p.to_c.String())
	}

	clientType, errorClient := consumeName(p.from_c, gammaNameTypesCtx)
	if errorClient != nil {
		return errorClient
	}

	if !types.EqualType(providerType, clientType, labelledTypesEnv) {
		return fmt.Errorf("problem in %s. Type of %s (%s) and %s (%s) do do not match", p.String(), p.to_c.String(), providerType.String(), p.from_c.String(), clientType.String())
	}

	// Set types
	p.to_c.Type = providerType
	p.from_c.Type = clientType

	// make sure that no variables are left in gamma
	err := linearGammaContext(gammaNameTypesCtx)

	return err
}

// drop w; ...
func (p *DropForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	// Drop
	logRule("rule Drp")

	// Can only wait for a client (not self)
	if !isProvider(p.client_c, providerShadowName) {
		clientType, errorClient := consumeName(p.client_c, gammaNameTypesCtx)
		if errorClient != nil {
			return errorClient
		}

		if types.IsWeakenable(clientType) {
			// Set type
			p.client_c.Type = clientType

			// Continue checking the remaining process
			continuationError := p.continuation_e.typecheckForm(gammaNameTypesCtx, providerShadowName, providerType, labelledTypesEnv, sigma)

			return continuationError
		} else {
			return fmt.Errorf("expected '%s' to have a weakenable type, but found the non-weakenable type '%s' instead", p.client_c.String(), clientType.String())
		}
	} else {
		// Wrongly dropping self
		return fmt.Errorf("expected '%s' to have a drop on a 'non-self' channel instead (i.e. incorrectly dropping '%s')", p.String(), p.client_c.String())
	}
}

func (p *CallForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {

	// Check that function exists
	functionSignature, exists := sigma[p.functionName]
	if !exists {
		return fmt.Errorf("function '%s' is undefined", p.String())
	}

	// Check that the arity matches
	if len(functionSignature.Parameters)+1 == len(p.parameters) {
		// Parameters being passed include a reference to 'self' as the first element

		foundProviderType, errorProvider := consumeNameMaybeSelf(p.parameters[0], providerShadowName, gammaNameTypesCtx, providerType)
		if errorProvider != nil {
			return errorProvider
		}

		// Check type of self
		if !types.EqualType(foundProviderType, functionSignature.Type, labelledTypesEnv) {
			return fmt.Errorf("type error in function call '%s'. Name '%s' has type '%s', but expected '%s'", p.String(), p.parameters[0].String(), foundProviderType.String(), functionSignature.Type.String())
		}

		// Check types of each parameter
		for i := 1; i < len(p.parameters); i++ {

			foundParamType, paramTypeError := consumeName(p.parameters[i], gammaNameTypesCtx)

			if paramTypeError != nil {
				return paramTypeError
			}

			expectedType := functionSignature.Parameters[i-1].Type

			if !types.EqualType(foundParamType, expectedType, labelledTypesEnv) {
				return fmt.Errorf("type error in function call '%s'. Name '%s' has type '%s', but expected '%s'", p.String(), p.parameters[i].String(), foundParamType.String(), expectedType.String())
			}

			// Set types
			p.parameters[i].Type = foundParamType
		}
	} else if len(functionSignature.Parameters) == len(p.parameters) {
		// 'self' is not included in the parameters

		// Check type of self
		if !types.EqualType(providerType, functionSignature.Type, labelledTypesEnv) {
			providerName := "self"
			if providerShadowName != nil {
				providerName = providerShadowName.String()
			}

			return fmt.Errorf("type error in function call '%s'. Provider '%s' has type '%s', but %s expects '%s'", p.String(), providerName, providerType.String(), p.functionName, functionSignature.Type.String())
		}

		// Check types of each parameter
		for i := 0; i < len(p.parameters); i++ {
			foundParamType, paramTypeError := consumeName(p.parameters[i], gammaNameTypesCtx)

			if paramTypeError != nil {
				return paramTypeError
			}

			expectedType := functionSignature.Parameters[i].Type

			if !types.EqualType(foundParamType, expectedType, labelledTypesEnv) {
				return fmt.Errorf("type error in function call '%s'. Name '%s' has type '%s', but expected '%s'", p.String(), p.parameters[i].String(), foundParamType.String(), expectedType.String())
			}

			// Set types
			p.parameters[i].Type = foundParamType
		}
	} else {
		// Wrong number of parameters
		return fmt.Errorf("wrong number of parameters in function call '%s'. Expected %d, but found %d parameters", p.String(), len(functionSignature.Parameters), len(p.parameters))
	}

	// Set type
	p.ProviderType = functionSignature.Type

	// make sure that no variables are left in gamma
	err := linearGammaContext(gammaNameTypesCtx)

	return err
}

func (p *SplitForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	return nil
}
func (p *CastForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	return nil
}
func (p *ShiftForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	return nil
}
func (p *PrintForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	return nil
}

// /// Fixed Environments: Set once and only read from

type FunctionType struct {
	FunctionName string
	Parameters   []Name
	Type         types.SessionType
}

func produceFunctionDefinitionsEnvironment(functionDefs []FunctionDefinition) FunctionTypesEnv {
	functionTypesEnv := make(FunctionTypesEnv)
	for _, j := range functionDefs {
		functionTypesEnv[j.FunctionName] = FunctionType{Type: j.Type, FunctionName: j.FunctionName, Parameters: j.Parameters}
	}

	return functionTypesEnv
}

type NamesType struct {
	Name Name
	Type types.SessionType
}

func produceNameTypesCtx(names []Name) NamesTypesCtx {
	namesTypesCtx := make(NamesTypesCtx)
	for _, j := range names {
		namesTypesCtx[j.Ident] = NamesType{Type: j.Type, Name: j}
	}

	return namesTypesCtx
}

func nameTypeExists(namesTypesCtx NamesTypesCtx, key string) bool {
	_, ok := namesTypesCtx[key]
	return ok
}

// /// Util functions

// we check whether a channel is the provider (i.e. either self of same as the explicit provider name)
// E.g. the channel being sent to is the provider for these cases:
// -> send self<b, c>
// -> send a<b, c> (where a is the provider)
func isProvider(name Name, providerShadowName *Name) bool {
	if name.IsSelf {
		return true
	}

	if providerShadowName != nil {
		return name.Ident == providerShadowName.Ident
	}

	return false
}

// Takes a name from gamma. If the name is not found, then return error
func consumeName(name Name, gammaNameTypesCtx NamesTypesCtx) (types.SessionType, error) {

	if name.IsSelf {
		return nil, fmt.Errorf("found self, expected a client")
	}

	foundName, ok := gammaNameTypesCtx[name.Ident]

	if ok {
		// If linear then remove
		delete(gammaNameTypesCtx, name.Ident)

		return foundName.Type, nil
	}

	// Problem since the requested name was not found in the gamma
	return nil, fmt.Errorf("problem since the requested name (%s) was not found in gamma. todo set cool error message", name.String())
}

// Takes a name from gamma. If the name is 'self', then return the provider type instead of fetching it from gamma
func consumeNameMaybeSelf(name Name, providerShadowName *Name, gammaNameTypesCtx NamesTypesCtx, providerType types.SessionType) (types.SessionType, error) {
	if name.IsSelf {
		return providerType, nil
	}

	if providerShadowName != nil && providerShadowName.Ident == name.Ident {
		return providerType, nil
	}

	foundName, ok := gammaNameTypesCtx[name.Ident]

	if ok {
		// If linear then remove
		delete(gammaNameTypesCtx, name.Ident)

		return foundName.Type, nil
	}

	// Problem since the requested name was not found in the gamma
	return nil, fmt.Errorf("problem since the requested name (%s) was not found in gamma. todo set cool error message", name.String())
}

func stringifyContext(gammaNameTypesCtx NamesTypesCtx) string {
	if len(gammaNameTypesCtx) == 0 {
		return ""
	}

	var buffer bytes.Buffer

	for k := range gammaNameTypesCtx {
		buffer.WriteString(k)
		buffer.WriteString("; ")
	}

	str := buffer.String()

	return str[:len(str)-2]
}

// Enforce linearity, i.e. ensure that there are no variables left in Gamma
func linearGammaContext(gammaNameTypesCtx NamesTypesCtx) error {
	// todo change to allow weakenable variables (although drop prevents this)
	if len(gammaNameTypesCtx) == 1 {
		return fmt.Errorf("linearity requires that no names are left behind, however there is one names (%s) left", stringifyContext(gammaNameTypesCtx))
	} else if len(gammaNameTypesCtx) > 1 {
		return fmt.Errorf("linearity requires that no names are left behind, however there were %d names (%s) left", len(gammaNameTypesCtx), stringifyContext(gammaNameTypesCtx))
	}

	// Ok, no unwanted variables left in gamma
	return nil
}

func extractUnusedLabels(branches []types.BranchOption, labels []string) string {
	// One or more branches are not exhausted
	uncheckedBranches := types.GetUncheckedBranches(branches, labels)

	var labelsNotChecked bytes.Buffer
	for i, j := range uncheckedBranches {
		labelsNotChecked.WriteString(j.Label)

		if i < len(uncheckedBranches)-1 {
			labelsNotChecked.WriteString(", ")
		}
	}

	return labelsNotChecked.String()
}

func logRule(s string) {
	fmt.Println(s)
}
