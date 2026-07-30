# KBox

KBox helps to run container from jupyter frontends (e.g. notebooks), and talk to the container through this frontend. This allows for :

- containerized code execution (of any kernel : python, bash, ai-agent, ...)

- easy container and kernel management

## Getting started

Compile the project :

```bash
go build -o kbox ./kbox
```

Then, run :

```bash
jupyter kernelspec list
```

to see which kernels are available to be containerized, for example :

```text
python3        /Users/mgg/Library/Jupyter/kernels/python3
```

To install a containerized version of this kernel :

```bash
kbox install image_name:latest python3 kbox-py3.13
```

> with image_name:latest being any docker image installed than can start the kernel you want to install (a container with python3.13 and ipykernel for this example) --- see [python_dockerfile](./python_dockerfile/) for examples.

Then, `jupyter kernelspec list` should output :

```bash
kbox-py3.13    /Users/mgg/Library/Jupyter/kernels/kbox-py3.13
python3        /Users/mgg/Library/Jupyter/kernels/python3
```

Any use of the kbox kernel should now launch an instance of image_name:latest, start a ipython3 kernel inside it, and communicate with the notebook (or other frontend) that started the kernel as usual. Closing the notebook will stop and remove the container.

To test :

```bash
jupyter console --kernel kbox-py3.13
```


> The directory where the kbox kernel is started is mounted on the container, allowing to access current datasets, ...

> Note : kbox only need the kernelspec file, but the kernel itself does not have to be startable from the host machine.

Example of kernels :

```text
  ir             /Users/mgg/Library/Jupyter/kernels/ir
  kbox-bash      /Users/mgg/Library/Jupyter/kernels/kbox-bash
  kbox-py3.11    /Users/mgg/Library/Jupyter/kernels/kbox-py3.11
  kbox-py3.12    /Users/mgg/Library/Jupyter/kernels/kbox-py3.12
  kbox-py3.13    /Users/mgg/Library/Jupyter/kernels/kbox-py3.13
  python3        /Users/mgg/Library/Jupyter/kernels/python3
```