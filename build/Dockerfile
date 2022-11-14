FROM golang:1.18

RUN apt-get update && \
    apt-get install -y unzip curl jq vim unzip less \
     && rm -rf /var/lib/apt/lists/* 


# setup user
RUN useradd -ms /bin/bash developer
USER developer
WORKDIR /home/developer/kubefirst

COPY --chown=developer:developer . .

# download dependencies and prepare fresh installation
RUN go mod download
RUN go build .

USER root
RUN  mv ./k1-watcher /usr/local/bin/ && chmod 777 /usr/local/bin/k1-watcher
USER developer
