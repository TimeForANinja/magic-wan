package cluster

type Cluster struct {
	self   *votingPeer
	peers  []*votingPeer
	votes  []*vote
	master *votingPeer
}

func (c *Cluster) CountVote(from string, to string) error {
	vote, err := c.parseVote(from, to)
	if err != nil {
		return err
	}

	c.updateVotes(vote)
	return nil
}
