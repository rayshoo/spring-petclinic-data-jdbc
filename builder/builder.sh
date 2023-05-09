#! /bin/sh
git clone https://github.com/rayshoo/spring-petclinic-data-jdbc.git petclinic
cd petclinic/was/
chmod +x ./gradlew
./gradlew jib