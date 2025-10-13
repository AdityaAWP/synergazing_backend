#!/bin/bash
set -e

echo "🚀 Memulai proses deployment..."

echo " menarik perubahan terbaru dari branch main..."
git pull origin main

echo " Membangun ulang dan menjalankan ulang kontainer Docker..."
sudo docker compose down
sudo docker compose up --build -d

echo "✅ Deployment selesai dengan sukses!"