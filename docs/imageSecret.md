# imageSecret 생성하기 | To create imageSecret
## Main, 메인
[Main](../README.md)

## Command, 커맨드

```sh
$ kubectl create secret -n default docker-registry <IMAGE_REPO_SECRET> \
[--docker-server=<registry server> ] \
--docker-username=<repo username> \
--docker-password=<repo password>
```