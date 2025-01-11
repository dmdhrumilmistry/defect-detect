# Build application
FROM golang:latest AS build

WORKDIR /go/src/defect-detect
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/defect-detect .

# Copy app from build stage
FROM gcr.io/distroless/static-debian12
COPY --from=build /go/bin/defect-detect /
CMD ["/defect-detect"]
