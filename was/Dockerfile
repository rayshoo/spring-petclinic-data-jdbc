FROM alpine:3.17.3 as builder
RUN apk add git && \
git clone https://github.com/rayshoo/spring-petclinic-data-jdbc.git petclinic

FROM khipu/openjdk17-alpine
COPY --from=builder petclinic/was .
RUN mkdir /root/.docker
ENTRYPOINT ["./gradlew"]
CMD ["jib"]