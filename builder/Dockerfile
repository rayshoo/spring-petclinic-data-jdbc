FROM khipu/openjdk17-alpine
COPY builder.sh .
RUN apk add git && chmod +x ./builder.sh
ENTRYPOINT ["./builder.sh"]