package constants

type OTPType string

const (
	OTPUserRegister   OTPType = "USER_REGISTER"
	OTPForgetPassword OTPType = "FORGET_PASSWORD"
	OTPVerifyEmail    OTPType = "VERIFY_EMAIL"
)
