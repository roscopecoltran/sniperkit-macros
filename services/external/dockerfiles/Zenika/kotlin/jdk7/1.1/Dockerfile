FROM openjdk:7-jdk-alpine

RUN apk add --no-cache bash && \
    apk add --no-cache -t build-dependencies wget && \
    cd /usr/lib && \
    wget https://github.com/JetBrains/kotlin/releases/download/v1.1.4-2/kotlin-compiler-1.1.4-2.zip && \
    unzip kotlin-compiler-1.1.4-2.zip && \
    rm kotlin-compiler-1.1.4-2.zip && \
    rm kotlinc/bin/*.bat && \
    apk del --no-cache build-dependencies

ENV PATH $PATH:/usr/lib/kotlinc/bin

CMD ["kotlinc"]
