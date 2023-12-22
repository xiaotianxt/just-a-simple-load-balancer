# Just A Simple Load Balancer

This repository contains the source code for a simple load balancer, written in Go.

## Features

- Supports both `amd64` and `arm64` architectures.
- Dockerized for easy deployment.
- Automated Docker image build and push using GitHub Actions.

## Usage

To build and run the Docker image, use the following command:

```sh
docker run -p 8088:8088 xiaotianxt/just-a-simple-load-balancer:latest
```

## Acknowledgements
99% of the code in this repository was written by ChatGPT, a language model developed by OpenAI.

## License
This project is licensed under the terms of the MIT license. See the LICENSE file for details.