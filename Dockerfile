FROM python:3.10.4-slim-buster
ENV PYTHONBUFFERED=1

RUN apt-get update && \
    rm -rf /var/lib/apt/lists/*

RUN pip install pipenv
COPY Pipfile /app/

WORKDIR /app

RUN pipenv lock --keep-outdated --requirements > requirements.txt && \
	pip install --no-cache-dir -r requirements.txt

COPY . .

CMD [ "python3", "-m" , "flask", "run"]
