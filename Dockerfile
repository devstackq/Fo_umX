# FROM golang:latest AS Builder

# # enable Go modules support
# ENV GO111MODULE=on
# #RUN mkdir /app
# WORKDIR /app

# # manage dependencies
# COPY go.mod .
# COPY go.sum .   
# RUN go mod download
# COPY . .

# #RUN  go build -o forumx
# # CGO_ENABLED=0 GOOS=linux GOARCH=amd64
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main 

# FROM  scratch
# COPY --from=Builder /app/forumx /app
# EXPOSE 6969
# ENTRYPOINT ["/app/forumx"]

FROM golang:latest
RUN mkdir /app
ADD . /app/ 
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]