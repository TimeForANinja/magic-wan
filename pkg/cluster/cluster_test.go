package cluster

import (
	"magic-wan/pkg/cluster/shared"
	"magic-wan/pkg/various"
	"testing"
)

var testPeers = []string{
	"1.0.0.0",
	"2.0.0.0",
	"3.0.0.0",
	"4.0.0.0",
}
var testStalePeerIDX = 2

func buildState_newCluster() *Cluster[CoreConfig] {
	c := NewCluster(testPeers[0])
	for _, peer := range testPeers[1:] {
		c.AddPeer(peer)
	}
	return c
}

func fakeVote(c *Cluster[CoreConfig], peer, selectedPeer string) {
	err := c.OnVoteReceived(&shared.VoteMessage{
		Voter: peer,
		Vote:  selectedPeer,
	})
	if err != nil {
		panic(err)
	}
}

func buildState_masterStale() *Cluster[CoreConfig] {

	c := NewCluster(testPeers[0])
	for _, peer := range testPeers[1:] {
		c.AddPeer(peer)
	}

	// set selected peer as master (n-1 votes)
	for _, peer := range testPeers[:len(testPeers)-1] {
		fakeVote(c, peer, testPeers[testStalePeerIDX])
	}
	// last peer votes different, to make the test more complex
	fakeVote(c, testPeers[len(testPeers)-1], testPeers[testStalePeerIDX+1])

	// set failed peer to stale
	peer := various.ArrayFind(c.peers, func(v *votingPeer) bool {
		return v.ip == testPeers[testStalePeerIDX]
	})
	peer.failedSends = maxFailedSends + 2

	return c
}

func Test_getValidVotes(t *testing.T) {
	cluster := buildState_newCluster()
	votes := cluster.getValidVotes()
	if len(votes) != 0 {
		t.Fatalf("Expected 0 votes, got %d", len(votes))
	}

	cluster = buildState_masterStale()
	votes = cluster.getValidVotes()
	// one vote, 3 Votes to the stale master are invalid, one to self is valid
	if len(votes) != 1 {
		t.Fatalf("Expected %d vote, got %d", len(testPeers)-1, len(votes))
	}
	stalePeer := various.ArrayFind(votes, func(p *peerVote) bool {
		return p.voter.isStale()
	})
	if stalePeer != nil {
		t.Fatalf("Expected stale peer to be nil, got %s", stalePeer.voter.ip)
	}
}

func Test_getMasterToVote(t *testing.T) {
	cluster := buildState_masterStale()
	newMaster := cluster.getMasterToVote()
	if newMaster.isStale() {
		t.Fatalf("Expected new master to be valid, got stale note \"%s\"", newMaster.ip)
	}
	if newMaster.ip == testPeers[testStalePeerIDX] {
		t.Fatal("Expected master not to be the test stale peer")
	}
}

func Test_masterCandidate(t *testing.T) {
	cluster := buildState_newCluster()
	candidates := cluster.masterCandidate()
	if len(candidates) != len(testPeers) {
		t.Fatalf("Expected %d candidates, got %d", len(testPeers), len(candidates))
	}

	cluster = buildState_masterStale()
	candidates = cluster.masterCandidate()
	// only one valid vote, as can be seen in the getValidVotes() test
	if len(candidates) != 1 {
		t.Fatalf("Expected 1 candidate, got %d", len(candidates))
	}
	for i, candidate := range candidates {
		if candidate.isStale() {
			t.Fatalf("Expected candidate %d to be valid, got stale note \"%s\"", i, candidate.ip)
		}
	}
}
