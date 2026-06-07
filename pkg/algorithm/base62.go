package algorithm

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func toBase62(n int) string {
	if n == 0 {
		return "0"
	}
	res := []byte{}
	for n > 0 {
		res = append([]byte{chars[n%62]}, res...)
		n /= 62
	}
	return string(res)
}

func CreateCode(n int) string {
	code := toBase62(n)
	return code
}
