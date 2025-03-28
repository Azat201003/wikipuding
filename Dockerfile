FROM alpine

WORKDIR /app
COPY /wikipuding .

# EXPOSE 8080
RUN chmod +x /app/wikipuding
# RUN apk update --no-cache && apk add --no-cache tzdata
# RUN apk add git make libc6-compat

EXPOSE 8080

CMD ./wikipuding
