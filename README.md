# KBox

KBox enables you to run Jupyter kernels inside Docker containers while maintaining a seamless connection via your preferred Jupyter frontend (e.g., Jupyter Notebooks, JupyterLab, or VS Code).

By wrapping the kernel execution in a container, KBox provides:
- **Isolated Execution**: Run any kernel (Python, Bash, AI agents, etc.) in a fully containerized environment.
- **Simplified Management**: Easily switch between different environment images and kernel versions without polluting your host machine.
- **Automatic Lifecycle**: Containers are automatically started when the kernel is launched and removed when the notebook is closed.

## Getting Started

### 1. Installation
Compile the project to create the `kbox` binary:

```bash
go build -o kbox ./kbox
```

### 2. Identify Target Kernels
List the kernels available on your host machine to see which ones you want to containerize:

```bash
jupyter kernelspec list
```

Example output:
```text
python3        /Users/mgg/Library/Jupyter/kernels/python3
```

### 3. Install a Containerized Kernel
Use `kbox` to create a new kernelspec that points to a Docker image.

```bash
kbox install <image_name:tag> <original_kernel_name> <new_kernel_name>
```

**Example:**
To create a containerized Python 3.13 kernel using a specific image:
```bash
kbox install my-python-env:latest python3 kbox-py3.13
```
*Note: The Docker image must contain the necessary runtime and `ipykernel` to support the kernel. See [python_dockerfile](./python_dockerfile/) for examples.*

### 4. Verification
Verify that the new kernel is installed:
```bash
jupyter kernelspec list
```

You should now see your new kernel in the list:
```text
kbox-py3.13    /Users/mgg/Library/Jupyter/kernels/kbox-py3.13
python3         /Users/mgg/Library/Jupyter/kernels/python3
```

## Usage

Once installed, simply select the `kbox-` kernel from your Jupyter dropdown. KBox will launch an instance of the specified Docker image and bridge the communication between the notebook and the container.

**Key Features:**
- **Working Directory**: The directory where the KBox kernel is started is automatically mounted into the container, providing easy access to local datasets and files.
- **Host Independence**: The original kernel only needs to exist as a `kernelspec` on the host; it does not actually need to be executable on the host machine itself.

### Quick Test
You can test the installation via the terminal:
```bash
jupyter console --kernel kbox-py3.13
```
