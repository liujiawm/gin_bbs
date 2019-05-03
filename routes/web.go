package routes

import (
	"gin_bbs/pkg/ginutils/captcha"
	"gin_bbs/pkg/ginutils/router"
	"gin_bbs/routes/middleware"
	"gin_bbs/routes/wrapper"

	"gin_bbs/app/controllers/auth/login"
	"gin_bbs/app/controllers/auth/password"
	"gin_bbs/app/controllers/auth/register"
	"gin_bbs/app/controllers/auth/verification"
	"gin_bbs/app/controllers/page"
	"time"
)

func registerWeb(r *router.MyRoute) {
	r.Register("GET", "root", "/", page.Root)
	r.Register("GET", "captcha", "/captcha/:id", captcha.Handler) // 验证码

	// ------------------------------------- Auth -------------------------------------
	// 用户身份验证相关的路由
	r.Register("GET", "login.show", "/login", middleware.Guest(), login.ShowLoginForm)
	r.Register("POST", "login", "/login", middleware.Guest(), login.Login)
	r.Register("POST", "logout", "/logout", login.Logout)

	// 用户注册相关路由
	r.Register("GET", "register.show", "/register", middleware.Guest(), register.ShowRegistrationForm)
	r.Register("POST", "register", "/register", middleware.Guest(), register.Register)

	// 密码重置相关路由
	pwdRouter := r.Group("/password", middleware.Guest())
	{
		pwdRouter.Register("GET", "password.request", "/reset", password.ShowLinkRequestForm)
		pwdRouter.Register("POST", "password.email", "/email", password.SendResetLinkEmail)
		pwdRouter.Register("GET", "password.reset", "/reset/:token", password.ShowResetForm)
		pwdRouter.Register("POST", "password.update", "/reset", password.Reset)
	}

	// Email 认证相关路由
	verificationRouter := r.Group("/email", middleware.Auth())
	{
		verificationRouter.Register("GET", "verification.notice", "/verify", wrapper.GetUser(verification.Show))
		verificationRouter.Register("GET", "verification.verify", "/verify/:token",
			middleware.RateLimiter(1*time.Minute, 6), // 1 分钟最多 6 次请求
			verification.Verify)
		verificationRouter.Register("GET", "verification.resend", "/resend",
			middleware.RateLimiter(1*time.Minute, 6),
			wrapper.GetUser(verification.Resend))
	}
}
