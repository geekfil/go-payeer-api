package payeer_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type CheckPayoutToExternalResponse struct {
	Error
	OutputParams struct {
		SumIn  int    `json:"sumIn"`
		CurIn  string `json:"curIn"`
		CurOut string `json:"curOut"`
		Ps     int    `json:"ps"`
		SumOut int    `json:"sumOut"`
	} `json:"outputParams"`
	HistoryID int `json:"historyId"`
}

func (p *Payeer) CheckPayoutToExternal(ps, sumIn, curIn, curOut string, fields map[string]string) (*CheckPayoutToExternalResponse, error) {
	data := &bytes.Buffer{}
	p.data.Add("action", "initOutput")
	p.data.Add("ps", ps)
	p.data.Add("sumIn", sumIn)
	p.data.Add("curIn", curIn)
	p.data.Add("curOut", curOut)
	for paramKey, paramValue := range fields {
		p.data.Add(fmt.Sprintf("param_%s", paramKey), paramValue)
	}

	_, err := data.WriteString(p.data.Encode())
	if err != nil {
		return nil, err
	}
	res, err := p.request(data)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resData := &CheckPayoutToExternalResponse{}
	if err := json.NewDecoder(res.Body).Decode(resData); err != nil {
		if strings.Contains(err.Error(),"errors") {
			resData.Errors = []string{}
		}else{
			return nil, err
		}
	}
	if len(resData.Errors) != 0 {
		return nil, errors.New(resData.Error.Error())
	} else {
		return resData, err
	}
}
