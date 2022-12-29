# signupin-api
## About
```
표준에 맞춰 클린 아키텍처 개념 활용하여 개발 (https://github.com/golang-standards/project-layout)

http error code, API format, regex, jwt, mongoDB wrapper 등을 별도 구현 (이하: 라이브러리)

해당 라이브러리를 참조해서 본 프로젝트 구현

📌 토큰 기반 인증 (만료 시간 ⏰ 1분)
📌 실행에 필요한 환경 변수 위치 signupin-api/config/.env
```

## Run (Local)
```
signupin-api % cd cmd
cmd % go run main.go
```

## Framework, database used
```
golang (1.19) / gin
MongoDB
```

## Libraries used
- `gin-gonic`: https://github.com/gin-gonic/gin
- `jwt-go`: https://github.com/dgrijalva/jwt-go
- `mgm`: https://github.com/Kamva/mgm
- `go-common`: https://github.com/kkodecaffeine/go-common

## APIs
```
전화번호 인증 API.  → POST. , /api/v1/auth/sms
회원 가입 API.     → POST. , /api/v1/auth/sign-up
회원 접속 API.     → POST. , /api/v1/auth/sign-in
비밀번호 수정 API.  → PUT.  , /api/v1/users/reset-password
회원 정보 조회 API. → GET.  , /api/v1/users/:userID

📌 전화번호 인증 시 임의로 생성한 6자리 문자열을 인증번호로 간주 (ex. 683577)
```
