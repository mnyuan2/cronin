FROM busybox:glibc
COPY ./configs/ /cronin/configs
COPY ./data/ /cronin/data
COPY ./cronin /cronin
WORKDIR /cronin
EXPOSE 9003
CMD ./cronin