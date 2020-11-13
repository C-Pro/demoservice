FROM scratch

COPY bin/service /service

CMD ['/service']
