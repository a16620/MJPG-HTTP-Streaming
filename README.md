# MJPG HTTP Streaming
 HTTP 서버를 통해 JPEG 이미지를 지속적으로 송출합니다. 오디오가 필요하지 않을 때 CCTV나 IP카메라 용도로 사용할 수 있습니다. 
 
# 요구 사항
 실행 환경에 golang과 openCV, gocv가 설치되어있어야 합니다.

# camera.go
 Camera는 일정 간격마다 캡쳐해 수신 대기중인 채널에 Fan-out시킵니다. select와 time.After를 통해 지금 대기중이 아닌 채널은 스킵하게 해뒀습니다. 대신 오래동안 채널을 사용하지 않았을때 채널 버퍼에 한참 전의 이미지가 포함되어있을 수 있습니다. 물론 몇번의 수신 후에는 다시 현재의 이미지를 수신할 수 있습니다.
