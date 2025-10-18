# Nordakademie *Bazzar*

<p align="center">
  <img src="./assets/banner.png" alt="Nordakademie Basar Banner" width="100%">
</p>

> # A digital marketplace for the Nordakademie community in Hamburg

## About

Nordakademie Bazzar is a marketplace platform designed for students and staff at Nordakademie Hochschule der Wirtschaft in Hamburg. Similar to eBay or classifieds platforms, it enables the community to buy, sell, and share various offerings including:

- Course materials and textbooks
- Flat and apartment listings
- Weekend seminars and events
- General marketplace items

## Features

- **Post Listings** - Create and manage marketplace posts with images and detailed descriptions
- **Community Interaction** - Comment on posts and connect with other members
- **Image Support** - Upload multiple images with automatic optimization
- **Easy Payments** - Integrated PayPal support for transactions
- **Search & Filter** - Find listings by creator, tags, or category

## Quick Start

### Using Docker (Recommended)

```bash
docker-compose up -d
```

The application will be available at `http://localhost:42000`

### Manual Setup

1. Ensure MongoDB is running
2. Copy `.env.example` to `.env` and configure your settings
3. Run the application:
   ```bash
   go run main.go
   ```

## Health Check

Check if the service is running:

```
GET http://localhost:42000/health
```

## Documentation

For detailed API documentation, see [API_DOCS.md](./API_DOCS.md)

## About Nordakademie

Nordakademie Hochschule der Wirtschaft is a state-recognized private university located in Hamburg, Germany, specializing in business and technology education.

---

Built with ❤️ for the Nordakademie community
