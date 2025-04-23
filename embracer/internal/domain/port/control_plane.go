package port

type ControlPlaneClient interface {
    Register(cluster, addr string) error
}
