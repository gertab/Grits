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

	fmt.Println("ok")
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
	labelledTypesEnv := produceLabelledSessionTypeEnvironment(*globalEnv.Types)
	functionDefinitionsEnv := produceFunctionDefinitionsEnvironment(*globalEnv.FunctionDefinitions)

	for _, funcDef := range *globalEnv.FunctionDefinitions {
		gammaNameTypesCtx := produceNameTypesCtx(funcDef.Parameters)
		providerType := funcDef.Type

		funcDef.Body.typecheckForm(gammaNameTypesCtx, providerType, labelledTypesEnv, functionDefinitionsEnv)
	}

	return nil
}

// ------- Syntax directed typechecking
// Each form has a dedicated typechecking function
type LabelledTypesEnv map[string]LabelledType
type FunctionTypesEnv map[string]FunctionType /* represented as sigma in the type system */
type NamesTypesCtx map[string]NamesType       /* maps names to their types */

// ⊗/: send w<u, v>
func (p *SendForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerType types.SessionType, a LabelledTypesEnv, sigma FunctionTypesEnv) {
	// if p.to_c.IsSelf || {

	// }
	// MulL

	// ImpL
}

func (p *ReceiveForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerType types.SessionType, a LabelledTypesEnv, sigma FunctionTypesEnv) {

}

func (p *SelectForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerType types.SessionType, a LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *BranchForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerType types.SessionType, a LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *CaseForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerType types.SessionType, a LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *NewForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerType types.SessionType, a LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *CloseForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerType types.SessionType, a LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *ForwardForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerType types.SessionType, a LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *SplitForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerType types.SessionType, a LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *CallForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerType types.SessionType, a LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *WaitForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerType types.SessionType, a LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *CastForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerType types.SessionType, a LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *ShiftForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerType types.SessionType, a LabelledTypesEnv, sigma FunctionTypesEnv) {

}
func (p *PrintForm) typecheckForm(gammaNameTypesCtx NamesTypesCtx, providerType types.SessionType, a LabelledTypesEnv, sigma FunctionTypesEnv) {

}

// /// Fixed Environments
//
// labelledTypesEnv: map of labels to their session type (wrapped in a LabelledType struct)
// This is constant and set once at the beginning. The information is obtained from the 'type A = ...' definitions.
type LabelledType struct {
	Name string
	Type types.SessionType
}

func produceLabelledSessionTypeEnvironment(types []types.SessionTypeDefinition) LabelledTypesEnv {
	labelledTypesEnv := make(LabelledTypesEnv)
	for _, j := range types {
		labelledTypesEnv[j.Name] = LabelledType{Type: j.SessionType, Name: j.Name}
	}

	return labelledTypesEnv
}

func labelledSessionTypeExists(labelledTypesEnv LabelledTypesEnv, key string) bool {
	_, ok := labelledTypesEnv[key]

	return ok
}

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
