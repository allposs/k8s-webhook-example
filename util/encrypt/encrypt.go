package encrypt

import "golang.org/x/crypto/bcrypt"

// Encry bcrypt.GenerateFromPassword加密封装
func Encry(key string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash), err
	//验证方式
	//bcrypt.CompareHashAndPassword(hashedPassword, passwordNotCheck)
}
