FROM golang:bullseye

RUN apt update
RUN apt install -y ffmpeg tap-plugins imagemagick libmagickcore-dev libmagickwand-dev build-essential

WORKDIR /app
COPY . .
RUN go build

CMD ["/app/vkbot"] 
