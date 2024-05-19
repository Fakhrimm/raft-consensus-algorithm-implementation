package Node

type ServerState struct {
	// Persistent
	currentTerm int
	votedFor    int
	log         []Entry

	// Volatile
	commitIndex int
	lastApplied int

	// Leaders
	nextIndex  []int
	matchIndex []int
}
