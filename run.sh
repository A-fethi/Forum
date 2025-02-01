#!/bin/bash

# Arrêter et supprimer tous les conteneurs en cours d'exécution
docker stop $(docker ps -aq) 2>/dev/null && docker rm $(docker ps -aq) 2>/dev/null

# Supprimer toutes les images Docker
docker rmi -f $(docker images -q) 2>/dev/null

# Construire l'image
docker build -t my_go_image .

# Exécuter le conteneur
docker run -d --name my_go_app -p 8080:8080 my_go_image
