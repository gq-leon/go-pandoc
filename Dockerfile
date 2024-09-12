FROM ubuntu:latest
LABEL authors="vagrant"

ENTRYPOINT ["top", "-b"]