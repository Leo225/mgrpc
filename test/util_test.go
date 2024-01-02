package test

import (
	"fmt"
	"testing"

	"github.com/Leo225/mgrpc"
)

func TestLocalIPv4(t *testing.T) {
	r, err := mgrpc.LocalIPv4("")
	if err != nil {
		fmt.Println("TestLocalIPv4err: ", err)
		return
	}
	fmt.Println("TestLocalIPv4 result", r)
}
