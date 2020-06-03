package crypto

import (
	sm2 "github.com/tjfoc/gmsm/sm2"
)

func generateKeyPairs() error {
	_, err := sm2.GenerateKey() // 生成密钥对
	if err != nil {
		return err
	}

	return nil
}
