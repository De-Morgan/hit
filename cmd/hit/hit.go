package main

import (
	"os"
)

const logo = `
__ __ __ ______
/\ \_\ \ /\ \ /\__ _\
\ \ __ \ \ \ \ \/_/\ \/
\ \_\ \_\ \ \_\ \ \_\
\/_/\/_/ \/_/ \/_/`

func main() {

	e := &env{
		stdout: os.Stdout,
		stderr: os.Stderr,
		args:   os.Args[1:],
	}
	if err := run(e); err != nil {
		os.Exit(1)
	}
}
