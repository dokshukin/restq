FROM        scratch
MAINTAINER  Ilya Dokshukin <dokshukin@gmail.com>

COPY restq-amd64 /usr/bin/restq

USER        nobody
EXPOSE      8080
ENTRYPOINT  [ "/usr/bin/restq" ]