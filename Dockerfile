FROM golang:1.13

# Installing chrome browser
RUN \
  wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - && \
  echo "deb http://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google.list && \
  apt update && \
  apt install -y google-chrome-stable && \
rm -rf /var/lib/apt/lists/*

# Installing fonts to support utf-8 for (Arabic, Indian, Chinese)
RUN apt update
RUN apt install software-properties-common -y
RUN apt install fonts-indic fonts-noto fonts-noto-cjk -y

RUN echo GOPATH
RUN mkdir -p $GOPATH/src/app
 
WORKDIR /app

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
