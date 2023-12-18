FROM golang:1.21.5-alpine3.18 as builder
WORKDIR /app
COPY . .


# Eger islemese "go mod vendor" ile lokalda vendoru yaradib birbasa gonderin built almaq daha mentiklidir
RUN go mod download


RUN go build -o ./app
CMD ["chmod", "+x", "./app/app" ]


FROM scratch
COPY --from=builder /app/app /app
EXPOSE 5000
CMD [ "./app" ]docker 