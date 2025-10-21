# SocioVault Landing Page & Server

Standalone Go server for the SocioVault community banking platform.

## Features

- **Contribution Management**: Accept and track community contributions
- **Help Requests**: Manage requests for support from community members
- **Live Statistics**: Real-time dashboard stats (balance, distribution %, stories funded)
- **Email Subscriptions**: Collect emails for updates
- **Voting System**: Allow community to vote on help requests
- **Data Persistence**: All data saved to JSON files in `./data/` directory

## Running the Server

### Option 1: Direct Go Command
```bash
cd c:\Users\aship\Desktop\Digital Orael\signal-bank-landing
go run main.go
```

The server will start on port 8081 by default.

### Option 2: Compiled Binary
```bash
go build -o sociovault.exe main.go
./sociovault.exe
```

### Option 3: Custom Port
```bash
set PORT=8080
go run main.go
```

## API Endpoints

### Contributions
- `POST /api/contribute` - Submit a contribution
- `GET /api/contribute` - Get all contributions

### Help Requests
- `POST /api/requests` - Submit a help request
- `GET /api/requests` - Get all help requests

### Statistics
- `GET /api/stats` - Get real-time statistics

### Voting
- `POST /api/vote` - Vote on a help request

### Subscriptions
- `POST /api/subscribe` - Subscribe to updates

## Data Storage

All data is stored in the `./data/` directory:
- `contributors.json` - List of all contributions
- `requests.json` - List of all help requests
- `subscribers.json` - Email subscribers

## Frontend

The landing page (`index.html`) includes:
- Interactive forms for contributions and help requests
- Real-time stats display
- Email subscription
- Beautiful responsive design

JavaScript (`script.js`) handles:
- Form submissions to API endpoints
- Live statistics updates
- CORS-enabled requests to the API

## Accessing the Page

Open your browser and go to:
```
http://localhost:8081
```

## Features in Action

1. **View Stats**: Stats auto-load from API on page load
2. **Make a Contribution**: Fill form and submit to contribute to the vault
3. **Request Help**: Submit your story and request amount
4. **Subscribe**: Get updates via email
5. **Vote**: Click vote button on requests to support them

## Development

To modify the server:
1. Edit `main.go` for API logic
2. Edit `script.js` for frontend interactions
3. Edit `index.html` and `styles.css` for layout/styling
4. Restart the server to see changes

## Notes

- CORS is enabled for local development
- All contributions require a valid email
- Data persists between server restarts
- Stats include both stored data and baseline values
