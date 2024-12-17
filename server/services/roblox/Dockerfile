FROM python:3.10-alpine

WORKDIR /usr/app

COPY requirements.txt /usr/app
RUN pip3 install -r requirements.txt

COPY app.py /usr/app
CMD python3 app.py
