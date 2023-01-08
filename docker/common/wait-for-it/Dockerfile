FROM busybox:glibc as build

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait

FROM scratch

COPY --from=build /wait ./

CMD ["/wait"]
