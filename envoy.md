# 1. 경로 기반 라우팅 설정 추가 (envoy.yaml)
```yaml
static_resources:
  listeners:
    - name: listener_0
      address:
        socket_address: { address: 0.0.0.0, port_value: 10000 }
      filter_chains:
        - filters:
            - name: envoy.filters.network.http_connection_manager
              typed_config:
                "@type": type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager
                stat_prefix: ingress_http
                codec_type: AUTO
                route_config:
                  name: local_route
                  virtual_hosts:
                    - name: backend
                      domains: ["*"]
                      routes:
                        - match: { prefix: "/users" }
                          route: { cluster: user_service }
                        - match: { prefix: "/orders" }
                          route: { cluster: order_service }
                http_filters:
                  - name: envoy.filters.http.router
  clusters:
    - name: user_service
      connect_timeout: 5s
      type: STRICT_DNS
      lb_policy: ROUND_ROBIN
      load_assignment:
        cluster_name: user_service
        endpoints:
          - lb_endpoints:
              - endpoint:
                  address:
                    socket_address: { address: user-service, port_value: 8080 }
    - name: order_service
      connect_timeout: 5s
      type: STRICT_DNS
      lb_policy: ROUND_ROBIN
      load_assignment:
        cluster_name: order_service
        endpoints:
          - lb_endpoints:
              - endpoint:
                  address:
                    socket_address: { address: order-service, port_value: 8080 }
```

# 2. Docker Compose로 Envoy + Backend 서비스 구성(hot reload)
```yaml
services:
  envoy:
    image: envoyproxy/envoy:latest
    container_name: envoy
    ports:
      - "10000:10000"
      - "19000:19000"  # Admin API 포트 추가 (상태 확인용)
    volumes:
      - ./envoy.yaml:/etc/envoy/envoy.yaml
      - ./entrypoint.sh:/entrypoint.sh
    command: ["/bin/sh", "/entrypoint.sh"]
    networks:
      - envoy-network

  user-service:
    image: my-user-service:latest
    container_name: user-service
    ports:
      - "8080:8080"
    networks:
      - envoy-network

  order-service:
    image: my-order-service:latest
    container_name: order-service
    ports:
      - "8081:8080"
    networks:
      - envoy-network

networks:
  envoy-network:
```

```sh
#!/bin/sh

# Envoy 실행 (백그라운드 모드)
envoy -c /etc/envoy/envoy.yaml --service-cluster my-service &

# Envoy 프로세스 ID 저장
ENVOY_PID=$!

# 설정 변경 감지 후 핫 리로드 실행 (무한 루프)
while inotifywait -e modify /etc/envoy/envoy.yaml; do
    echo "Config changed, reloading Envoy..."
    
    # 설정 파일 검증 (오류가 있으면 reload 중단)
    envoy --mode validate -c /etc/envoy/envoy.yaml
    if [ $? -ne 0 ]; then
        echo "Config validation failed. Skipping reload."
        continue
    fi

    # Envoy 핫 리로드 (SIGHUP 신호 전송)
    kill -s SIGHUP $ENVOY_PID

    # 새 Envoy 프로세스가 준비될 때까지 대기
    while true; do
        STATUS=$(curl -s http://localhost:19000/ready)
        if [ "$STATUS" = "live" ]; then
            echo "New Envoy process is live!"
            break
        fi
        echo "Waiting for Envoy to be ready..."
        sleep 2
    done

    echo "Hot reload completed successfully!"
done
```

