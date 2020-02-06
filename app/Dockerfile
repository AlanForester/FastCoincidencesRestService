FROM golang

# if left blank app will run with dev settings
# to build production image run:
# $ docker build ./app --build-args app_env=production
ARG app_env
ENV APP_ENV $app_env

COPY . /go/src/app
WORKDIR /go/src/app

## added vendor modules
RUN go mod vendor
#
RUN go build -o ./bin/app

CMD if [ $APP_ENV = production ]; \
	then \
	chmod +x ./bin/app; \
	./bin/app; \
	else \
	go install github.com/pilu/fresh && \
	fresh; \
	fi

EXPOSE 8080