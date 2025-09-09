# readonly-policy.hcl
path "*" {
    capabilities = ["read", "list"]
}

path "sys/mounts" {
    capabilities = ["read", "list"]
}