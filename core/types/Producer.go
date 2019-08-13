package types

type Vote_info struct {
	Producer_public_key string `json:",omitempty"`
	Vote_type           string `json:",omitempty"`
	Txid                string `json:",omitempty"`
	N                   int    `json:",omitempty"`
	Value               string `json:",omitempty"`
	Outputlock          int    `json:",omitempty"`
	Address             string `json:",omitempty"`
	Block_time          int64  `json:",omitempty"`
	Height              int64  `json:",omitempty"`
	Rank                int64  `json:",omitempty"`
	Producer_info       `json:",omitempty"`
	Is_valid            string `json:",omitempty"`
	Reward              string `json:",omitempty"`
	EstRewardPerYear    string `json:",omitempty"`
}

type Vote_statistic_header struct {
	Value      string   `json:",omitempty"`
	Node_num   int      `json:",omitempty"`
	Txid       string   `json:",omitempty"`
	Height     int64    `json:",omitempty"`
	Nodes      []string `json:",omitempty"`
	Block_time int64    `json:",omitempty"`
	Is_valid   string   `json:",omitempty"`
}

type Vote_statistic struct {
	Vote_Header Vote_statistic_header `json:",omitempty"`
	Vote_Body   []Vote_info           `json:",omitempty"`
}

type Vote_statisticSorter []Vote_statistic

func (a Vote_statisticSorter) Len() int      { return len(a) }
func (a Vote_statisticSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Vote_statisticSorter) Less(i, j int) bool {
	return a[i].Vote_Header.Height > a[j].Vote_Header.Height
}

type Producer_info struct {
	Ownerpublickey string
	Nodepublickey  string
	Nickname       string
	Url            string
	Location       int64
	Active         int
	Votes          string
	Netaddress     string
	State          string
	Registerheight int64
	Cancelheight   int64
	Inactiveheight int64
	Illegalheight  int64
	Index          int64
}
