FROM golang AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -o /contactus


FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /contactus .
EXPOSE 8888
USER nonroot:nonroot
ENTRYPOINT ["/contactus"]
