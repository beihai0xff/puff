package server

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	stats "go.etcd.io/etcd/server/v3/etcdserver/api/v2stats"
	"go.uber.org/zap"

	"github.com/beihai0xff/puff/pkg/types"
	"github.com/beihai0xff/puff/server/transport"
)

type raftNodeConfig struct {
	id      uint64 // client ID for raft session
	peerMap map[uint64]string

	logger *zap.Logger

	// to check if msg receiver is removed from cluster
	isIDRemoved func(id uint64) bool
	raft.Node
	raftStorage *raft.MemoryStorage
	// heartbeat 心跳消息发送间隔
	heartbeat time.Duration // for logging

}

type raftNode struct {
	tickMu *sync.Mutex

	node raft.Node
	// MemoryStorage 是 etcd raft 提供的一个基于内存的实现，并不能进行持久化
	raftStorage *raft.MemoryStorage

	// transport specifies the transport to send and receive msgs to members.
	// Sending messages MUST NOT block. It is okay to drop messages, since
	// clients should timeout and reissue their messages.
	// If transport is nil, server will panic.
	transport transport.Transporter

	// 提供一个周期性的时钟定时触发 Tick 方法
	ticker *time.Ticker

	raftNodeConfig

	done <-chan struct{}
}

func newRaftNode(config raftNodeConfig) *raftNode {
	n := &raftNode{
		raftNodeConfig: config,
		raftStorage:    raft.NewMemoryStorage(),
		// TODO: 添加通信模块实现
		transport: &transport.GRPCTransport{},
		done:      make(chan struct{}),
	}
	if n.heartbeat == 0 {
		n.ticker = &time.Ticker{}
	} else {
		n.ticker = time.NewTicker(n.heartbeat)
	}

	go n.startNode()
	return n
}

func (rn *raftNode) startNode() {
	peers := []raft.Peer{}
	for i := range rn.peerMap {
		peers = append(peers, raft.Peer{ID: i})
	}
	c := &raft.Config{
		ID:              rn.id,
		ElectionTick:    10,
		HeartbeatTick:   1,
		Storage:         rn.raftStorage,
		MaxSizePerMsg:   4096,
		MaxInflightMsgs: 256,
	}
	rn.node = raft.StartNode(c, peers)
	rn.transport = &transport.GRPCTransport{
		Logger:      rn.logger,
		ID:          types.ID(rn.id),
		ClusterID:   0x1000,
		Raft:        rn,
		ServerStats: stats.NewServerStats("", ""),
		LeaderStats: stats.NewLeaderStats(zap.NewExample(), strconv.FormatUint(rn.id, 10)),
		ErrorC:      make(chan error),
	}
	rn.transport.Start()
	for peer, addr := range rn.peerMap {
		if peer != rn.id {
			rn.transport.AddPeer(types.ID(peer), []string{addr})
		}
	}
	go rn.serveRaft()
	go rn.run()
}

func (rn *raftNode) serveRaft() {
	addr := rn.peerMap[rn.id][strings.LastIndex(rn.peerMap[rn.id], ":"):]
	server := http.Server{
		Addr:    addr,
		Handler: rn.transport.Handler(),
	}
	server.ListenAndServe()
}

func (rn *raftNode) run() {

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rn.node.Tick()
		case rd := <-rn.node.Ready():
			rn.raftStorage.Append(rd.Entries)
			rn.transport.Send(rd.Messages)
			if !raft.IsEmptySnap(rd.Snapshot) {
				rn.raftStorage.ApplySnapshot(rd.Snapshot)
			}
			for _, entry := range rd.CommittedEntries {
				// TODO: handle Entries
				switch entry.Type {
				case raftpb.EntryNormal:
					log.Printf("Receive committed data on node %v: %v\n", rn.id, string(entry.Data))
				case raftpb.EntryConfChangeV2:
					var cc raftpb.ConfChangeV2
					cc.Unmarshal(entry.Data)
					rn.node.ApplyConfChange(cc)
				}
			}
			rn.node.Advance()
		case <-rn.done:
			// stop raft state machine and thus stop the Transport.
			rn.transport.Stop()
			return
		}
	}

}

// raft.Node does not have locks in Raft package
func (rn *raftNode) tick() {
	rn.tickMu.Lock()
	defer rn.tickMu.Unlock()
	rn.Tick()

}

func (rn *raftNode) Process(ctx context.Context, m raftpb.Message) error {
	return rn.node.Step(ctx, m)
}
func (rn *raftNode) IsIDRemoved(id uint64) bool                           { return false }
func (rn *raftNode) ReportUnreachable(id uint64)                          {}
func (rn *raftNode) ReportSnapshot(id uint64, status raft.SnapshotStatus) {}
