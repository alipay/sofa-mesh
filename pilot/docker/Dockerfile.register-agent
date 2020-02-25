FROM ubuntu:xenial

RUN mkdir -p ./config
ADD agent.yaml ./config/
ADD register-agent /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/register-agent"]