#!/bin/bash

# Exibir informações de debug
echo "Current directory: $(pwd)"
echo "Listing files:"
ls -la

# Garantir que as pastas templates e static existam no diretório de build
echo "Copying templates and static directories..."
mkdir -p .vercel/output/static/templates
mkdir -p .vercel/output/static/static

# Copiar os arquivos para o diretório de saída da Vercel
cp -r templates/* .vercel/output/static/templates/
cp -r static/* .vercel/output/static/static/

echo "Build completed" 