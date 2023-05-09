# About, 설명
## Main, 메인
[Main](../README.md)

1. Gradle로 어플리케이션, 도커이미지 빌드
```sh
# 어플리케이션 빌드
./gradle build

# 도커이미지 빌드
./gradlew jib
```
jib 세팅은 [build.gradle](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/master/was/build.gradle) [#11](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/was/build.gradle#L11), [#69](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/was/build.gradle#L69) line에 작성</br>
*jib를 사용한 이유: 도커 데몬없이 컨테이너 이미지 빌드를 쉽게 빌드할 수 있다는 점이 매력적으로 느껴져서 궁금했는데, 이번 기회에 사용함.*

2. 어플리케이션 log - host /logs 디렉터리에 적재
로그 설정: [logback-spring.xml](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/master/was/src/main/resources/logback-spring.xml)</br>
host /logs 디렉터리 적재: [quickstart](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/quickstart/server/manifests.yaml#L73), [template](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/k8s/templates/deploy.yaml#L89)

3. 정상 작동 api 구현, 10초마다 체크
정상 작동 api 구현: [health probe enable](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/was/src/main/resources/application.yaml#L39)</br>
10초마다 체크: [quickstart](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/quickstart/server/manifests.yaml#L61), [template](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/k8s/templates/deploy.yaml#L77)</br>
*perionSeconds 기본값이 10이기에 readiness는 작성 생략*

4. 3번 연속 체크 실패시, 어플리케이션 restart
[quickstart](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/quickstart/server/manifests.yaml#L60), [template](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/k8s/templates/deploy.yaml#L76)

5. 종료 시 30초 이내 프로세스 종료 안될 시, SIGKILL 강제 종료
server: [appication.yaml](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/was/src/main/resources/application.yaml#L7)</br> 
k8s: [quickstart](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/quickstart/server/manifests.yaml#L39), [template](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/k8s/templates/deploy.yaml#L55)

6. 배포, scale in/out 유실 트래픽 방지
- 기 liveness, readiness 설정
- graceful shutdown setting: [application.yaml](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/was/src/main/resources/application.yaml#L3)
- k8s rolling update startegy: [quickstart](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/quickstart/server/manifests.yaml#L27), [template](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/k8s/templates/deploy.yaml#L43)

7. 어플리케이션 프로세스 uid:1000으로 실행
[quickstart](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/quickstart/server/manifests.yaml#L41), [template](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/k8s/templates/deploy.yaml#L57)

8. DB kubernetes 실행, 재 실행시 변경된 데이터 유실 방지
volume mount setting: [quickstart](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/quickstart/database/manifests.yaml#L50), [template](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/k8s/templates/deploy.yaml#L27)</br>
volume: [quickstart](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/quickstart/database/manifests.yaml#L54), [template](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/k8s/templates/deploy.yaml#L31)</br>
*pvc 설정 - 적절한 stroage를 연결하도록 유도*

9. 어플리케이션 - DB cluster domain 통신
- DB Headless svc: [quickstart](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/quickstart/database/manifests.yaml#L65), [template](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/k8s/templates/svc.yaml#L11)
- To DB URL: [application.yaml](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/was/src/main/resources/application.yaml#L14), [quickstart](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/quickstart/server/manifests.yaml#L12), [template](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/k8s/templates/secret.yaml#L22)
*어플리케이션 소스에서는 환경변수를 읽어들이도록 처리하고, 루트디렉터리의 [.env](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/cd433335fb563ee4cbfa441a8e9b48cd621876e8/.env.example#L8)에 따라 templating해서 [manifest파일을 생성](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/master/README.md#without-docker-engine---on-k8s)하도록 유도*

10. nginx-ingress-controller를 통해 어플리케이션 접속
ingress: [quickstart](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/master/quickstart/example/ing.yaml), [template](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/master/k8s/files/ing.yaml)</br>
*nginx-ingress-controller svc가 외부로 어떻게 노출되어있는지를 몰라서 서비스를 [포트포워드하고 접속하도록 README 작성](https://github.com/rayshoo/spring-petclinic-data-jdbc/blob/master/README.md#quick-start-%EB%B9%A0%EB%A5%B8-%EC%8B%9C%EC%9E%91)함*

11. namespace default 사용
manifests 모든 리소스에 default 네임스페이스 명시적으로 작성.

12. README.md 파일 실행 방법 기술
[README.md](../README.md)