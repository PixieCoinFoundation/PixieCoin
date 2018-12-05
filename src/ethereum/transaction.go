package ethereum

type Transaction struct {
	ID    string `json:"id,omitempty"`
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

type TokenTransaction struct {
	ID   string `json:"id,omitempty"`
	From string `json:"from"`
	To   string `json:"to"`
	Data string `json:"data"`
}
