package did

import (
	"encoding/json"

	"github.com/tendermint/abci/types"
)

func writeBurnTokenReport(nodeID string, method string, price float64, data string, app *DIDApplication) error {
	key := "SpendGas" + "|" + nodeID
	chkExists := app.state.db.Get(prefixKey([]byte(key)))
	newReport := Report{
		method,
		price,
		data,
	}
	if chkExists != nil {
		var reports []Report
		err := json.Unmarshal([]byte(chkExists), &reports)
		if err != nil {
			return err
		}
		reports = append(reports, newReport)
		value, err := json.Marshal(reports)
		if err != nil {
			return err
		}
		app.SetStateDB([]byte(key), []byte(value))
	} else {
		var reports []Report
		reports = append(reports, newReport)
		value, err := json.Marshal(reports)
		if err != nil {
			return err
		}
		app.SetStateDB([]byte(key), []byte(value))
	}
	return nil
}

func getUsedTokenReport(param string, app *DIDApplication) types.ResponseQuery {
	app.logger.Infof("GetUsedTokenReport, Parameter: %s", param)
	var funcParam GetUsedTokenReportParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.Height, app)
	}
	key := "SpendGas" + "|" + funcParam.NodeID
	value := app.state.db.Get(prefixKey([]byte(key)))
	if value == nil {
		value = []byte("")
		return ReturnQuery(value, "not found", app.state.Height, app)
	}
	return ReturnQuery(value, "success", app.state.Height, app)
}
