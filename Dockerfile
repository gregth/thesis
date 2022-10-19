FROM debian:bullseye

RUN apt-get update \
&&  apt-get install -y \
      build-essential \
      curl \
      git \
      gosu \
      python3 \
      sudo \
      texlive \
      texlive-xetex \
      texlive-lang-greek \
      texlive-science \
      texlive-extra-utils \
      vim \
      wget

COPY /files/entrypoint.py /usr/local/bin

ENTRYPOINT ["python3", "/usr/local/bin/entrypoint.py"]
