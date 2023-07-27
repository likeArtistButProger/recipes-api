FROM golang:1.18

WORKDIR /usr/recipes_api
COPY . .

EXPOSE 8000

WORKDIR /usr/recipes_api/src/
RUN ["go", "build", "main.go"]

CMD ./main
