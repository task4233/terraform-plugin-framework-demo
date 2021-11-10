terraform {
  required_providers {
    log = {
      version = "0.1.0"
      source  = "terraform.local/local/log"
    }
  }
}

provider "log" {
  host = "http://localhost:19090"
}

resource "log_order" "item" {
  items = [
    {
      log = {
        body = "hoge"
      }
    },
    {
      log = {
        body = "piyo"
      }
    },
  ]
}

