package pingtest

import (
	"fmt"
	"os/exec"
	"strconv"
)

type Pingtest struct {
	Count   int
	Address string
}

func (p *Pingtest) Setup(data map[string]string) {
	address, ok := data["address"]
	if !ok {
		panic("Address not found in Pingtest init")
	}
	count, ok := data["count"]
	if !ok {
		panic("Count not found in Pingtest init")
	}
	p.Address = address
	p.Count, _ = strconv.Atoi(count)
}

func (p *Pingtest) Run() (float64, error) {
	cmd := exec.Command("/bin/ping", "-c", strconv.Itoa((*p).Count), (*p).Address)
	fmt.Println("path:", cmd.Path)
	fmt.Println("args:", cmd.Args)
	fmt.Println("dir:", cmd.Dir)

	err := cmd.Run()
	if err != nil {
		return 0.0, err
	}
	return 1.0, nil
}
