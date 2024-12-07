package shared

type Cluster interface {
	OnVoteReceived(message *VoteMessage) error
}
