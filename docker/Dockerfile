FROM python:3.10-slim

ARG USER=hefactory
ARG GROUP=hefactory
ARG UID=1000
ARG GID=1000
ARG PW=docker_hefactory
ARG ENV_NAME=base
ENV DEBIAN_FRONTEND="noninteractive" TZ="Europe/Madrid"

# Base image is Ubuntu and needed tools for Linux usage
# Getting python and related tools
RUN apt-get update -y && \
apt-get install -y vim git wget curl unzip sudo 


RUN groupadd -g ${GID} ${GROUP}
RUN useradd -s /bin/bash -d /home/${USER} --uid=${UID} -m -g sudo ${USER} && echo "${USER}:${PW}" | chpasswd


RUN echo "PS1='\[\033[01;31m\](docker) \[\e[01;36m\]\u\[\e[01;37m\]@\[\e[01;33m\]\H\[\e[01;37m\]:\[\e[01;32m\]\w\[\e[01;37m\]\$\[\033[0;37m\] '" >> /root/.bashrc && \
echo "alias python=python3" >> /root/.bashrc && \
echo "alias pip=pip3" >> /root/.bashrc && \
echo "export PATH=/usr/local/go/bin:/home/${USER}/miniconda3/bin:$PATH" >> /root/.bashrc

RUN echo "PS1='\[\033[01;31m\](docker) \[\e[01;36m\]\u\[\e[01;37m\]@\[\e[01;33m\]\H\[\e[01;37m\]:\[\e[01;32m\]\w\[\e[01;37m\]\$\[\033[0;37m\] '" >> /home/${USER}/.bashrc && \
echo "alias python=python3" >> /home/${USER}/.bashrc && \
echo "alias pip=pip3" >> /home/${USER}/.bashrc && \
echo "export PATH=/usr/local/go/bin:/home/${USER}/miniconda3/bin:$PATH" >> /home/${USER}/.bashrc

RUN mkdir -p /workspace/HEFactory/

COPY examples /workspace/HEFactory/examples/

RUN chown -R ${USER}:${GROUP} /workspace

USER ${USER}:${GID}
WORKDIR /home/${USER}

COPY dist/HEFactory-0.0.1-cp310-cp310-linux_x86_64.whl /tmp/HEFactory-0.0.1-cp310-cp310-linux_x86_64.whl
RUN python3 --version && pip3 install /tmp/*.whl

WORKDIR /workspace/HEFactory/


CMD ["bash"]
