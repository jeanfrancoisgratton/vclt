package env

import (
	cerr "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
)

func CreateEnvFile(fname string) *cerr.CustomError {
	if EnvName == "" || VAddress == "" || VUserName == "" || VPassword == "" || KVstorePath == "" {
		return &cerr.CustomError{Fatality: cerr.Warning, Title: "Missing parameters", Message: "Use : vclt env create -h for more info", Code: 2}
	}

	es := Config_s{
		EnvironmentName: EnvName,
		VaultAddress:    VAddress,
		VaultUsername:   VUserName,
		VaultPassword:   hf.EncodeString(VPassword, ""),
		KVEnginePath:    KVstorePath,
		Comments:        EnvComments,
	}

	return es.SaveEnvironmentFile(fname)
}
