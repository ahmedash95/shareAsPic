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
# all fonts from universe 
RUN apt install -y fonts-dejavu-core fonts-dejavu-extra fonts-droid-fallback fonts-guru fonts-guru-extra fonts-horai-umefont fonts-kacst fonts-kacst-one fonts-lao fonts-liberation fonts-lklug-sinhala fonts-lohit-guru fonts-nanum fonts-noto-cjk fonts-opensymbol fonts-roboto fonts-roboto-hinted fonts-sil-abyssinica fonts-sil-padauk fonts-stix fonts-symbola fonts-thai-tlwg fonts-tibetan-machine fonts-tlwg-garuda fonts-tlwg-kinnari fonts-tlwg-laksaman fonts-tlwg-loma fonts-tlwg-mono fonts-tlwg-norasi fonts-tlwg-purisa fonts-tlwg-sawasdee fonts-tlwg-typewriter fonts-tlwg-typist fonts-tlwg-typo fonts-tlwg-umpush fonts-tlwg-waree fonts-unfonts-core

RUN echo GOPATH
RUN mkdir -p $GOPATH/src/app
 
WORKDIR /app

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
