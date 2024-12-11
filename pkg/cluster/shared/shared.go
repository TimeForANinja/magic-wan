package shared

type Cluster interface {
	OnVoteReceived(message *VoteMessage) error
	OnNewConfig(listener func(config any)) error
}
