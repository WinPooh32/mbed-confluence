package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/anacrolix/confluence/confluence"
	"github.com/anacrolix/dht/v2"
	"github.com/anacrolix/dht/v2/int160"
	peer_store "github.com/anacrolix/dht/v2/peer-store"
	"github.com/anacrolix/missinggo/v2/filecache"
	"github.com/anacrolix/publicip"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
)

const httpShutdownTimeout = 15 * time.Second

type Config struct {
	Addr                string   `json:"addr"`
	PublicIp4           net.IP   `json:"public_ip_4"`
	PublicIp6           net.IP   `json:"public_ip_6"`
	ImplicitTrackers    []string `json:"implicit_trackers"`
	TorrentGraceSeconds int32    `json:"torrent_grace_seconds"`
	CacheCapacityBytes  int64    `json:"cache_capacity_bytes"`
	CacheDir            string   `json:"cache_dir"`
}

func (c Config) TorrentGrace() time.Duration {
	return time.Duration(c.TorrentGraceSeconds) * time.Second
}

type Server struct {
	mut   sync.Mutex
	http  *http.Server
	doneC chan struct{}
	errs  []error
}

func (srv *Server) Start(ctx context.Context, cfg Config) (string, error) {
	laddr, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return "", fmt.Errorf("listening on addr: %w", err)
	}

	handler, err := srv.newConfluenceHandler(ctx, cfg)
	if err != nil {
		return "", fmt.Errorf("init confluencec handler: %w", err)
	}

	srv.http = &http.Server{
		Addr:    cfg.Addr,
		Handler: handler,
	}

	srv.doneC = make(chan struct{})

	go func() {
		defer close(srv.doneC)
		defer laddr.Close()

		if err := srv.http.Serve(laddr); err != nil {
			srv.addErr(fmt.Errorf("serve http: %w", err))
			return
		}
	}()

	return laddr.Addr().String(), nil
}

func (srv *Server) newConfluenceHandler(ctx context.Context, cfg Config) (http.Handler, error) {
	cs, err := newFileCacheClientStorage(cfg.CacheDir, cfg.CacheCapacityBytes)
	if err != nil {
		return nil, err
	}

	cl, err := newTorrentClient(ctx, cfg, cs, torrent.Callbacks{})
	if err != nil {
		return nil, fmt.Errorf("creating torrent client: %w", err)
	}

	ch := confluence.Handler{
		Logger:                   nil,
		TC:                       cl,
		TorrentGrace:             cfg.TorrentGrace(),
		MetainfoCacheDir:         new(string),
		MetainfoStorage:          nil,
		MetainfoStorageInterface: nil,
		OnNewTorrent: func(t *torrent.Torrent, mi *metainfo.MetaInfo) {
			var spec *torrent.TorrentSpec
			if mi != nil {
				spec, _ = torrent.TorrentSpecFromMetaInfoErr(mi)
			} else {
				spec = new(torrent.TorrentSpec)
			}
			spec.Trackers = append(spec.Trackers, cfg.ImplicitTrackers)
			t.MergeSpec(spec)
		},
		DhtServers: nil,
		Storage:    storage.NewClient(cs),
		ModifyUploadMetainfo: func(mi *metainfo.MetaInfo) {
			mi.AnnounceList = append(mi.AnnounceList, cfg.ImplicitTrackers)
			for _, ip := range cl.PublicIPs() {
				mi.Nodes = append(mi.Nodes, metainfo.Node(net.JoinHostPort(
					ip.String(),
					strconv.FormatInt(int64(cl.LocalPort()), 10))))
			}
		},
	}

	ch.OnTorrentGrace = func(t *torrent.Torrent) {
		t.Drop()
	}

	for _, s := range cl.DhtServers() {
		ch.DhtServers = append(ch.DhtServers, s.(torrent.AnacrolixDhtServerWrapper).Server)
	}

	return &ch, nil
}

func (srv *Server) Wait() {
	<-srv.doneC
}

func (srv *Server) Stop(ctx context.Context) error {
	if srv.http == nil {
		return errors.New("http server is not listening")
	}

	ctx, cancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
	defer cancel()

	if err := srv.http.Shutdown(ctx); err != nil {
		return fmt.Errorf("http: shutdown: %w", err)
	}

	return nil
}

func (srv *Server) Err() error {
	errs := srv.errs
	srv.errs = nil
	return errors.Join(errs...)
}

func (srv *Server) addErr(err error) {
	srv.mut.Lock()
	srv.errs = append(srv.errs, err)
	srv.mut.Unlock()
}

func newFileCacheClientStorage(dir string, capacity int64) (storage.ClientImpl, error) {
	cache, err := filecache.NewCache(dir)
	if err != nil {
		return nil, fmt.Errorf("make a new file cache: %w", err)
	}

	cache.SetCapacity(capacity)

	return storage.NewResourcePieces(cache.AsResourceProvider()), nil
}

func newTorrentClient(ctx context.Context, cfg Config, storage storage.ClientImpl, callbacks torrent.Callbacks) (tc *torrent.Client, err error) {
	cconf := torrent.NewDefaultClientConfig()
	cconf.DefaultStorage = storage

	cconf.PublicIp4 = cfg.PublicIp4
	if cconf.PublicIp4 == nil {
		cconf.PublicIp4, err = publicip.Get4(ctx)
		if err != nil {
			log.Printf("error getting public ipv4 address: %v", err)
		}
	}

	cconf.PublicIp6 = cfg.PublicIp6
	if cfg.PublicIp6 == nil {
		cfg.PublicIp6, err = publicip.Get6(ctx)
		if err != nil {
			log.Printf("error getting public ipv6 address: %v", err)
		}
	}

	cconf.ConfigureAnacrolixDhtServer = func(cfg *dht.ServerConfig) {
		cfg.InitNodeId()
		if cfg.PeerStore == nil {
			cfg.PeerStore = &peer_store.InMemory{
				RootId: int160.FromByteArray(cfg.NodeId),
			}
		}
	}

	cconf.Callbacks = callbacks

	cl, err := torrent.NewClient(cconf)
	if err != nil {
		return nil, fmt.Errorf("new torrent client: %w", err)
	}

	return cl, nil
}
