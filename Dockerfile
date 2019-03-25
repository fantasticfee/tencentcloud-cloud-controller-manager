FROM ubuntu:16.04

RUN apt-get update
RUN apt-get install -y  ca-certificates
ADD tencentcloud-cloud-controller-manager /bin/

CMD ["/bin/tencentcloud-cloud-controller-manager"]
