package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
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
	jsonFile    = "endpoints.json"
	defaultPort = 80
)

type EndpointMap map[string][]string

func main() {
	snapshotCache := cache.NewSnapshotCache(false, cache.IDHash{}, nil)
	srv := server.NewServer(context.Background(), snapshotCache, nil)

	go runGRPCServer(srv)

	endpoints := loadEndpointsFromJSON(jsonFile)

	resources := []types.Resource{}
	for cluster, addressList := range endpoints {
		cla := generateClusterLoadAssignment(cluster, addressList)
		resources = append(resources, cla)
		log.Printf("✔ Loaded endpoints for [%s]: %v\n", cluster, addressList)
	}

	snapshotVersion := strconv.FormatInt(time.Now().UnixNano(), 10)
	snap, err := cache.NewSnapshot(snapshotVersion, map[resource.Type][]types.Resource{
		resource.EndpointType: resources,
	})
	if err != nil {
		log.Fatalf("failed to create snapshot: %v", err)
	}

	if err := snapshotCache.SetSnapshot(context.Background(), nodeID, snap); err != nil {
		log.Fatalf("failed to set snapshot: %v", err)
	}

	log.Println("🚀 xDS EDS server is running and snapshot loaded.")
	select {}
}

func runGRPCServer(srv server.Server) {
	lis, err := net.Listen("tcp", ":18000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, srv)
	servicev3.RegisterEndpointDiscoveryServiceServer(grpcServer, srv)
	log.Println("📡 gRPC server listening on :18000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func loadEndpointsFromJSON(path string) EndpointMap {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	var m EndpointMap
	if err := json.Unmarshal(file, &m); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
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
