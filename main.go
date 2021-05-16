package main

import (
	"fmt"
	"os"
	"os/exec"
    "syscall"
)

func main() {
    // switch on first argument, commonly "run"
    switch os.Args[1] {
    case "run":
        run()
    // otherwise, unintended behavior
    default:
        panic("oops")
    }
}

func run() {
    // print command info
	fmt.Printf("Running %v as PID %d\n", os.Args[2:], os.Getpid())

    // setup Command
    cmd := exec.Command(os.Args[2], os.Args[3:]...)

    // setup pipelines
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

    // clone new UTS (hostname)
    cmd.SysProcAttr = &syscall.SysProcAttr {
        Cloneflags: syscall.CLONE_NEWUTS,
    }

    // finally, run the command
	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
