ARG PYTHON_VERSION=3.12
FROM python:${PYTHON_VERSION}-slim
COPY ./requirements.txt /build/requirements.txt
RUN pip install -r /build/requirements.txt
RUN python -m bash_kernel.install
