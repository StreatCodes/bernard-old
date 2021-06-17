package db

//Next Sequence
type Host struct {
	ID          int64  `db:"id"`
	HostName    string `db:"host_name"`
	Description string `db:"description"`
	Key         []byte `db:"key"`
	Checks      []Check
}
