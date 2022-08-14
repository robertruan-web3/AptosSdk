package aptos

type CommonRsp struct {
	Success bool `json:"success"`
}

type MintRsp struct {
	Success bool   `json:"success"`
	Type    string `json:"type"`
}

type AccountBalanceRsp struct {
	Type string `json:"type"`
	Data struct {
		Coin struct {
			Value string `json:"value"`
		} `json:"coin"`
	} `json:"data"`
}
