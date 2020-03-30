FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/poolprovider-for-k8s \
    USER_UID=1001 \
    USER_NAME=poolprovider-for-k8s

# install operator binary
COPY build/_output/bin/poolprovider-for-k8s ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
