# nev-home Terraform

This repository contains Terraform configuration to manage a local RouterOS (MikroTik) router for a home setup.

## Prerequisites

- Terraform (recommended >= 1.0). Verify with `terraform version`.
- Network access to the RouterOS device (the configuration uses `https://nev-core-01.local` by default). Ensure DNS or /etc/hosts resolves the name, or update the provider `hosturl` in `main.tf` to an IP or hostname you can reach.
- RouterOS credentials: a user with sufficient privileges to manage interfaces, IP addresses, DHCP and certificates. The configuration defaults to username `tfuser` in `main.tf` and reads the password from the Terraform variable `tfuserpass`.

## terraform.tfvars file

An example file is included at `terraform.tfvars.example`. Copy it to `terraform.tfvars` and update the values before running Terraform:

1. Copy the example:

   cp terraform.tfvars.example terraform.tfvars

2. Edit `terraform.tfvars` and set:
   - `tfuserpass` — password for the RouterOS account (sensitive). Do NOT commit this file if it contains real secrets.
   - `ipv6_prefix` — the IPv6 prefix to use for the RouterOS IPv6 pool (example: `fd00:dead:beef::` or `fd00:dead:beef::/48` depending on your addressing plan).
   - `ipv6_bridge_ip` — the IPv6 address (with prefix length) to assign to the bridge interface (example: `fd00:dead:beef::1/64`).

You can also supply the password and variables via environment variables instead of a tfvars file. Terraform supports `TF_VAR_<name>` for variables (e.g., `TF_VAR_tfuserpass`). The provider block in `main.tf` includes comments indicating support for environment variables like `ROS_HOSTURL`, `ROS_USERNAME`, and `ROS_PASSWORD` — adjust your environment or `main.tf` if you prefer that method.

## Running Terraform

From the repository root:

1. Initialize the working directory and providers:

   terraform init

2. Review the planned changes:

   terraform plan -out=tfplan

3. Apply the changes:

   terraform apply tfplan

Or run `terraform apply` and confirm when prompted.

## Security notes

- Do NOT commit real credentials. Keep `terraform.tfvars` out of version control (it is common to add it to `.gitignore`).
- The provider is configured with `insecure = true` in `main.tf` to skip TLS verification — this is convenient for local self-signed certs but is insecure for production. If you have proper certificates, set `insecure = false`.

## Contact / Troubleshooting

- If Terraform cannot reach the device, check network connectivity and DNS. Try `curl -vk https://nev-core-01.local` from the machine running Terraform to verify TLS/connectivity.
- If authentication fails, verify username and password and that the `tfuser` account has the required privileges.

