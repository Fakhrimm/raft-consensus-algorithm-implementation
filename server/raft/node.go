package raft

type State int

const (
	Leader State = iota // https://go.dev/wiki/Iota
	Candidate
	Follower
)

type Address struct {
	IpAddress string
	Port      int
}

type Node struct {
	state                State
	electionTerm         int
	app                  Application
	address              Address
	clusterLeaderAddress Address
	clusterAddresses     []Address
	log                  []string
}
