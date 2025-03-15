FROM gcr.io/distroless/base-debian12:debug
COPY ./bin/product-store /usr/local/bin/product-store
ENTRYPOINT [ "/usr/local/bin/product-store" ]
