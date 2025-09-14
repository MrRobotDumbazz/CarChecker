package handlers

import (
	"car-status-backend/pkg/utils"
	"io/ioutil"
	"net/http"
)

type SwaggerHandler struct {
	specPath string
}

func NewSwaggerHandler(specPath string) *SwaggerHandler {
	return &SwaggerHandler{
		specPath: specPath,
	}
}

func (h *SwaggerHandler) ServeSwaggerUI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	swaggerHTML := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Car Status Detection API - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui.css" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin:0;
            background: #fafafa;
        }
        .swagger-ui .topbar {
            background-color: #2d5aa0;
        }
        .swagger-ui .topbar .download-url-wrapper {
            display: none;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: '/api/docs/openapi.yaml',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                tryItOutEnabled: true,
                requestInterceptor: (request) => {
                    console.log('Request:', request);
                    return request;
                },
                responseInterceptor: (response) => {
                    console.log('Response:', response);
                    return response;
                }
            });
        }
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(swaggerHTML))
}

func (h *SwaggerHandler) ServeOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	specContent, err := ioutil.ReadFile(h.specPath)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to read OpenAPI specification")
		return
	}

	w.Header().Set("Content-Type", "application/yaml")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(specContent)
}

func (h *SwaggerHandler) ServeSwaggerRedoc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	redocHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Car Status Detection API - ReDoc</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
    <style>
        body {
            margin: 0;
            padding: 0;
        }
    </style>
</head>
<body>
    <div id="redoc-container"></div>
    <script src="https://cdn.jsdelivr.net/npm/redoc@2.1.3/bundles/redoc.standalone.js"></script>
    <script>
        Redoc.init('/api/docs/openapi.yaml', {
            scrollYOffset: 50,
            hideHostname: false,
            theme: {
                colors: {
                    primary: {
                        main: '#2d5aa0'
                    }
                },
                typography: {
                    fontSize: '14px',
                    fontFamily: 'Roboto, sans-serif',
                    headings: {
                        fontFamily: 'Montserrat, sans-serif'
                    }
                }
            }
        }, document.getElementById('redoc-container'));
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(redocHTML))
}

func (h *SwaggerHandler) ApiDocsIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	indexHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Car Status Detection API Documentation</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background: #f5f5f5;
        }
        .container {
            background: white;
            border-radius: 10px;
            padding: 30px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #2d5aa0;
            margin-bottom: 10px;
        }
        .subtitle {
            color: #666;
            margin-bottom: 30px;
        }
        .docs-grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 20px;
            margin-top: 30px;
        }
        .doc-card {
            border: 1px solid #ddd;
            border-radius: 8px;
            padding: 20px;
            text-decoration: none;
            color: inherit;
            transition: box-shadow 0.2s;
        }
        .doc-card:hover {
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        }
        .doc-card h3 {
            margin: 0 0 10px 0;
            color: #2d5aa0;
        }
        .doc-card p {
            margin: 0;
            color: #666;
            font-size: 14px;
        }
        .endpoints {
            margin-top: 30px;
        }
        .endpoint {
            background: #f8f9fa;
            border-left: 4px solid #2d5aa0;
            padding: 15px;
            margin-bottom: 10px;
            border-radius: 4px;
        }
        .method {
            font-weight: bold;
            color: #2d5aa0;
        }
        @media (max-width: 600px) {
            .docs-grid {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚗 Car Status Detection API</h1>
        <p class="subtitle">API для системы определения состояния автомобиля по фотографии</p>

        <div class="docs-grid">
            <a href="/api/docs/swagger" class="doc-card">
                <h3>📋 Swagger UI</h3>
                <p>Интерактивная документация API с возможностью тестирования endpoints</p>
            </a>

            <a href="/api/docs/redoc" class="doc-card">
                <h3>📖 ReDoc</h3>
                <p>Красивая и детальная документация API в современном стиле</p>
            </a>

            <a href="/api/docs/openapi.yaml" class="doc-card">
                <h3>📄 OpenAPI Spec</h3>
                <p>Raw OpenAPI 3.0 спецификация в формате YAML</p>
            </a>

            <a href="/" class="doc-card">
                <h3>🏠 API Root</h3>
                <p>Основная информация о сервисе и список endpoints</p>
            </a>
        </div>

        <div class="endpoints">
            <h3>🔗 Основные Endpoints:</h3>

            <div class="endpoint">
                <span class="method">POST</span> /api/v1/images/upload
                <p>Загрузка изображения автомобиля</p>
            </div>

            <div class="endpoint">
                <span class="method">POST</span> /api/v1/predict/{image_id}
                <p>Запуск анализа состояния автомобиля</p>
            </div>

            <div class="endpoint">
                <span class="method">GET</span> /api/v1/predictions/{id}
                <p>Получение результата анализа</p>
            </div>

            <div class="endpoint">
                <span class="method">GET</span> /api/v1/health
                <p>Проверка состояния системы</p>
            </div>
        </div>

        <div style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; color: #666; font-size: 14px;">
            <p><strong>Для ML команды:</strong> Смотрите секцию "ML Service" в Swagger UI для требований к Python API</p>
            <p><strong>Для Frontend команды:</strong> Все endpoints задокументированы с примерами запросов и ответов</p>
        </div>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(indexHTML))
}