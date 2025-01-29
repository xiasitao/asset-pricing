terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.16.0"
    }
  }
  backend "azurerm" {
    resource_group_name  = "tfstate"
    storage_account_name = "tfstatexiasitao"
    container_name       = "tfstate"
    key                  = "asset-pricing-dev-containerapps.tfstate"
  }
}

provider "azurerm" {
  features {
  }
}

resource "azurerm_resource_group" "resource_group" {
  name     = "asset-pricing-dev-containerapps"
  location = "germanywestcentral"
}

resource "azurerm_container_app_environment" "container_app_environment" {
  name = "asset-pricing-api-dev"
  location = azurerm_resource_group.resource_group.location
  resource_group_name = azurerm_resource_group.resource_group.name
}

resource "azurerm_container_app" "container_app" {
  name = "asset-pricing-api-dev"
  resource_group_name = azurerm_resource_group.resource_group.name
  container_app_environment_id = azurerm_container_app_environment.container_app_environment.id
  revision_mode = "Single"
  ingress {
    external_enabled = true
    traffic_weight {
      percentage = 100
      latest_revision = true
    }
    target_port = 8080
  }

  template {
    container {
      name = "asset-pricing-api"
      image = "docker.io/xiasitao/asset-pricing-api:latest"
      cpu = 0.25
      memory = "0.5Gi"
    }
  }
}

