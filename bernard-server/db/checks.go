package db

import "time"

//TODO use client version
type Check struct {
	ID       string   `db:"id"`
	HostID   int64    `db:"host_id"`
	Name     string   `db:"name"`
	Command  string   `db:"command"`
	Args     []string `db:"args"`
	Env      []string `db:"env"`
	Dir      string   `db:"dir"`
	Interval int64    `db:"interval"`
	Timeout  int64    `db:"timeout"`
}

type CheckResult struct {
	CheckID string    `db:"check_id"`
	Time    time.Time `db:"time"`
	Status  int64     `db:"status"`
	Output  []byte    `db:"output"`
}
