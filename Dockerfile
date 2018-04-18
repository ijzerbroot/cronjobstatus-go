
FROM ubuntu:xenial

RUN mkdir /home/ubuntu

ADD jobstatus /jobstatus
ADD entrypoint.sh /entrypoint.sh

ADD test/etl.log /home/ubuntu/etl.log
ADD test/backupportainer.log /tmp/backupportainer.log
ADD test/backupswarm.log /tmp/backupswarm.log

RUN chmod a+rx /jobstatus /entrypoint.sh
ENTRYPOINT /entrypoint.sh

