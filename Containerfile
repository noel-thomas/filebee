FROM quay.io/fedora/fedora:35

RUN dnf -y update

RUN dnf -y install python3-pip && dnf clean all

EXPOSE 8000
# file store directory
RUN mkdir -p /var/filebee

COPY ./requirements.txt /app/requirements.txt

WORKDIR /app

RUN pip3 install -r requirements.txt

COPY ./app.py /app/main.py

ENTRYPOINT [ "python" ]

CMD [ "/app/main.py" ]


