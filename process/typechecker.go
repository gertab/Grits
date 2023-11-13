package process

import (
	"bytes"
	"fmt"
	"phi/types"
)

func Typecheck(processes []*Process, globalEnv *GlobalEnvironment) error {
	errorChan := make(chan error)
	doneChan := make(chan bool)

	fmt.Println("Initiating typechecking")

	// Running in a separate process allows us to break the typechecking part as soon as the first
	// error is found
	go typecheckFunctionsAndProcesses(processes, globalEnv, errorChan, doneChan)

	select {
	case err := <-errorChan:
		return err
	case <-doneChan:
		fmt.Println("Typecheck successful")
	}

	return nil
}

func typecheckFunctionsAndProcesses(processes []*Process, globalEnv *GlobalEnvironment, errorChan chan error, doneChan chan bool) {
	defer func() {
		doneChan <- true
	}()

	// Start with some preliminary check on the types
	err := lightweightChecks(processes, globalEnv)
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
}

// Perform some preliminary checks about the type type definitions
// Ensures that types only referred to existing labelled types
func lightweightChecks(processes []*Process, globalEnv *GlobalEnvironment) error {

	// First check the labelled types (i.e. type A = ...)
	err := types.SanityChecksTypeDefinitions(*globalEnv.Types)

	if err != nil {
		return err
	}

	// Check the types for the function declarations (i.e. let f() : A = ...)
	var all_types []types.SessionType
	for _, i := range *globalEnv.FunctionDefinitions {
		if i.Type != nil {
			all_types = append(all_types, i.Type)
		}

		for _, i := range i.Parameters {
			if i.Type != nil {
				all_types = append(all_types, i.Type)
			}
		}
	}

	// Check the types for the processes (i.e. prc[a] : A = ...)
	for _, i := range processes {
		if i.Type != nil {
			all_types = append(all_types, i.Type)
		}
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
// ... gammaNameTypesCtx NamesTypesCtx   		<- names in context to be used (in case of linearity, ...)
// ... providerShadowName *Name          		<- name of the process providing on (typically nil, since self is used instead)
// ... providerType types.SessionType    		<- the type of the provider (i.e. type of self)
// ... labelledTypesEnv types.LabelledTypesEnv 	<- keeps the mapping of pre-defined types (type A = ...)
// ... sigma FunctionTypesEnv           	 	<- keeps the mapping of pre-defined function definitions (let f() : A = ...)

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
	return nil, fmt.Errorf("problem since the requested name was not found in the gamma. todo set cool error message")
}

func consumeNameMaybeSelf(name Name, gammaNameTypesCtx NamesTypesCtx, providerType types.SessionType) (types.SessionType, error) {
	if name.IsSelf {
		return providerType, nil
	}

	foundName, ok := gammaNameTypesCtx[name.Ident]

	if ok {
		// If linear then remove
		delete(gammaNameTypesCtx, name.Ident)

		return foundName.Type, nil
	}

	// Problem since the requested name was not found in the gamma
	return nil, fmt.Errorf("problem since the requested name was not found in the gamma. todo set cool error message")
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
	if len(gammaNameTypesCtx) > 0 {
		return fmt.Errorf("linearity requires that no names are left behind, however there were %d names (%s) left", len(gammaNameTypesCtx), stringifyContext(gammaNameTypesCtx))
	}

	// Ok, no unwanted variables left in gamma
	return nil
}

// */-o: send w<u, v>
func (p *SendForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	if isProvider(p.to_c, providerShadowName) {
		// MulR: *
		logRule("rule MulR")

		// The type of the provider must be SendType
		providerSendType, sendTypeOk := providerType.(*types.SendType)

		if sendTypeOk {
			expectedLeftType := providerSendType.Left
			expectedRightType := providerSendType.Right
			foundLeftType, errorLeft := consumeName(p.payload_c, gammaNameTypesCtx)
			foundRightType, errorRight := consumeName(p.continuation_c, gammaNameTypesCtx)

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
		} else {
			// wrong type: A * B
			return fmt.Errorf("expected '%s' to have a send type (A * B), but found type '%s' instead", p.String(), providerType.String())
		}
	} else if isProvider(p.continuation_c, providerShadowName) {
		// ImpL: -o
		logRule("rule ImpL")

		clientType, errorClient := consumeName(p.to_c, gammaNameTypesCtx)
		if errorClient != nil {
			return errorClient
		}

		// The type of the client must be ReceiveType
		clientReceiveType, clientTypeOk := clientType.(*types.ReceiveType)

		if clientTypeOk {
			expectedLeftType := clientReceiveType.Left
			expectedRightType := clientReceiveType.Right
			foundLeftType, errorLeft := consumeName(p.payload_c, gammaNameTypesCtx)
			foundRightType, errorRight := consumeNameMaybeSelf(p.continuation_c, gammaNameTypesCtx, providerType)

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
		} else {
			// wrong type: A -o B
			return fmt.Errorf("expected '%s' to have a send type (A -o B), but found type '%s' instead", p.to_c.String(), clientType.String())
		}
	} else if isProvider(p.payload_c, providerShadowName) {
		return fmt.Errorf("the send construct requires that you use the self name or send self as a continuation. In '%s', self was used as a payload", p.String())
	} else {
		return fmt.Errorf("the send construct requires that you use the self name or send self as a continuation. In '%s', self was not used appropriately", p.String())
	}

	// make sure that no variables are left in gamma
	err := linearGammaContext(gammaNameTypesCtx)

	return err
}

// */-o: <x, y> <- recv w; P
func (p *ReceiveForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	if isProvider(p.from_c, providerShadowName) {
		// ImpR: -o
		logRule("rule ImpR")
		// The type of the provider must be ReceiveType
		providerReceiveType, receiveTypeOk := providerType.(*types.ReceiveType)

		if receiveTypeOk {
			newLeftType := providerReceiveType.Left
			newRightType := providerReceiveType.Right

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
			// gammaNameTypesCtx[p.continuation_c.Ident] = NamesType{Type: newRightType}

			p.from_c.Type = providerReceiveType
			p.payload_c.Type = newLeftType
			p.continuation_c.Type = newRightType

			checkContinuation := p.continuation_e.typecheckForm(gammaNameTypesCtx, &p.continuation_c, newRightType, labelledTypesEnv, sigma)

			return checkContinuation
		} else {
			// wrong type: A -o B
			return fmt.Errorf("expected '%s' to have a receive type (A -o B), but found type '%s' instead", p.String(), providerType.String())

		}
	} else if isProvider(p.payload_c, providerShadowName) || isProvider(p.continuation_c, providerShadowName) {
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

		// The type of the client must be SendType
		clientSendType, clientTypeOk := clientType.(*types.SendType)

		if clientTypeOk {
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

			checkContinuation := p.continuation_e.typecheckForm(gammaNameTypesCtx, providerShadowName, providerType, labelledTypesEnv, sigma)

			return checkContinuation
		} else {
			// wrong type, expected A * B
			return fmt.Errorf("expected '%s' to have a send type (A -o B), but found type '%s' instead", p.from_c.String(), clientType.String())

		}
	}
}

func (p *SelectForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	return nil
}
func (p *BranchForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	return nil
}
func (p *CaseForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	return nil
}
func (p *NewForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	return nil
}

// 1 : close w
func (p *CloseForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	// EndR: 1
	logRule("rule EndR")
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

		// The type of the client must be UnitType
		clientUnitType, clientTypeOk := clientType.(*types.UnitType)

		if clientTypeOk {
			// Set type
			p.to_c.Type = clientUnitType

			// Continue checking the remaining process
			checkContinuation := p.continuation_e.typecheckForm(gammaNameTypesCtx, providerShadowName, providerType, labelledTypesEnv, sigma)

			return checkContinuation
		} else {
			return fmt.Errorf("expected '%s' to have a unit type (1), but found type '%s' instead", p.to_c.String(), clientUnitType.String())
		}
	} else {
		// Waiting on the wrong name
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
	// EndL: 1
	logRule("rule Drp")

	return nil

	// // Can only wait for a client (not self)
	// if !isProvider(p.to_c, providerShadowName) {
	// 	clientType, errorClient := consumeName(p.to_c, gammaNameTypesCtx)
	// 	if errorClient != nil {
	// 		return errorClient
	// 	}

	// 	// The type of the client must be UnitType
	// 	clientUnitType, clientTypeOk := clientType.(*types.UnitType)

	// 	if clientTypeOk {
	// 		// Set type
	// 		p.to_c.Type = clientUnitType

	// 		// Continue checking the remaining process
	// 		checkContinuation := p.continuation_e.typecheckForm(gammaNameTypesCtx, providerShadowName, providerType, labelledTypesEnv, sigma)

	// 		return checkContinuation
	// 	} else {
	// 		return fmt.Errorf("expected '%s' to have a unit type (1), but found type '%s' instead", p.to_c.String(), clientUnitType.String())
	// 	}
	// } else {
	// 	// Waiting on the wrong name
	// 	_, unitTypeOk := providerType.(*types.UnitType)

	// 	if unitTypeOk {
	// 		return fmt.Errorf("expected '%s' to have a wait on a 'non-self' channel instead (%s is acting as self)", p.String(), p.to_c.String())
	// 	} else {
	// 		return fmt.Errorf("expected '%s' to have a unit type (1), but found type '%s' instead", p.String(), providerType.String())
	// 	}
	// }
}

func (p *SplitForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
	return nil
}
func (p *CallForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) error {
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

// func functionExists(functionTypesEnv FunctionTypesEnv, key string) bool {
// 	_, ok := functionTypesEnv[key]

// 	return ok
// }

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

func logRule(s string) {
	fmt.Println(s)
}
