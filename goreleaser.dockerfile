FROM scratch
COPY depends-on /depends-on
ENTRYPOINT ["/depends-on"]