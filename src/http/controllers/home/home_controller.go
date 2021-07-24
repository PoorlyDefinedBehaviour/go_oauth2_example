package homecontroller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HomePage(ctx echo.Context) error {
	const html = `
		<html>
			<body>
				<a href="/auth/signin">Sign in with spotify</a>
			</body>
		</html>
	`
	return ctx.HTML(http.StatusOK, html)
}
