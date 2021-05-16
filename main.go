package main

import (
	"fmt"
	"os"
	"os/exec"
    "syscall"
    "path/filepath"
    "io/ioutil"
    "strconv"
)

func main() {
    // switch on first argument, commonly "run"
    switch os.Args[1] {
    case "run":
        run()
    case "child":
        child()
    // otherwise, unintended behavior
    default:
        panic("oops")
    }
}

func run() {
    // print command info
	fmt.Printf("Running %v as PID %d\n", os.Args[2:], os.Getpid())

    // setup Command to run self exec
    // this is done in order to create the new namespace and enter it
    // before we exec the real command
    cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

    // setup pipelines
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

    // clone vals
    cmd.SysProcAttr = &syscall.SysProcAttr{
        // create new UTS, PID space, and UID space, as well as new mount space
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
        
        // set up mapping between UID 0 and current OS UID
        Credential: &syscall.Credential {
            Uid: 0,
            Gid: 0,
        },
        UidMappings: []syscall.SysProcIDMap {
            {
                ContainerID: 0,
                HostID: os.Getuid(),
                Size: 1,
            },
        },
        GidMappings: []syscall.SysProcIDMap {
            {
                ContainerID: 0,
                HostID: os.Getgid(),
                Size: 1,
            },
        },
    }

    // finally, run the command
	must(cmd.Run())
}

func child() {
    // print command info
	fmt.Printf("Running %v as PID %d\n", os.Args[2:], os.Getpid())

    cgroup()

    // setup actual command to be runCommand
    cmd := exec.Command(os.Args[2], os.Args[3:]...)

    // setup pipelines
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

    // set hostname within new UTS
    must(syscall.Sethostname([]byte("container")))

    // chroot and chdir into alpine rootfs
    must(syscall.Chroot("/home/skovati/code/git/gocon/rootfs"))
    must(syscall.Chdir("/"))

    //
    must(syscall.Mount("proc", "proc", "proc", 0, ""))

    // finally, run the command
	must(cmd.Run())

    // and unmount proc
    must(syscall.Unmount("proc", 0))
}

func cgroup() {
    cgroups := "/sys/fs/cgroup/"
    pids := filepath.Join(cgroups, "pids")

    must(os.Mkdir(filepath.Join(pids, "alp"), 0755))
    must(ioutil.WriteFile(filepath.Join(pids, "alp/pids.max"), []byte("20"), 0700))

    // removes cgroup in place after the container exits
    must(ioutil.WriteFile(filepath.Join(pids, "alp/notify_on_release"), []byte("1"), 0700))
    must(ioutil.WriteFile(filepath.Join(pids, "alp/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
