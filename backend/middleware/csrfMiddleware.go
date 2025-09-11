package middleware

import (
    "net/http"
    "github.com/alimosavifard/zyros-backend/utils"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/csrf"
    "os"
)

// The auth key should be a 32-byte secret. It's best to set this as an environment variable.
var csrfAuthKey = []byte(os.Getenv("CSRF_AUTH_KEY"))

func CSRFMiddleware() gin.HandlerFunc {
    // Check if the key is set. If not, generate a random one for development only.
    if len(csrfAuthKey) != 32 {
        utils.InitLogger().Fatal().Msg("CSRF_AUTH_KEY environment variable is not set to a 32-byte string.")
    }

    // CSRF middleware with required settings
    csrfMiddleware := csrf.Protect(
        csrfAuthKey,
        csrf.Path("/"),
        csrf.HttpOnly(true),
        csrf.Secure(true), // Set to false in a local environment for HTTP
        csrf.SameSite(csrf.SameSiteStrictMode),
    )

    return func(c *gin.Context) {
        csrfMiddleware(c.Writer, c.Request)
        c.Next()
    }
}