package main

import (
	"testing"

	"github.com/google/syzkaller/pkg/rpcserver"
	"github.com/google/syzkaller/pkg/rpcserver/mocks"
	"github.com/stretchr/testify/assert"
)

func TestInitRPC(t *testing.T) {
	mgr := &Manager{}

	emptyServer := rpcserver.Server{}

	serv := mocks.NewServerInterface[rpcserver.Server](t)
	serv.On("Start").Return(nil)
	serv.On("Ptr").Return(&emptyServer)

	mgr.startRPC(serv)
	assert.Equal(t, mgr.serv, serv)
}
