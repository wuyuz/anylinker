package webssh

import (
	"fmt"
	"testing"
)

func TestSsh(t *testing.T) {
	_, err := NewSshClient("root","4477123Wl!","114.215.84.163")
	if err !=nil {
		fmt.Println("err:",err)
	}
}
