# AI Agent Kernel

To have access to an agent in a contairized kernel, use the dockerfile here to build an image with pydantic-ai-kernel.

## 1 - Create agent configuration

Rename [ex_jupyter_pydantic_ai_config.yaml](./ex_jupyter_pydantic_ai_config.yaml) into [jupyter_pydantic_ai_config.yaml](./jupyter_pydantic_ai_config.yaml). See [https://mariusgarenaux.github.io/pydantic-ai-kernel/agent_config/](https://mariusgarenaux.github.io/pydantic-ai-kernel/agent_config/) to fill the config file.

Then, build the image :

```bash
docker build -t <image_name> .
```

> ⚠️ the config is copied inside the container, don't share it if it contains credentials ⚠️

## 2 - Install kernel as usual

```bash
./kbox install <image_name> pydantic_ai <kernel_name>
```

choose anything for kernel_name

`<image_name>` : the image you just built above


