package trusted

func Init() {
    trustedKeys = map[string]string{}
    keyManagers = map[string]string{}
    trustedKeys["demokey"] = "AAAAB3NzaC1yc2EAAAADAQABAAABAQCsxLekrKxHNTzfH0Qzeq9VUnScK+hCpC97bJPJGDgHvynYxy/x7mvZPF/p6X5lvSs6HA/tVsaVCnztdmE5sYQ/RgLerdlIvPs3o3HCcEcVr/YKGnNMXC923Gs2cKkcbhDHqIcimgQ+yQdf+tgyw/xtK7WgInxgh7rxJTRqmhQP0LkhzIGNHXW//JF1f3R/i+wpQ6X/cAXiLi26ZIcVPEqU3o8RsyAnYJG3nQbhH0I+5Snl38bOER3seG3H9zbOQLk+I0YSS/A60ko4gxBl8gN0/6bzCyYRfh6EkAmlilYnTggDhvmF8WuuYrKtndlpwYpygO4Qp6Ez4yf8dDrRPdkt"
    keyManagers["demokey"] = trustedKeys["demokey"]
}
