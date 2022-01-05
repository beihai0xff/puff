package transport

import (
	"net/http"
	"time"

	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp"
	stats "go.etcd.io/etcd/server/v3/etcdserver/api/v2stats"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/beihai0xff/puff/pkg/types"
)

// Transporter 通信模块接口定义
type Transporter interface {
	// Start starts the given Transporter.
	// Start MUST be called before calling other functions in the interface.
	Start() error
	// Handler returns the HTTP handler of the transporter.
	// A transporter HTTP handler handles the HTTP requests
	// from remote peers.
	// The handler MUST be used to handle RaftPrefix(/raft)
	// endpoint.
	Handler() http.Handler
	// Send sends out the given messages to the remote peers.
	// Each message has a To field, which is an id that maps
	// to an existing peer in the transport.
	// If the id cannot be found in the transport, the message
	// will be ignored.
	Send(m []raftpb.Message)
	// SendSnapshot sends out the given snapshot message to a remote peer.
	// The behavior of SendSnapshot is similar to Send.
	SendSnapshot()
	// AddRemote adds a remote with given peer urls into the transport.
	// A remote helps newly joined member to catch up the progress of cluster,
	// and will not be used after that.
	// It is the caller's responsibility to ensure the urls are all valid,
	// or it panics.
	AddRemote(id types.ID, peers []string)
	// AddPeer adds a peer with given peer urls into the transport.
	// It is the caller's responsibility to ensure the urls are all valid,
	// or it panics.
	// Peer urls are used to connect to the remote peer.
	AddPeer(id types.ID, peers []string)
	// RemovePeer removes the peer with given id.
	RemovePeer(id types.ID)
	// RemoveAllPeers removes all the existing peers in the transport.
	RemoveAllPeers()
	// UpdatePeer updates the peer urls of the peer with the given id.
	// It is the caller's responsibility to ensure the urls are all valid,
	// or it panics.
	UpdatePeer(id types.ID, urls []string)
	// ActiveSince returns the time that the connection with the peer
	// of the given id becomes active.
	// If the connection is active since peer was added, it returns the adding time.
	// If the connection is currently inactive, it returns zero time.
	ActiveSince(id types.ID) time.Time
	// ActivePeers returns the number of active peers.
	ActivePeers() int
	// Stop closes the connections and stops the transporter.
	Stop()
}

// GRPCTransport implements Transporter interface. It provides the functionality
// to send raft messages to peers, and receive raft messages from peers.
// User should call Handler method to get a handler to serve requests
// received from peers.
// User needs to call Start before calling other functions, and call
// Stop when the GRPCTransport is no longer used.
type GRPCTransport struct {
	Logger *zap.Logger

	DialTimeout time.Duration // maximum duration before timing out dial of the request
	// DialRetryFrequency defines the frequency of streamReader dial retrial attempts;
	// a distinct rate limiter is created per every peer (default value: 10 events/sec)
	DialRetryFrequency rate.Limit

	ID        types.ID      // local member ID
	Peers     types.Peers   // local peer URLs
	ClusterID types.ID      // raft cluster ID for request validation
	Raft      rafthttp.Raft // raft state machine, to which the GRPCTransport forwards received messages and reports status

	ServerStats *stats.ServerStats // used to record general transportation statistics
	// used to record transportation statistics with followers when
	// performing as leader in raft protocol
	LeaderStats *stats.LeaderStats
	// ErrorC is used to report detected critical errors, e.g.,
	// the member has been permanently removed from the cluster
	// When an error is received from ErrorC, user should stop raft state
	// machine and thus stop the Transport.
	ErrorC chan error
}

func (s *GRPCTransport) Start() error                          { return nil }
func (s *GRPCTransport) Handler() http.Handler                 { return nil }
func (s *GRPCTransport) Send(m []raftpb.Message)               {}
func (s *GRPCTransport) SendSnapshot()                         {}
func (s *GRPCTransport) AddRemote(id types.ID, peers []string) {}
func (s *GRPCTransport) AddPeer(id types.ID, peers []string)   {}
func (s *GRPCTransport) RemovePeer(id types.ID)                {}
func (s *GRPCTransport) RemoveAllPeers()                       {}
func (s *GRPCTransport) UpdatePeer(id types.ID, us []string)   {}
func (s *GRPCTransport) ActiveSince(id types.ID) time.Time     { return time.Time{} }
func (s *GRPCTransport) ActivePeers() int                      { return 0 }
func (s *GRPCTransport) Stop()                                 {}
func (s *GRPCTransport) Pause()                                {}
func (s *GRPCTransport) Resume()                               {}
