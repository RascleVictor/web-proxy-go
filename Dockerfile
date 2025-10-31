# Étape 1 : build (compilation)
FROM golang:1.24-alpine AS builder

# Variables d'environnement
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Installation de git et ca-certificates (utile pour les modules privés)
RUN apk add --no-cache git ca-certificates

# Définir le répertoire de travail
WORKDIR /app


# Copier les fichiers de dépendances Go
COPY go.mod go.sum ./

# Télécharger les dépendances
RUN go mod download

# Copier le code source
COPY . .

# Compiler le binaire
RUN go build -o /proxy ./cmd/proxy

# Étape 2 : runtime (exécution)
FROM alpine:3.20

# Ajouter les certificats SSL
RUN apk add --no-cache ca-certificates

# Créer un utilisateur non-root pour la sécurité
RUN adduser -D -g '' proxyuser

WORKDIR /app
COPY config/config.yml /app/config/config


# Copier le binaire depuis le builder
COPY --from=builder /proxy /app/proxy

# Donner les droits à l'utilisateur
RUN chown -R proxyuser:proxyuser /app
USER proxyuser

# Exposer le port HTTP (à adapter si nécessaire)
EXPOSE 8080

# Démarrage du service
ENTRYPOINT ["/app/proxy"]
