package raft

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.etcd.io/etcd/client/pkg/v3/types"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp"
	stats "go.etcd.io/etcd/server/v3/etcdserver/api/v2stats"
	"go.uber.org/zap"
)

type raftNode struct {
	id      uint64
	peerMap map[uint64]string

	node        raft.Node
	raftStorage *raft.MemoryStorage

	transport *rafthttp.Transport

	logger *zap.Logger
}

func newRaftNode(id uint64, peerMap map[uint64]string) *raftNode {
	n := &raftNode{
		id:          id,
		peerMap:     peerMap,
		raftStorage: raft.NewMemoryStorage(),
		logger:      zap.NewExample(),
	}
	go n.startRaft()
	return n
}

func (rn *raftNode) startRaft() {
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
	rn.transport = &rafthttp.Transport{
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
	go rn.serveChannels()
}

func (rn *raftNode) serveRaft() {
	addr := rn.peerMap[rn.id][strings.LastIndex(rn.peerMap[rn.id], ":"):]
	server := http.Server{
		Addr:    addr,
		Handler: rn.transport.Handler(),
	}
	server.ListenAndServe()
}

func (rn *raftNode) serveChannels() {

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
				case raftpb.EntryConfChange:
					var cc raftpb.ConfChange
					cc.Unmarshal(entry.Data)
					rn.node.ApplyConfChange(cc)
				}
			}
			rn.node.Advance()
		case err := <-rn.transport.ErrorC:
			// stop raft state machine and thus stop the Transport.
			log.Fatal(err)
		}
	}

}

func (rn *raftNode) Process(ctx context.Context, m raftpb.Message) error {
	return rn.node.Step(ctx, m)
}
func (rn *raftNode) IsIDRemoved(id uint64) bool                           { return false }
func (rn *raftNode) ReportUnreachable(id uint64)                          {}
func (rn *raftNode) ReportSnapshot(id uint64, status raft.SnapshotStatus) {}
