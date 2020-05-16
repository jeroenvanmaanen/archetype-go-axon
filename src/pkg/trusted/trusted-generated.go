package trusted

func Init() {
	trustedKeys = map[string]string{}
	keyManagers = map[string]string{}
	trustedKeys["jeroen@aenea.entreact.com"] = "AAAAB3NzaC1yc2EAAAABIwAAAQEAzu3J6nPQjN71F7rkvuBoy3DuoRK144z9CrpuNuU9U86rHl33mTSCiOaWFXvArR5nUpG8Oe1qRzGnHqczLP74L8CGXmq9rmh3zXGS8goudPx9iAc1dpZSGumnffY1/o/PKKU6mEudY/KIP4ZRxZZ8l4moUCH9xwip+YIEHiUm0XGVJLoBUc8Gx/v1nzZGdKgbCMBx78SizF6rIN77pcHqCiFa5j7p7QcGwa7pPmZw7Mwuqnu7/qpRdyqmnu1q4h+f+UjsReEUH5MEWPCzhxCLOy3iN7qunWavxNjWNHMa6/JjAvyilO6FaHYkcn5uQCvM+wleMUtXuLdNx/gpVUGsHQ=="
	keyManagers["jeroen@aenea.entreact.com"] = trustedKeys["jeroen@aenea.entreact.com"]
}
