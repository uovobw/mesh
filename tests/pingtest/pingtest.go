package pingtest

import (
	"fmt"
	"os/exec"
	"strconv"
)

type Pingtest struct {
	count   int
	address string
}

func NewPingTest(address string, count int) *Pingtest {
	ret := Pingtest{}
	ret.address = address
	ret.count = count
	return &ret
}

func (p *Pingtest) Run() (bool, error) {
	cmd := exec.Command("/bin/ping", "-c", strconv.Itoa(p.count), p.address)
	fmt.Println("path:", cmd.Path)
	fmt.Println("args:", cmd.Args)
	fmt.Println("dir:", cmd.Dir)

	err := cmd.Run()
	if err == nil {
		return true, nil
	}
	return false, err
}
