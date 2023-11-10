package process

import (
	"fmt"
	"phi/types"
)

func Typecheck(processes []*Process, globalEnv *GlobalEnvironment) error {
	errorChan := make(chan error)
	doneChan := make(chan bool)

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
	typecheckFunctionDefinitions(globalEnv)

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

		funcDef.Body.typecheckForm(gammaNameTypesCtx, nil, providerType, labelledTypesEnv, functionDefinitionsEnv)
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
// ... gammaNameTypesCtx NamesTypesCtx   <- names in context to be used (in case of linearity, ...)
// ... providerShadowName *Name          <- name of the process providing on (typically nil, since self is used instead)
// ... providerType types.SessionType    <- the type of the provider (i.e. type of self)
// ... labelledTypesEnv LabelledTypesEnv <- keeps the mapping of pre-defined types (type A = ...)
// ... sigma FunctionTypesEnv            <- keeps the mapping of pre-defined function definitions (let f() : A = ...)

func consumeName(name Name, gammaNameTypesCtx NamesTypesCtx) types.SessionType {
	foundName, ok := gammaNameTypesCtx[name.Ident]

	if ok {
		// If linear then remove
		delete(gammaNameTypesCtx, name.Ident)

		return foundName.Type
	}

	// Problem since the requested name was not found in the gamma
	panic("Problem since the requested name was not found in the gamma. todo set cool error message")
	// return ok
}

// */-o: send w<u, v>
func (p *SendForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) {
	if isProvider(p.to_c, providerShadowName) {
		// MulR
		provider, sendTypeOk := providerType.(*types.SendType)

		if sendTypeOk {
			expectedLeftType := provider.Left
			expectedRightType := provider.Right
			foundLeftType := consumeName(p.payload_c, gammaNameTypesCtx)
			foundRightType := consumeName(p.continuation_c, gammaNameTypesCtx)

			// The expected and found must match
			if !types.EqualType(expectedLeftType, foundLeftType, labelledTypesEnv) {
				panic("errorrrr")
			}

			if !types.EqualType(expectedRightType, foundRightType, labelledTypesEnv) {
				panic("errorrrr")
			}

			// Set the types for the names
			p.payload_c.Type = foundLeftType
			p.continuation_c.Type = foundRightType
		} else {
			// error
		}

		fmt.Println(p.String())
	}

	// ensure that the remaining names in gamma are allow (i.e. memmx names imdendlin)
	if len(gammaNameTypesCtx) > 0 {
		panic("non linear names left in gamma")
	}

	// ImpL

	// at this point gammaNameTypesCtx should not contain linear names
}

func (p *ReceiveForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) {

}

func (p *SelectForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *BranchForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *CaseForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *NewForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *CloseForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *ForwardForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *SplitForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *CallForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *WaitForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *CastForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *ShiftForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *PrintForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerShadowName *Name, providerType types.SessionType, labelledTypesEnv types.LabelledTypesEnv, sigma FunctionTypesEnv) {

}

// /// Fixed Environments
//
// This is constant and set once at the beginning. The information is obtained from the 'type A = ...' definitions.

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

func functionExists(functionTypesEnv FunctionTypesEnv, key string) bool {
	_, ok := functionTypesEnv[key]

	return ok
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

func nameTypeExists(labelledTypesEnv NamesTypesCtx, key string) bool {
	_, ok := labelledTypesEnv[key]
	return ok
}

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
