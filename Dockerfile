FROM --platform=$BUILDPLATFORM golang:alpine AS build

RUN apk add --no-cache git

RUN git clone https://github.com/TheTipo01/videoDownloader /videoDownloader
WORKDIR /videoDownloader
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go mod download
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o videoDownloader

FROM alpine

RUN apk add --no-cache ffmpeg python3

RUN wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -O /usr/bin/yt-dlp
RUN chmod a+rx /usr/bin/yt-dlp

COPY --from=build /videoDownloader/videoDownloader /usr/bin/

CMD ["videoDownloader"]
