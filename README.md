# SmartBudget

An intelligent budgeting app that uses AI to provide personalized financial insights and recommendations.

## Features

- Real-time spending analysis using Hugging Face DistilBERT
- Adaptive budgeting based on spending patterns
- Micro-savings challenges
- Voice-based transaction logging
- Hyper-local price insights
- Bill predictions
- Mood-driven spending analysis
- Social savings goals
- Crypto tips
- Eco-conscious spending tracking

## Tech Stack

### Frontend
- Swift/SwiftUI
- MVVM Architecture
- Combine for reactive programming
- UserDefaults for local storage

### Backend
- Go with Gin framework
- PostgreSQL on Supabase
- Hugging Face DistilBERT for AI
- SendGrid for email notifications

### Infrastructure
- Fly.io for backend hosting
- Supabase for database
- Various free APIs (Plaid sandbox, Numbeo, OpenWeatherMap, CoinGecko)

## Project Structure

```
smartbudget/
├── ios/                 # iOS app (Swift/SwiftUI)
│   ├── Models/         # Data models
│   ├── Views/          # SwiftUI views
│   ├── ViewModels/     # MVVM view models
│   └── Services/       # API services
│
├── backend/            # Go backend
│   ├── api/           # API handlers
│   ├── config/        # Configuration
│   ├── db/            # Database models
│   ├── services/      # Business logic
│   └── utils/         # Utility functions
│
└── docs/              # Documentation
```

## Setup Instructions

### Prerequisites
- Xcode 14+
- Go 1.20+
- PostgreSQL
- Supabase account
- Fly.io account
- SendGrid account

### Backend Setup
1. Navigate to backend directory
2. Run `go mod tidy`
3. Set up environment variables
4. Run `go run main.go`

### iOS Setup
1. Open `ios/SmartBudget.xcodeproj`
2. Install dependencies
3. Run in simulator

## Development

### Running Tests
- Backend: `go test ./...`
- iOS: Use Xcode test navigator

### Deployment
- Backend: `fly deploy`
- iOS: Build and run in simulator

## License
MIT License

## Contributing
1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request 