FROM khipu/openjdk17-alpine
COPY . .
RUN mkdir /root/.docker
ENTRYPOINT ["./gradlew"]
CMD ["jib"]