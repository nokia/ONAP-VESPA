FROM fedora

RUN dnf makecache && dnf update -y
RUN dnf install -y rpm-build make git gcc

RUN curl https://dl.google.com/go/go1.11.4.linux-amd64.tar.gz | tar -C /usr/local -xz

RUN mkdir /go
RUN chmod 777 /go
ENV GOPATH=/go
ENV GOROOT=/usr/local/go

ENV PATH=${PATH}:${GOPATH}/bin:${GOROOT}/bin


RUN groupadd -r jenkins
RUN useradd -r -g jenkins -u 1000 -s /sbin/nologin -c "jenkins services" jenkins
# USER jenkins