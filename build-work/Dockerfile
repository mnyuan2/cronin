FROM busybox:latest
COPY cronin /biz-code/cronin/
COPY configs/ /biz-code/cronin/configs
WORKDIR /biz-code
CMD cd ./cronin && ./cronin
EXPOSE 9003