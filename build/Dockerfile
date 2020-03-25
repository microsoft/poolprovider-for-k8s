FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/k8s-poolprovider \
    USER_UID=1001 \
    USER_NAME=k8s-poolprovider

# install operator binary
COPY build/_output/bin/k8s-poolprovider ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
