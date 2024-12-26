package real

func (r *RealBackwardsInvocation) mlchainPath(path ...string) string {
	path = append([]string{"inner", "api"}, path...)
	return r.mlchainInnerApiBaseurl.JoinPath(path...).String()
}
