package opinion

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type Opinion struct {
	Map map[string]float64
}

func NewOpinion() (o *Opinion) {
	o = &Opinion{}
	o.Map = make(map[string]float64)
	return
}

func (o Opinion) GetOpinionForHost(host string) (hostOpinion float64, err error) {
	opinion, present := o.Map[host]
	if !present {
		return 0.0, errors.New("Could not find opinion for host " + host)
	}
	return opinion, nil
}

func (o *Opinion) SetOpinionForHost(host string, hostOpinion float64) (err error) {
	opinion, present := o.Map[host]
	if present {
		log.Print("Overwriting old opinion ", opinion, " for ", host, " with new one ", hostOpinion)
	} else {
		log.Print("Creating opinion for ", host, " value ", hostOpinion)
	}
	o.Map[host] = hostOpinion
	return nil
}

func (o *Opinion) Print() (data string) {
	for k, v := range o.Map {
		data += k + ":" + fmt.Sprintf("%f", v) + " "
	}
	return
}

func (o *Opinion) ToJson() []byte {
	ret, err := json.Marshal(*o)
	if err != nil {
		panic(err)
	}
	return ret
}
