# Petclinic

## About, 설명
[About](docs/about.md)

## Quick Start, 빠른 시작

```sh
# Test pv, for persistence, you need to connect an appropriate pv to the default/petclinic pvc.
# 테스트용 pv, 영속성을 위해서는 default/petclinic pvc에 적절한 pv를 연결해야 함.
$ kubectl apply -f https://raw.githubusercontent.com/rayshoo/spring-petclinic-data-jdbc/master/quickstart/example/pv.yaml

# Create petclinic manifest, forward ingress-nginx controller service port.
# petclinic 매니패스트 생성, ingress-nginx controller 서비스 포트포워드.
$ kubectl apply -f https://raw.githubusercontent.com/rayshoo/spring-petclinic-data-jdbc/master/quickstart/database/manifests.yaml && \
kubectl apply -f https://raw.githubusercontent.com/rayshoo/spring-petclinic-data-jdbc/master/quickstart/server/manifests.yaml && \
kubectl port-forward -n ingress-nginx --address=0.0.0.0 svc/ingress-nginx-controller 3000:80

$ curl -H "Host: petclinic.example.com" <hostIP>:3000

# Use quickstart/example/ing.yaml if you don't want host input.
# host 입력을 원하지 않을 경우, quickstart/example/ing.yaml 사용.
$ kubectl apply -f https://raw.githubusercontent.com/rayshoo/spring-petclinic-data-jdbc/master/quickstart/example/ing.yaml
$ curl <hostIP>:3000
```

## Build Requirements, 빌드 전 필요 사항
### .env file Setting, .env 파일 설정
```sh
$ wget https://raw.githubusercontent.com/rayshoo/spring-petclinic-data-jdbc/master/.env.example && \
wget https://raw.githubusercontent.com/rayshoo/spring-petclinic-data-jdbc/master/.env.mysql.example && \
mv .env.example .env && mv .env.mysql.example .env.mysql && ls -al
```
Create an imageSecret for the registry to pull and push images to.
이미지를 풀,푸시할 레지스트리의 imageSecret을 미리 생성한다.</br>
[imageSecret 생성하기 | To create imageSecret](docs/imageSecret.md)

```yaml
# For scalability, the password environment variable used by each app is separated, but the same password is used in practice.
# 확장성을 위해 각 앱이 사용하는 Password 환경 변수를 분리했으나, 실습에서는 동일한 비밀번호를 사용하도록 한다.
# .env
IMAGE=<registry/imageName>
IMAGE_REPO_SECRET=<imageSecretName>

MYSQL_PASS=<mysql pass> # Use same password!

# .env.mysql
MYSQL_ROOT_PASSWORD=<mysql pass> # Use same password!
```

## Build, 빌드
### Without Docker Engine - on K8S
```sh
# https://github.com/rayshoo/spring-petclinic-data-jdbc/tree/master/k8s
# Downloading petclinic binary files. It creates petclinic k8s manifest files according to envs.
# petclinic 바이너리 파일 다운로드, 환경변수에 따라 petclinic k8s 매니패스트 파일들을 생성해주는 도구.
$ VERSION=v1.0.1
$ OS=<linux|windows|darwin>
$ ARCH=<amd64|arm64>
$ wget https://github.com/rayshoo/spring-petclinic-data-jdbc/releases/download/$VERSION/petclinic-$OS-$ARCH -O ./petclinic && \
chmod +x ./petclinic

# Run petclinic binary with .env, .env.mysql file in current path, Create manifests.yaml file with stdout.
# 현재 경로에 .env, .env.mysql 파일이 있는 상태에서 petclinic 바이너리 실행, stdout으로 manifests.yaml 파일 생성.
$ ls .env .env.mysql
-rw-rw-r-- .env
-rw-rw-r-- .env.mysql

$ ./petclinic > manifests.yaml

# Wait for the build to complete. Note that it is executed synchronously in the order of base > builder > build.
# 빌드가 완료될때까지 대기. base > builder > build 순서로 동기적으로 실행됨에 유의한다.
$ kubectl apply -f manifests.yaml && \
watch kubectl get pods

# If the mysql pod is in pending state, make sure you have connected the appropriate pv. Create pv for testing if there is none.
# mysql pod가 pending 상태일 경우, 적절한 pv를 연결했는지 확인한다. 없을 시 테스트용 pv 연결.
$ kubectl get pods -l app=petclinic-mysql
$ kubectl apply -f https://raw.githubusercontent.com/rayshoo/spring-petclinic-data-jdbc/master/quickstart/example/pv.yaml

# https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#restart-policy
# When pod creation fails, the restart time increases as time elapses, so restart the deployment using the image you just built.
# 파드 생성 실패 시 시간이 경과할 수록 재시작 시간이 늘어나기 때문에 방금 전 빌드한 이미지를 사용하는 deployment를 재시작해준다.
$ kubectl rollout restart deployment petclinic

# Check routing through Ingress Controller
# 인그레스 컨트롤러를 통한 라우팅 확인
$ kubectl port-forward -n ingress-nginx --address=0.0.0.0 svc/ingress-nginx-controller 3000:80
$ curl -H "Host: petclinic.example.com" <hostIP>:3000
```

### With Docker Engine - on Docker
[docker-compose 설치](https://github.com/docker/compose/releases)
```sh
# Docker login to the registry to pull and push the image.
# 이미지를 풀,푸시할 레지스트리에 도커 로그인을 한다.
$ echo <registry password> | docker login [<registry name>] -u <registry username> --password-stdin

$ git clone https://github.com/rayshoo/spring-petclinic-data-jdbc.git && cd spring-petclinic-data-jdbc

# Cleans up the path where mysql data is mounted.
# mysql 데이터가 마운트 되는 경로를 정리해준다.
$ sudo cp -r /mnt/petclinic/mysql /mnt/petclinic/mysql.backup && \
sudo rm -rf /mnt/petclinic/mysql

# Run the following command in the path where the docker-compose.yml, .env, .env.mysql, Makefile files are located
# docker-compose.yml, .env, .env.mysql, Makefile 파일이 위치한 경로에서 다음 명령어 실행
$ ls -al docker-compose.yml .env .env.mysql Makefile
-rw-rw-r-- docker-compose.yml
-rw-rw-r-- .env
-rw-rw-r-- .env.mysql
-rw-rw-r-- Makefile

$ docker-compose build base && docker-compose push base && \
docker-compose build --no-cache was_builder && docker-compose up was_builder && \
docker-compose down && docker-compose pull mysql was && \
docker-compose up -d mysql was

# If make is installed, you can enter the make command instead of the docker-compose command.
# make가 설치되어 있다면, docker-compose 명령어 대신 make 명령어를 입력해도 된다.
$ make

$ curl localhost:8080
```