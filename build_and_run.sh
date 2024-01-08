set -e

npx tailwindcss -i ./src/styles/index.css -o ./static/main.css

templ generate

