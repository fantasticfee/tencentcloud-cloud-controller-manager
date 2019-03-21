FROM ubuntu:16.04

ADD tencentcloud-cloud-controller-manager /bin/

CMD ["/bin/tencentcloud-cloud-controller-manager"]
