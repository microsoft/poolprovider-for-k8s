FROM ubuntu:16.04

# To make it easier for build and release pipelines to run apt-get,
# configure apt to not require confirmation (assume the -y argument by default)
ENV DEBIAN_FRONTEND=noninteractive
RUN echo "APT::Get::Assume-Yes \"true\";" > /etc/apt/apt.conf.d/90assumeyes

RUN apt-get update \
&& apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        jq \
        git \
        iputils-ping \
        libcurl3 \
        libicu55 \
        libunwind8 \
        netcat

RUN /usr/bin/curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl \
     && chmod +x ./kubectl  \
     &&  mv ./kubectl /usr/local/bin/kubectl

RUN /usr/bin/curl -LO https://github.com/genuinetools/img/releases/download/$(curl -s "https://api.github.com/repos/genuinetools/img/releases/latest" | jq -r .tag_name)/img-linux-amd64 \
     && chmod +x ./img-linux-amd64 \
     && mv ./img-linux-amd64 /usr/local/bin/img

RUN BUILDKITVERSION=$(curl -s "https://api.github.com/repos/moby/buildkit/releases/latest" | jq -r .tag_name) \
     && curl -LO https://github.com/moby/buildkit/releases/download/${BUILDKITVERSION}/buildkit-${BUILDKITVERSION}.linux-amd64.tar.gz \
     && tar -zxvf buildkit-${BUILDKITVERSION}.linux-amd64.tar.gz \
     && chmod +x ./bin/buildctl && chmod +x ./bin/buildkitd && chmod +x ./bin/buildkit-runc \
     && mv ./bin/buildctl /usr/local/bin/ && mv ./bin/buildkitd /usr/local/bin/ && mv ./bin/buildkit-runc /usr/local/bin/

# setting work directory
WORKDIR /azp

COPY ./start.sh .
RUN chmod +x start.sh

CMD ["./start.sh"]