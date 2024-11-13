FROM golang:1.21.4-alpine3.18 AS builder

ENV GOPROXY=https://goproxy.cn

WORKDIR /src

COPY . .

RUN go build -o /bin ./cmd/api

FROM alpine

WORKDIR /usr/src/app/

RUN apk add --no-cache ca-certificates tzdata

ENV TZ=Asia/Shanghai

COPY --from=builder /bin ./
COPY --from=builder /src/configs ./

EXPOSE 9000

CMD [ "./api","-f","./configs/config.local.yaml" ]