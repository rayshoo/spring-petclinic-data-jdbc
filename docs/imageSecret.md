# imageSecret 생성하기 | To create imageSecret

```sh
$ kubectl create secret -n default docker-registry <IMAGE_REPO_SECRET> \
[--docker-server=<registry server> ] \
--docker-username=<repo username> \
--docker-password=<repo password>
```