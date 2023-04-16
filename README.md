# study_rabbitmq
rabbitmq를 통한 이벤트 송수신 테스트

https://www.notion.so/RabbitMQ-6ea7c13a635f4eddb1e79b10a8bee020 


exchage
fanout 타입

queue
2개의 큐를 생성 위 exchange 바인딩

producer
HTTP API(PUT)을 통해 큐에 메시지를 저장.

customer
스레드로 동작 바인딩된 큐에 이벤트 발생시 수신 후 메시지를 로그로 출력


customer 로직 제거(주석처리)후 재실행시 message가 queue에 해소되지 않은채 쌓여있는 것을 확인.
management console(localhost:15672)로 접속하여 모니터링.
