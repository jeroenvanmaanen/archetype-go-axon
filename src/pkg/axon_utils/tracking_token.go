package axon_utils

type TokenStore interface {
    ReadToken()*int64
    WriteToken(int64) error
}

type NullTokenStore struct {}

func (tokenStore *NullTokenStore) ReadToken() *int64 {
    return nil
}

func (tokenStore *NullTokenStore) WriteToken(int64) error {
    return nil
}
