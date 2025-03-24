# backend-project
Projects for backend developers

## 한국 유튜브 동시접속자 수
![image](https://github.com/user-attachments/assets/9ccbe627-df8b-4216-958f-c86794c4f16d)

## 1. 핵심 고려사항
* **고가용성(High Availability)**: 50만 동시 접속을 처리할 수 있어야 한다.
* **확장성(Scalability)**: 트래픽 증가에 따라 수평 확장이 가능해야 한다.
* **데이터 일관성(Consistency)**: 프로필 데이터가 중요하므로 적절한 일관성 모델을 유지해야 한다.
* **성능(Performance)**: 빠른 조회와 트랜잭션 처리가 가능해야 한다.

## 2. 어플리에이션 아키텍쳐

**Hexagonal**
![golang clean architecture](https://github.com/bxcodec/go-clean-arch/raw/master/clean-arch.png)
  
