package hasher

import "golang.org/x/crypto/bcrypt"

type PasswordHasher struct {
	cost int
}

func New() *PasswordHasher {
	return &PasswordHasher{cost: bcrypt.DefaultCost}
}

func (h *PasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	return string(bytes), err
}

func (h *PasswordHasher) Compare(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
