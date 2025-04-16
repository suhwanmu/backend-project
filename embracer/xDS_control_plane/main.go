package main

import (
	"context"
	"control_plane/registry"
	"encoding/json"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"embracer/utils"
	"embracer/utils/log"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointv3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	servicev3 "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	types "github.com/envoyproxy/go-control-plane/pkg/cache/types"

	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	server "github.com/envoyproxy/go-control-plane/pkg/server/v3"

	_ "net/http/pprof"

	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
)

const (
	nodeID      = "envoy-test-node"
	defaultPort = 80
	httpPort    = ":2222"
	grpcPort    = ":18000"
)

type EndpointMap map[string][]string

var (
	snapshotCache cache.SnapshotCache
	endpoints     = EndpointMap{}
	mu            sync.Mutex
)

func main() {

	log.Info().Msgf("\n version: %s\n commit: %s\n build-time: %s\n",
		utils.Version, utils.Commit, utils.BuildTime)

	var cfg utils.Config
	var err error
	err = envconfig.Process("EMBRACER", &cfg)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	log.Info().Msgf("config: %+v", cfg)

	if cfg.Debug {
		_ = log.Initialize("debug")
	} else {
		_ = log.Initialize("info")
	}

	// check connection to kafka broker
	err = utils.CheckKafkaConn(cfg.KafkaBrokers)
	if err != nil {
		log.Fatal().Err(err).Msg("Kafka connection check failed")
	}
	log.Info().Msgf("connection successful to kafka brokers %v", cfg.KafkaBrokers)

	if cfg.Debug {
		runtime.SetBlockProfileRate(1)
		runtime.SetMutexProfileFraction(1)
	}

	kafkaCancel, err := newWork(cfg)
	defer kafkaCancel()
	if err != nil {
		log.Fatal().Err(err).Msg("New Work Run failed")
	}
	snapshotCache = cache.NewSnapshotCache(false, cache.IDHash{}, nil)
	srv := server.NewServer(context.Background(), snapshotCache, nil)

	go runGRPCServer(srv)
	go runHTTPServer()
	select {}
}

func newWork(cfg utils.Config) (context.CancelFunc, error) {

	ingestC, cancel := utils.StartKafka(cfg)

	time.Sleep(time.Duration(time.Second * 2))
	registry.NewRegistry(ingestC).StartRegistry()

	return cancel, nil
}

func runGRPCServer(srv server.Server) {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Error().Msgf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, srv)
	servicev3.RegisterEndpointDiscoveryServiceServer(grpcServer, srv)
	log.Info().Msgf("📡 gRPC server listening on %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Error().Msgf("Failed to serve: %v", err)
	}
}

func runHTTPServer() {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Service string `json:"service"`
			Addr    string `json:"addr"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		if req.Service == "" || req.Addr == "" {
			http.Error(w, "Missing service or addr", http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		// 중복 제거
		for _, existing := range endpoints[req.Service] {
			if existing == req.Addr {
				log.Info().Msgf("ℹ️ [%s] already registered to [%s]", req.Addr, req.Service)
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		endpoints[req.Service] = append(endpoints[req.Service], req.Addr)
		log.Info().Msgf("✅ Registered new endpoint [%s] to [%s]", req.Addr, req.Service)

		updateSnapshot()
		w.WriteHeader(http.StatusOK)
	})

	log.Info().Msgf("🌐 HTTP registration server listening on %s", httpPort)
	if err := http.ListenAndServe(httpPort, nil); err != nil {
		log.Error().Msgf("Failed to run HTTP server: %v", err)
	}
}

func updateSnapshot() {
	resources := []types.Resource{}
	for cluster, addressList := range endpoints {
		cla := generateClusterLoadAssignment(cluster, addressList)
		resources = append(resources, cla)
	}

	version := strconv.FormatInt(time.Now().UnixNano(), 10)
	snap, err := cache.NewSnapshot(version, map[resource.Type][]types.Resource{
		resource.EndpointType: resources,
	})
	if err != nil {
		log.Error().Msgf("❌ failed to create snapshot: %v", err)
		return
	}

	if err := snapshotCache.SetSnapshot(context.Background(), nodeID, snap); err != nil {
		log.Error().Msgf("❌ failed to set snapshot: %v", err)
		return
	}
	log.Info().Msgf("📦 Snapshot updated! Version: %s", version)
}

func generateClusterLoadAssignment(cluster string, endpoints []string) *endpointv3.ClusterLoadAssignment {
	var lbEndpoints []*endpointv3.LbEndpoint

	for _, ep := range endpoints {
		host, portStr, found := strings.Cut(ep, ":")
		port := uint32(defaultPort)
		if found {
			if p, err := strconv.ParseUint(portStr, 10, 32); err == nil {
				port = uint32(p)
			} else {
				log.Error().Msgf("⚠️ Invalid port in endpoint '%s', defaulting to %d", ep, defaultPort)
			}
		} else {
			log.Error().Msgf("⚠️ No port specified in endpoint '%s', defaulting to %d", ep, defaultPort)
		}

		lbEndpoints = append(lbEndpoints, &endpointv3.LbEndpoint{
			HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
				Endpoint: &endpointv3.Endpoint{
					Address: &core.Address{
						Address: &core.Address_SocketAddress{
							SocketAddress: &core.SocketAddress{
								Address: host,
								PortSpecifier: &core.SocketAddress_PortValue{
									PortValue: port,
								},
							},
						},
					},
				},
			},
		})
	}

	return &endpointv3.ClusterLoadAssignment{
		ClusterName: cluster,
		Endpoints: []*endpointv3.LocalityLbEndpoints{{
			LbEndpoints: lbEndpoints,
		}},
	}
}
