FROM python:3.10-slim

WORKDIR /usr/app

RUN apt-get update && \
    apt-get install -y libpq-dev gcc

COPY requirements.txt /usr/app
RUN pip3 install -r requirements.txt

COPY app.py /usr/app
CMD python3 app.py
