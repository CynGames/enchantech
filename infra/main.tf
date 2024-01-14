terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

variable "image_repo" {
  type    = string
  default = "enchantech-codex"
}

variable "do_token" {
  type = string
}

provider "digitalocean" {
  token = var.do_token
}

resource "digitalocean_container_registry" "enchantech-codex-registry" {
  name                   = "enchantech-codex-registry"
  subscription_tier_slug = "starter"
}

resource "null_resource" "enchantech-codex-packaging" {
  provisioner "local-exec" {
    command = <<EOF
        doctl auth init -t "${var.do_token}"
        doctl registry login
        docker build --platform="linux/amd64" -t "${digitalocean_container_registry.enchantech-codex-registry.endpoint}/${var.image_repo}:latest" -f ../Dockerfile.prod ..
        docker push "${digitalocean_container_registry.enchantech-codex-registry.endpoint}/${var.image_repo}:latest"
        echo "FINISHED"
	    EOF
  }

  triggers = {
    "run_at" = timestamp()
  }

  depends_on = [
    digitalocean_container_registry.enchantech-codex-registry,
  ]
}

resource "digitalocean_app" "enchantech-codex-app" {
  spec {
    name   = "enchantech-codex-app"
    region = "nyc3"

    service {
      name               = "api"
      environment_slug   = "go"
      instance_count     = 1
      instance_size_slug = "basic-xxs"
      http_port          = 11001

      image {
        registry_type = "DOCR"
        repository    = var.image_repo
        tag           = "latest"
        deploy_on_push {
          enabled = true
        }
      }

      env {
        key   = "DATABASE_URI"
        value = data.digitalocean_database_cluster.cluster-data.uri
      }

      run_command = "bash ../run_seed.sh"
    }

    database {
      cluster_name = "enchantech-codex-cluster-tf"
      name         = "enchantech-codex-db-tf"
      engine       = "MYSQL"
      production   = true
    }
  }

  depends_on = [
    null_resource.enchantech-codex-packaging,
  ]
}

resource "digitalocean_project" "enchantech-codex-project" {
  name      = "enchantech-codex"
  resources = [
    digitalocean_app.enchantech-codex-app.urn,
    digitalocean_database_cluster.enchantech-codex-cluster.urn,
  ]
}

resource "digitalocean_database_cluster" "enchantech-codex-cluster" {
  name       = "enchantech-codex-cluster-tf"
  engine     = "mysql"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

## Outputs

data "digitalocean_database_cluster" "cluster-data" {
  name = digitalocean_database_cluster.enchantech-codex-cluster.name
}

output "database_output" {
  value = data.digitalocean_database_cluster.cluster-data.uri
  sensitive = true
}

## Template

data "template_file" "env" {
  template = file("${path.module}/env.tpl")

  vars = {
    database_uri = data.digitalocean_database_cluster.cluster-data.uri
  }
}

resource "local_file" "env_file" {
  filename = "${path.module}/.env"
  content  = data.template_file.env.rendered
}
