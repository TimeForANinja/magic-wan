package shared

type VoteMessage struct {
	Voter         string `json:"voter"`
	Vote          string `json:"vote"`
	ConfigVersion uint16 `json:"configVersion"`
}
