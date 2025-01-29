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
    key                  = "asset-pricing-dev-appservice.tfstate"
  }
}

provider "azurerm" {
  features {
  }
}

resource "azurerm_resource_group"  "resource_group" {
  name = "asset-pricing-dev-webapps"
  location = "westus2"
}

resource "azurerm_service_plan" "service_plan" {
  name = "asset-pricing-dev-service-plan"
  resource_group_name = azurerm_resource_group.resource_group.name
  location = azurerm_resource_group.resource_group.location
  os_type = "Linux"
  sku_name = "B1"
}

resource "azurerm_linux_web_app" "webapp" {
  name = "asset-pricing-dev"
  resource_group_name = azurerm_resource_group.resource_group.name
  location =  azurerm_resource_group.resource_group.location
  service_plan_id = azurerm_service_plan.service_plan.id

  site_config {
    application_stack {
      docker_registry_url = "https://index.docker.io"
      docker_image_name = "xiasitao/asset-pricing-api:latest"
    }
    always_on = true
  }

  app_settings = {
    "WEBSITES_PORT" = 8080
    "PORT" = 8080
  }

}
