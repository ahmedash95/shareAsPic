FROM golang:1.11.2

# Installing chrome browser
RUN \
  wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - && \
  echo "deb http://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google.list && \
  apt-get update && \
  apt-get install -y --no-install-recommends google-chrome-stable

# Installing fonts to support utf-8 for (Arabic, Indian, Chainees)
RUN apt-get install --no-install-recommends -y software-properties-common \
    fonts-indic fonts-noto fonts-noto-cjk \
    && rm -rf /var/lib/apt/lists/* /etc/apt/sources.list.d/google.list

WORKDIR /go/src/github.com/ahmedash95/shareAsPic
COPY . .

RUN go get -d -v ./...

RUN go install -v ./...
