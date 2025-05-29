title:	cloudflare_ruleset shows constant changes despite no fields being updated
state:	OPEN
author:	paul-zah-iptiq
labels:	kind/bug, version/5
comments:	2
assignees:	
projects:	
milestone:	
number:	5600
--
### Confirmation

- [x] This is a bug with an existing resource and is not a feature request or enhancement. Feature requests should be submitted with Cloudflare Support or your account team.
- [x] I have searched the issue tracker and my issue isn't already found.
- [x] I have replicated my issue using the latest version of the provider and it is still present.

### Terraform and Cloudflare provider version

Terraform v1.10.4
on darwin_arm64
+ provider registry.terraform.io/cloudflare/cloudflare v5.4.0
+ provider registry.terraform.io/hashicorp/aws v5.97.0
+ provider registry.terraform.io/hashicorp/cloudinit v2.3.7
+ provider registry.terraform.io/hashicorp/http v3.5.0
+ provider registry.terraform.io/hashicorp/kubernetes v2.35.1
+ provider registry.terraform.io/hashicorp/random v3.7.2
+ provider registry.terraform.io/hashicorp/tls v4.1.0

### Affected resource(s)

cloudflare_ruleset

### Terraform configuration files

```hcl
resource "cloudflare_ruleset" "waf_ruleset_customisation" {
  zone_id     = "redacted"
  name        = "default"
  kind        = "zone"
  phase       = "http_request_firewall_managed"
  description = "Created by the Cloudflare security team, this ruleset is designed to provide fast and effective protection for all your applications. It is frequently updated to cover new vulnerabilities and reduce false positives."

  rules = [{
      expression        = "true"
      action            = "execute"
      action_parameters = {
        id        = "4814384a9e5d4991b9815dcfc25d2f1f" # OWASP Rules
        ref       = "4814384a9e5d4991b9815dcfc25d2f1f"
        overrides = {
          categories = [
            {
              category = "paranoia-level-2"
              enabled  = false
            },
            {
              category = "paranoia-level-3"
              enabled  = false
            },
            {
              category = "paranoia-level-4"
              enabled  = false
            }
          ]
        }
      }
    }
    ,{
      expression        = "true"
      action            = "execute"
      action_parameters = {
        id        = "efb7b8c949ac4650a09736fc376e9aee" # Cloudflare Managed Ruleset
        ref       = "efb7b8c949ac4650a09736fc376e9aee"
        # Multiple rules disabled as not relevant (specific to wordpress, drupal, php, etc)
        overrides = {
          categories = [
            {
              category = "wordpress"
              enabled  = false
            },
            {
              category = "drupal"
              enabled  = false
            },
            {
              category = "joomla"
              enabled  = false
            },
            {
              category = "ivanti"
              enabled  = false
            },
            {
              category = "microsoft-asp-net"
              enabled  = false
            },
            {
              category = "microsoft-exchange"
              enabled  = false
            },
            {
              category = "microsoft-iis"
              enabled  = false
            },
            {
              category = "default-windows-user"
              enabled  = false
            },
            {
              category = "microsoft-sharepoint"
              enabled  = false
            },
            {
              category = "microsoft-sql-server"
              enabled  = false
            },
            {
              category = "servicenow"
              enabled  = false
            },
            {
              category = "oracle-weblogic"
              enabled  = false
            },
            {
              category = "php"
              enabled  = false
            },
            {
              category = "phpcms"
              enabled  = false
            },
            {
              category = "phpmailer"
              enabled  = false
            },
            {
              category = "php-cgi"
              enabled  = false
            },
            {
              category = "ruby"
              enabled  = false
            },
            {
              category = "ruby-on-rails"
              enabled  = false
            },
            {
              category = "adobe-coldfusion"
              enabled  = false
            },
            {
              category = "adobe-flash"
              enabled  = false
            },
            {
              category = "magento"
              enabled  = false
            },
            {
              category = "plone"
              enabled  = false
            },
            {
              category = "jenkins"
              enabled  = false
            },
            {
              category = "practico-cms"
              enabled  = false
            },
            {
              category = "vmware-vcenter"
              enabled  = false
            },
            {
              category = "citrix-netscaler-adc"
              enabled  = false
            },
            {
              category = "jetbrains-teamcity"
              enabled  = false
            },
            {
              category = "progress-ws-ftp"
              enabled  = false
            },
            {
              category = "automation-anywhere"
              enabled  = false
            },
            {
              category = "solarwinds"
              enabled  = false
            },
            {
              category = "palo-alto-networks"
              enabled  = false
            },
            {
              category = "fortios"
              enabled  = false
            },
            {
              category = "fortinet-fortimanager"
              enabled  = false
            },
            {
              category = "progress-software-whatsup-gold"
              enabled  = false
            },
            {
              category = "sonicwall-sslvpn"
              enabled  = false
            }
          ]
        }
      }
    }
  ]
  lifecycle {
    ignore_changes = [ # id's and ref's (same value) change upon replacement of rules, but provider bug on continuous changes
      id,
      rules[0].id,
      rules[0].ref,
      rules[1].id,
      rules[1].ref,
    ]
  }
}
```

### Link to debug output

https://gist.github.com/paul-zah-iptiq/a3bf914c3b967c3ab04e5cf912875e39

### Panic output

_No response_

### Expected output

Expected no changes to be detected.

### Actual output

```
Terraform will perform the following actions:

  # cloudflare_ruleset.waf_ruleset_customisation will be updated in-place
  ~ resource "cloudflare_ruleset" "waf_ruleset_customisation" {
        id          = "fe59571x365d422basd3563a10a1baea"
        name    = "default"
        # (5 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```
Note the presence of a resource update despite the lack of updates to any field.


### Steps to reproduce

1. Define a customisation ruleset entrypoint for the Managed Ruleset
2. Import the existing ruleset entrypoint
3. Apply customisation configuration
4. Apply customisation configuration - change detected despite no fields being updated

### Additional factoids

_No response_

### References

_No response_
