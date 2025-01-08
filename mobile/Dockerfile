FROM ghcr.io/cirruslabs/flutter:3.27.1 AS build

RUN apt-get update && apt-get install -y openjdk-17-jdk

ENV JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64
ENV PATH=$JAVA_HOME/bin:$PATH

WORKDIR APP

COPY . .

RUN flutter doctor -v
RUN flutter pub get

RUN flutter build apk --release