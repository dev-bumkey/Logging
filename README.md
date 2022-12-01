# Acloud Alarm Collector
Acloud alarm agent가 전송한 Prometheus Alert을 DB에 저장한다.

* 수신 방식
  - HTTP 직접 전송: dashboard proxy를 통하며, json 문서를 gzip으로 압축한 형태
  - Kafka consumer

</br>

---
## 환경 변수
### version 4.8.0
|이름|설명|값|필수|기본값|참고|버전|
|------|---|---|---|---|---|--|
|SERVICE_HOST|서비스를 실행하는 host ip 주소 지정|string|O|0.0.0.0||>= 4.8.0|
|SERVICE_PORT|서비스를 실행하는 host port 번호|int|O|9308||>= 4.8.0|
|DB_ENGINE|알람을 저장할 디비 엔진|string|X|postgress||>= 4.8.0|
|DB_USER|알람을 저장할 디비 사용자 아이디|string|O||보통 secret을 mount해서 사용한다.|>= 4.8.0|
|DB_PASSWORD|알람을 저장할 디비 사용자 암호|string|O||보통 secret을 mount해서 사용한다.|>= 4.8.0|
|DB_NAME|알람을 저장할 디비 이름|string|O|||>= 4.8.0|
|DB_HOST|알람을 저장할 디비 ip 주소|string|O|||>= 4.8.0|
|DB_PORT|알람을 저장할 디비 포트|int|O|||>= 4.8.0|
|DB_MAX_OPEN|알람을 저장할 디비 동시 연결 수|int|X||postgres 문서 참조|>= 4.8.0|
|DB_MAX_IDLE|알람을 저장할 디비 유지 연결 수|int|X||postgres 문서 참조|>= 4.8.0|
|DB_TLS|데이터베이스 TLS 접속 여부|bool|X|false||>= 4.8.0|
|DB_CERT_FILE_PATH|데이터베이스 인증파일 경로|string|X||DB_TLS 값이 true이면 반드시 필요|>= 4.8.0|
|DB_CA_CERT|디비 root ca cert파일 이름|string|X||DB_TLS 값이 true이면 반드시 필요|>= 4.8.0|
|DB_CLIENT_CERT|디비 사용자 cert 파일 이름 |string|X||DB_TLS 값이 true이면 반드시 필요|>= 4.8.0|
|DB_CLIENT_KEY|디비 사용자 key 파일 이름|string|X||DB_TLS 값이 true이면 반드시 필요|>= 4.8.0|
|USE_KAFKA|카프카 사용 여부|bool|O|false||>= 4.8.0|
|KAFKA_BROKER_ADDRRESSES|카프카 브로커 주소(배열)|[]string|O||콤마로 구분된 목록을 입력|>= 4.8.0|
|KAFKA_TOPIC|카프카 토픽|[]string|O|audit-topic|콤마로 구분된 목록을 입력|>= 4.8.0|
|KAFKA_CONSUMER_GROUP|카프카 컨슈머 그룹|string|O|audit-consumer-group||>= 4.8.0|
|KAFKA_TLS|카프카 사용 여부|bool|O|false||>= 4.8.0|
|KAFKA_CERT_FILE_PATH|카프카 인증파일 저장 경로|string|X||KAFKA_TLS 값이 true이면 반드시 필요|>= 4.8.0|
|KAFKA_CA_CERT|카프카 root ca cert파일 이름|string|X||KAFKA_TLS 값이 true이면 반드시 필요|>= 4.8.0|
|KAFKA_CLIENT_CERT|카프카 사용자 cert 파일 이름|string|X||KAFKA_TLS 값이 true이면 반드시 필요|>= 4.8.0|
|KAFKA_CLIENT_KEY|카프카 사용자 key 파일 이름|string|X||KAFKA_TLS 값이 true이면 반드시 필요|>= 4.8.0|
|AUTH_SERVER_URL|칵테일 API 서버 URL|string|O|||>= 4.8.0|
|DO_ALARM_NOTIFY|알람 전달 여부 (정해진 알람 처리자에게 노티합니다.)|bool|X|false||>= 4.8.0|
|ALARM_NOTIFICATION_API_URL|알람 처리자 URL|string|X||DO_ALARM_NOTIFY가 true인 경우 반드시 필요|>= 4.8.0|
|LOGGING_LEVEL|프로그램 로그 레벨|string|X|info||>= 4.8.0|
|LOGGING_FILE_USE|프로그램 로그 저장 여부|bool|X|false||>= 4.8.0|
|LOGGING_FILE_PATH|프로그램 로그 저장 위치|string|X||LOGGING_FILE_USE 값이 true이면 반드시 필요|>= 4.8.0|
<br>
---
### version 4.6.3
|이름|설명|값|필수|기본값|참고|버전|
|------|---|---|---|---|---|--|
|SERVICE_PORT|서비스를 실행하는 host port 번호|int|X|9000||>= 4.6.3.x|
|AUTH_SERVER_URL|칵테일 API 서버 URL|string|O||http 방식 전송 시 인증을 위해 칵테일 서버와 통신.|>= 4.6.3.x|
|MONITORING_DB_URL|알람을 저장할 디비 ip 주소|string|O|||>= 4.6.3.x|
|MONITORING_DB|알람을 저장할 디비 이름|string|O|||>= 4.6.3.x|
|MONITORING_DB_USER|알람을 저장할 디비 사용자 아이디|string|O||보통 secret을 mount해서 사용한다.|>= 4.6.3.x|
|MONITORING_DB_PASSWORD|알람을 저장할 디비 사용자 암호|string|O||보통 secret을 mount해서 사용한다.|>= 4.6.3.x|
|DB_MAX_CONNECTION|알람을 저장할 디비 동시 연결 수|int|X|10|postgres 문서 참조|>= 4.6.3.x|
|DB_MAX_IDLE_CONNECTION|알람을 저장할 디비 유지 연결 수|int|X|10|postgres 문서 참조|>= 4.6.3.x|
|PROCESSING_MAP_PATH|설정 파일 마운트 경로|string|X||절대 경로를 입력합니다.|>= 4.6.3.x|
|LOGGING_LEVEL|로거 레벨|string|X|info||>= 4.6.3.x|
|REQUEST_TRACE|Request 상세 로깅 활성화|bool|X|false||>= 4.6.3.x|
|DEV_MODE|개발자 모드 활성화|bool|X|false||>= 4.6.3.x|
|DEV_PROFILING|서버 프로파일링 활성화|bool|X|false||>= 4.6.3.x|
|USE_KAFKA|카프카 사용 여부를 지정한다|bool|X|true|false로 지정하면 이전 방식의 수신만 가능|>= 4.6.3|
|KAFKA_BROKER_ADDRRESSES|카프카 서버 주소|string|X||주소가 하나 이상이면 콤마(,)로 분리, USE_KAFKA값이 true면 반드시 필요|>= 4.6.3|
|KAFKA_TOPIC|데이터를 읽어올 카프카 topic이름|sring|X||USE_KAFKA값이 true면 반드시 필요|>= 4.6.3|
|KAFKA_CONSUMER_GROUP|카프카 consumer group 이름|string|X||하나 이상의 consumer를 하는 경우 필요|>= 4.6.3|
|KAFKA_TLS|카프카 접속 시 SSL 사용 여부를 지정|bool|X|false||>= 4.6.3|
|KAFKA_CERT_FILE_PATH|SSL 인증서 파일 경로|string|X||USE_TLS 값이 true이면 반드시 필요|>= 4.6.3|
|KAFKA_CA_CERT|카프카 root ca cert파일 이름|string|X||USE_TLS 값이 true이면 반드시 필요|>= 4.6.3|
|KAFKA_CLIENT_CERT|카프카 사용자 cert파일 이름|string|X||USE_TLS 값이 true이면 반드시 필요|>= 4.6.3|
|KAFKA_CLIENT_KEY|카프카 사용자 key 파일 이름|string|X||USE_TLS 값이 true이면 반드시 필요|>= 4.6.3|
|MONITORING_DB_TLS|모니터링 디비 접속 시 SSL 사용 여부 지정|bool|X|false||>= 4.6.3|
|MONITORING_DB_CERT_FILE_PATH|디비 접속 시 사용할 인증서 저장 경로|string|X||MONITORING_DB_TLS 값이 true이면 반드시 필요|>= 4.6.3|
|MONITORING_DB_CA_CERT|디비 root ca cert파일 이름|string|X||MONITORING_DB_TLS 값이 true이면 반드시 필요|>= 4.6.3|
|MONITORING_DB_CLIENT_CERT|디비 사용자 cert 파일 이름 |string|X||MONITORING_DB_TLS 값이 true이면 반드시 필요|>= 4.6.3|
|MONITORING_DB_CLIENT_KEY|디비 사용자 key 파일 이름|string|X||MONITORING_DB_TLS 값이 true이면 반드시 필요|>= 4.6.3|
</br>
---

## version 4.8.0
### 주요 변경 내용
  * 코드 공통화 (아키텍처 변경) - grop사용

## version 4.6.3
### 주요 변경 내용
  * 모니터링 디비 SSL 접속 기능 추가
  * refactoring