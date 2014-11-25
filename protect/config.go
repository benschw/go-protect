package protect

type Config struct {
	RaftHost  string
	RaftPort  int
	ApiHost   string
	ApiPort   int
	DataDir   string
	JoinAddr  string
	Bootstrap bool
}
