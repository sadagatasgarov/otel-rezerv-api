FROM golang:1.21.5-alpine3.18 as builder
WORKDIR /app
COPY . .

#RUN go mod tidy
#COPY *.go .
RUN go build -o ./app
CMD ["chmod", "+x", "./app/app" ]


FROM scratch
COPY --from=builder /app/app /app
EXPOSE 5000
CMD [ "./app" ]docker 