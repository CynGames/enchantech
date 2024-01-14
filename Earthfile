VERSION 0.7
FROM ubuntu:focal

terraform:
    ARG --required DO_TOKEN
    RUN apt-get update && apt-get install -y gnupg software-properties-common wget
    RUN wget -O- https://apt.releases.hashicorp.com/gpg | \
        gpg --dearmor | \
        tee /usr/share/keyrings/hashicorp-archive-keyring.gpg
    RUN gpg --no-default-keyring \
        --keyring /usr/share/keyrings/hashicorp-archive-keyring.gpg \
        --fingerprint
    RUN echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] \
        https://apt.releases.hashicorp.com $(lsb_release -cs) main" | \
        tee /etc/apt/sources.list.d/hashicorp.list
    RUN apt update
    RUN apt-get install -y terraform
    RUN cd ~ && wget https://github.com/digitalocean/doctl/releases/download/v1.101.0/doctl-1.101.0-linux-amd64.tar.gz
    RUN cd ~ && tar xf ~/doctl-1.101.0-linux-amd64.tar.gz
    RUN mv ~/doctl /usr/local/bin

    WORKDIR /app
    COPY . .
    WORKDIR /app/infra
    RUN terraform init

apply:
    ARG --required DO_TOKEN
    FROM +terraform
    WITH DOCKER
        RUN terraform apply -auto-approve -var="do_token=$DO_TOKEN"
    END

destroy:
    ARG --required DO_TOKEN
    FROM +terraform
    RUN terraform destroy -auto-approve -var="do_token=$DO_TOKEN"