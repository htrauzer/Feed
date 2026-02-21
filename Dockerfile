#on top of the image "golang", version 1.17, do these commands:
FROM golang:1.17

#create a folder inside the image
RUN mkdir /app

#Set the working directory for any subsequent ADD, COPY, CMD, ENTRYPOINT, or RUN instructions that follow it in the Dockerfile.
WORKDIR /app

#copy all files from the Dockerfile directory and all subfolders to WORKDIR (inside image)
ADD . .

# --- gets all imported packages ---
RUN go mod download

# --- documents the port(s) this containerized app is listening on ---
EXPOSE 8080

#build executable binary from our go files and output it to (flag -o) file called foorum.exe. Dot represent all files and folders.
RUN go build -o forum .

#when image is run, exec this file:
CMD ["/app/forum"]

