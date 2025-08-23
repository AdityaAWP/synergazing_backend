set -e

PROJECT_DIR="/home/rebbeca/synergazing/synergazing_backend"

echo "🚀 Starting deployment for Synergize..."

cd "$PROJECT_DIR" || { echo "Error: Project directory not found. Aborting."; exit 1; }

echo "Pulling from main branch..."
git pull origin main

echo "Restart Docker container"
sudo docker compose down
sudo docker compose up --build -d

echo "Deployment finished successfully!"
