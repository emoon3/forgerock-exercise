FROM --platform=linux/amd64 alpine:3.16
WORKDIR /app
COPY main /app
ENV SYMBOL=MSFT
ENV API_KEY=C227WD9W3LUVKVV9
ENV NUM_DAYS=7
RUN chmod +X /app/main
CMD /app/main