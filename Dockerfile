FROM golang:1.22-rc-bookworm

RUN apt-get -y update && apt-get -y install libfontconfig

WORKDIR /tmp
RUN wget https://mirror.ctan.org/systems/texlive/tlnet/install-tl-unx.tar.gz
RUN zcat < install-tl-unx.tar.gz | tar xf -
# TODO: update year when the next millenium comes
RUN cd install-tl-2* && perl install-tl --no-interaction --scheme=small
ENV PATH="/usr/local/texlive/2023/bin/x86_64-linux:${PATH}"
RUN tlmgr install enumitem titlesec

WORKDIR /app
COPY . /app

RUN mkdir -p /outputs /inputs

RUN GCO_ENABLED=0 go build -o generate-resumes main.go

ENTRYPOINT [ "./generate-resumes" ]
CMD ["examples/example.toml", "-o", "/outputs"]

