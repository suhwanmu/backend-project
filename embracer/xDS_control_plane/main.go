package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointv3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	servicev3 "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	types "github.com/envoyproxy/go-control-plane/pkg/cache/types"

	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	server "github.com/envoyproxy/go-control-plane/pkg/server/v3"

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
	snapshotCache = cache.NewSnapshotCache(false, cache.IDHash{}, nil)
	srv := server.NewServer(context.Background(), snapshotCache, nil)

	go runGRPCServer(srv)
	go runHTTPServer()
	select {}
}

func runGRPCServer(srv server.Server) {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, srv)
	servicev3.RegisterEndpointDiscoveryServiceServer(grpcServer, srv)
	log.Println("📡 gRPC server listening on ", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
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
				log.Printf("ℹ️ [%s] already registered to [%s]", req.Addr, req.Service)
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		endpoints[req.Service] = append(endpoints[req.Service], req.Addr)
		log.Printf("✅ Registered new endpoint [%s] to [%s]", req.Addr, req.Service)

		updateSnapshot()
		w.WriteHeader(http.StatusOK)
	})

	log.Println("🌐 HTTP registration server listening on", httpPort)
	if err := http.ListenAndServe(httpPort, nil); err != nil {
		log.Fatalf("Failed to run HTTP server: %v", err)
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
		log.Printf("❌ failed to create snapshot: %v", err)
		return
	}

	if err := snapshotCache.SetSnapshot(context.Background(), nodeID, snap); err != nil {
		log.Printf("❌ failed to set snapshot: %v", err)
		return
	}
	log.Printf("📦 Snapshot updated! Version: %s", version)
}

func loadEndpointsFromJSON(path string) EndpointMap {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Printf("⚠️ Failed to read file, starting with empty endpoints: %v", err)
		return EndpointMap{}
	}
	var m EndpointMap
	if err := json.Unmarshal(file, &m); err != nil {
		log.Printf("⚠️ Failed to parse JSON: %v", err)
		return EndpointMap{}
	}
	return m
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
				log.Printf("⚠️ Invalid port in endpoint '%s', defaulting to %d", ep, defaultPort)
			}
		} else {
			log.Printf("⚠️ No port specified in endpoint '%s', defaulting to %d", ep, defaultPort)
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
