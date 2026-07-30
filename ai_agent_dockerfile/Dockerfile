ARG PYTHON_VERSION=3.13
FROM python:${PYTHON_VERSION}-slim
COPY ./requirements.txt /build/requirements.txt
RUN mkdir /root/.jupyter
COPY ./jupyter_pydantic_ai_config.yaml /root/.jupyter/jupyter_pydantic_ai_config.yaml
RUN pip install -r /build/requirements.txt
