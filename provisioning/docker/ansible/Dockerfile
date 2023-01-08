FROM alpine:3.12

ENV ANSIBLE_VERSION=2.13.6

RUN apk --no-cache add \
        sudo \
        python3 \
        py3-pip \
        openssl \
        ca-certificates \
        sshpass \
        openssh-client \
        rsync \
        git \
    && apk --no-cache add --virtual build-dependencies \
        python3-dev \
        libffi-dev \
        openssl-dev \
        build-base \
    && pip3 install --upgrade pip \
    && pip3 install ansible-core==${ANSIBLE_VERSION} \
    && ansible-galaxy collection install ansible.posix \
#    && pip3 install --upgrade cffi \
#    && pip3 install mitogen ansible-lint jmespath \
#    && pip3 install --upgrade pywinrm \
#    && pip3 install passlib \
    && apk del build-dependencies \
    && rm -rf /var/cache/apk/*

RUN mkdir /ansible && \
    mkdir -p /etc/ansible && \
    echo 'localhost' > /etc/ansible/hosts

COPY ./ssh_keys /root/ssh_keys
RUN chmod -R 600 /root/ssh_keys/ && chmod 750 /root/ssh_keys/

COPY ./ansible /ansible
WORKDIR /ansible

CMD ["ansible-playbook", "--version"]
