//go:build integration && scylla
// +build integration,scylla

package gocql

import (
	"context"
	"fmt"
	"net"
	"testing"
)

// unixSocketDialer is a special dialer which connects only to the maintenance_socket.
type unixSocketDialer struct {
	dialer     net.Dialer
	socketPath string
}

func (d unixSocketDialer) DialContext(_ context.Context, _, _ string) (net.Conn, error) {
	return d.dialer.Dial("unix", d.socketPath)
}

func TestUnixSockets(t *testing.T) {
	socketPath := "/tmp/scylla/cql.m"

	c := createCluster()
	c.NumConns = 1
	c.DisableInitialHostLookup = true
	c.ProtoVersion = 3
	c.ReconnectInterval = 0
	c.WriteCoalesceWaitTime = 0

	c.Events.DisableNodeStatusEvents = true
	c.Events.DisableTopologyEvents = true
	c.Events.DisableSchemaEvents = true

	d := net.Dialer{
		Timeout: c.Timeout,
	}
	if c.SocketKeepalive > 0 {
		d.KeepAlive = c.SocketKeepalive
	}

	c.Dialer = unixSocketDialer{
		dialer:     d,
		socketPath: socketPath,
	}

	sess, err := c.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("unable to create session: %v", err))
	}

	defer sess.Close()

	keyspace := "test1"

	err = createTable(sess, `DROP KEYSPACE IF EXISTS `+keyspace)
	if err != nil {
		t.Fatal("unable to drop keyspace if exists:", err)
	}

	err = createTable(sess, fmt.Sprintf(`CREATE KEYSPACE %s
	WITH replication = {
		'class' : 'NetworkTopologyStrategy',
		'replication_factor' : 1
	}`, keyspace))
	if err != nil {
		t.Fatal("unable to create keyspace:", err)
	}
}
