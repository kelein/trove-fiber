package docs

import "github.com/kelein/trove-fiber/pkg/version"

// InitSwaggerInfo setup swagger info
// * Basic Info
// @license.name Apache 2.0
// @contact.name trove-gin
// @contact.url https://github.com/kelein/trove-gin
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// * Authentication Info
// @securityDefinitions.apiKey Bearer
// @name Authorization
// @in header
func InitSwaggerInfo() {
	SwaggerInfo.BasePath = "v1"
	SwaggerInfo.Title = "Trove API Server"
	SwaggerInfo.Version = version.AppVersion
	SwaggerInfo.Description = version.String()
	SwaggerInfo.InfoInstanceName = "swagger"
	SwaggerInfo.Schemes = []string{"http", "https"}
}
