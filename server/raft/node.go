package raft

type State int

const (
	Leader State = iota // https://go.dev/wiki/Iota
	Candidate
	Follower
)

type Node struct {
	state State
	// TODO: Complete this
}
