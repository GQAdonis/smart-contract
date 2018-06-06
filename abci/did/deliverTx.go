package did

import (
	"fmt"
	"reflect"

	"github.com/ndidplatform/smart-contract/abci/code"
	"github.com/tendermint/abci/types"
)

// ReturnDeliverTxLog return types.ResponseDeliverTx
func ReturnDeliverTxLog(code uint32, log string, extraData string) types.ResponseDeliverTx {
	return types.ResponseDeliverTx{
		Code: code,
		Log:  fmt.Sprintf(log),
		Data: []byte(extraData),
	}
}

// DeliverTxRouter is Pointer to function
func DeliverTxRouter(method string, param string, nonce string, signature string, nodeID string, app *DIDApplication) types.ResponseDeliverTx {
	funcs := map[string]interface{}{
		"InitNDID":                   initNDID,
		"RegisterNode":               registerNode,
		"RegisterMsqDestination":     registerMsqDestination,
		"AddAccessorMethod":          addAccessorMethod,
		"CreateRequest":              createRequest,
		"CreateIdpResponse":          createIdpResponse,
		"SignData":                   signData,
		"RegisterServiceDestination": registerServiceDestination,
		"RegisterMsqAddress":         registerMsqAddress,
		"AddNodeToken":               addNodeToken,
		"ReduceNodeToken":            reduceNodeToken,
		"SetNodeToken":               setNodeToken,
		"SetPriceFunc":               setPriceFunc,
		"CloseRequest":               closeRequest,
		"TimeOutRequest":             timeOutRequest,
		"AddNamespace":               addNamespace,
		"DeleteNamespace":            deleteNamespace,
		"UpdateNode":                 updateNode,
		"CreateIdentity":             createIdentity,
		"UpdateValidator":            updateValidator,
		"RegisterService":            registerService,
		"DeleteService":              deleteService,
	}

	// ---- check authorization ----
	checkTxResult := CheckTxRouter(method, param, nonce, signature, nodeID, app)
	if checkTxResult.Code != code.OK {
		// return result = false
		var result types.ResponseDeliverTx
		result.Code = checkTxResult.Code
		result.Log = checkTxResult.Log
		return result
	}

	value, _ := callDeliverTx(funcs, method, param, app, nodeID)
	result := value[0].Interface().(types.ResponseDeliverTx)
	// ---- Burn token ----
	if result.Code == code.OK {
		if !isNDIDMethod[method] {
			needToken := getTokenPriceByFunc(method, app)
			err := reduceToken(nodeID, needToken, app)
			if err != nil {
				result.Code = code.TokenAccountNotFound
				result.Log = err.Error()
				return result
			}
			// Write burn token report
			// only have result.Data in some method
			writeBurnTokenReport(nodeID, method, needToken, string(result.Data), app)
		}
	}
	return result
}

func callDeliverTx(m map[string]interface{}, name string, param string, app *DIDApplication, nodeID string) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	in := make([]reflect.Value, 3)
	in[0] = reflect.ValueOf(param)
	in[1] = reflect.ValueOf(app)
	in[2] = reflect.ValueOf(nodeID)
	result = f.Call(in)
	return
}
