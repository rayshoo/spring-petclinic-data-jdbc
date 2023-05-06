FROM khipu/openjdk17-alpine

RUN addgroup --system --gid 1000 spring && \
adduser --system --uid 1000 --ingroup spring --disabled-password spring && \
mkdir /var/log/petclinic && \
chown -R spring:spring /var/log/petclinic
USER spring
WORKDIR /var/log/petclinic
