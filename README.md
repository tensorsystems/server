![Tensor EMR](./logo.png)

[![Docker Publish](https://github.com/Tensor-Systems/tensoremr-server/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/Tensor-Systems/tensoremr-server/actions)

Tensor EMR is a web based Open Source Electronic Medical Record software application. It features patient registration, patient chart management, appointment & scheduling, queue management, eRx, and more. 

### Quick Start

1. Pull repository 
2. Start the Go Server by executing `make run`
3. Explore at `http://localhost:8080/api`

#### Use Docker
1. Pull `https://github.com/Tensor-Systems/tensoremr-deploy`
2. `docker-compose up -d`