package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	port     = ":8088"
	yamlPath = "../../api/contract/v1/workmate.yaml"
)

// HTML шаблон для Swagger UI
const swaggerHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>WorkMate API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui.css" />
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
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '/swagger.yaml',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>
`

func main() {
	log.Printf("WorkMate Swagger UI Server starting...")

	// Читаем YAML файл
	yamlData, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		log.Printf("Warning: Could not read %s: %v", yamlPath, err)
		log.Printf("Trying alternative paths...")

		// Пробуем альтернативные пути
		altPaths := []string{
			"api/contract/v1/workmate.yaml",
			"./api/contract/v1/workmate.yaml",
			"../api/contract/v1/workmate.yaml",
		}

		for _, path := range altPaths {
			yamlData, err = ioutil.ReadFile(path)
			if err == nil {
				log.Printf("Found YAML at: %s", path)
				break
			}
		}

		if err != nil {
			log.Fatalf("Error: Could not read workmate.yaml from any location: %v", err)
		}
	}

	// Обработчик для главной страницы Swagger UI
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl, err := template.New("swagger").Parse(swaggerHTML)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// Обработчик для YAML файла
	http.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Write(yamlData)
	})

	// Обработчик для проверки здоровья
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "service": "swagger-ui"}`))
	})

	server := &http.Server{
		Addr:    port,
		Handler: nil, // используем DefaultServeMux
	}

	log.Printf("Swagger UI доступен по адресу: http://localhost%s/", port)
	log.Printf("YAML спецификация доступна по адресу: http://localhost%s/swagger.yaml", port)
	log.Printf("Health check: http://localhost%s/health", port)

	log.Fatal(server.ListenAndServe())
}
