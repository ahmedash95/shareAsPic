FROM golang:1.11.2

# Installing chrome browser
RUN \
  wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - && \
  echo "deb http://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google.list && \
  apt-get update && \
  apt-get install -y google-chrome-stable && \
rm -rf /var/lib/apt/lists/*

# Installing fonts to support utf-8 for (Arabic, Indian, Chainees)
RUN apt-get update
RUN apt-get install software-properties-common -y
RUN apt-get install fonts-indic fonts-noto fonts-noto-cjk -y

WORKDIR /go/src/github.com/ahmedash95/shareAsPic
COPY . .

RUN go get -d -v ./...

RUN go install -v ./...