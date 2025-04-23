package app

import "embracer/internal/domain/port"

/*
물론 지금은 단순히 cp.Register를 호출하는 얇은 로직이지만,
실제 서비스에서는 다음처럼 확장 가능해:
service/addr에 대한 유효성 검사
중복 등록 방지
등록 전에 로그 기록 또는 이벤트 발행
내부 서비스 레지스트리에 캐싱
즉, 도메인 안에서 일어나는 업무적인 판단, 흐름이 여기에 있어야 해.
*/
type RegisterUsecase struct {
	cp port.ControlPlaneClient
}

func NewRegisterUsecase(cp port.ControlPlaneClient) *RegisterUsecase {
	return &RegisterUsecase{cp: cp}
}

func (u *RegisterUsecase) Register(cluster, addr string) error {
	return u.cp.Register(cluster, addr)
}
