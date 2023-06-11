terraform {
    required_providers {
        azurerm = {
            source = "hashicorp/azurerm"
            version = "~> 3.5.0"
        }
    }
    backend "azurerm" {
      
    }
}

provider "azurerm" {
    features {}
}

variable "storage_account_name" {
    type = string
    description = "Name of the storage account that is created"
}

variable "resource_group_name" {
    type = string
    description = "Name of the resource group that is created"
}

data "azurerm_resource_group" "test_example" {
    name = "rg-lab-ben.kooijman"
}

resource "azurerm_storage_account" "test_example" {
  name                     = "staterratesteuwtest01"
  resource_group_name      = data.azurerm_resource_group.test_example.name
  location                 = data.azurerm_resource_group.test_example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "test"
  }
}

resource "azurerm_storage_container" "test_example" {
  name                  = "repository"
  storage_account_name  = azurerm_storage_account.test_example.name
  container_access_type = "private"
}

resource "azurerm_storage_blob" "test_example" {
  name                   = "hello-world.txt"
  storage_account_name   = azurerm_storage_account.test_example.name
  storage_container_name = azurerm_storage_container.test_example.name
  type                   = "Block"
  source                 = "hello-world.txt"
}

output "storage_account_name" {
    value = azurerm_storage_account.test_example.name
}

output "container_name" {
    value = azurerm_storage_container.test_example.name
}
