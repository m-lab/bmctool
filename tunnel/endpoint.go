package tunnel

import "fmt"

type Endpoint struct {
	Host string
	Port int32
}

func (ep *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", ep.Host, ep.Port)
}
