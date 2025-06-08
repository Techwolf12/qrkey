FROM scratch
COPY qrkey /usr/bin/qrkey
ENTRYPOINT ["/usr/bin/qrkey"]