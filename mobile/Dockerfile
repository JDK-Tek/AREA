FROM ghcr.io/cirruslabs/flutter:3.27.1 AS build

RUN apt-get update && apt-get install -y openjdk-17-jdk

ENV JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64
ENV PATH=$JAVA_HOME/bin:$PATH

WORKDIR /app

COPY . .

RUN flutter doctor -v
RUN flutter pub get
RUN flutter pub upgrade --major-versions

RUN flutter build apk --release
RUN mkdir -p /apk_volume && cp build/app/outputs/flutter-apk/app-release.apk /apk_volume/area_jepgo.apk