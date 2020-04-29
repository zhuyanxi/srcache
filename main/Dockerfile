FROM golang:1.14.2-alpine AS build
RUN apk --no-cache add ca-certificates 
WORKDIR /app
COPY . .
RUN http_proxy=http://proxy-mtl.ubisoft.org:3128 CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o srcache

FROM scratch AS final
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/srcache /app/srcache
EXPOSE 3001
ENTRYPOINT [ "/app/srcache" ]