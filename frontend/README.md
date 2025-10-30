# SmartResume Frontend

Modern Next.js 14 frontend for the SmartResume AI-powered resume analysis platform.

## Features

- 🎨 **Modern UI**: Blue/Purple gradient theme with dark mode support
- 🔐 **Authentication**: JWT-based login and registration
- 📤 **Resume Upload**: Drag-and-drop PDF upload with optional job description
- 📊 **ATS Analysis**: Visual score gauge with detailed metrics
- 🎯 **Skill Analysis**: Color-coded matching and missing skills
- 💼 **Job Recommendations**: Real-time job listings from multiple APIs
- 🌙 **Dark Mode**: Toggle between light and dark themes
- 📱 **Responsive**: Mobile-first design with Tailwind CSS

## Tech Stack

- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: React Hooks
- **API Client**: Fetch API
- **Authentication**: JWT (localStorage)

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Backend server running on `http://localhost:8080`

### Installation

1. Install dependencies:
```bash
npm install
```

2. Configure environment variables:
Create `.env.local` file with:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

3. Run the development server:
```bash
npm run dev
```

4. Open [http://localhost:3000](http://localhost:3000)

## Project Structure

```
frontend/
├── src/
│   ├── app/
│   │   ├── dashboard/
│   │   │   └── page.tsx         # Main dashboard with upload
│   │   ├── login/
│   │   │   └── page.tsx         # Login page
│   │   ├── register/
│   │   │   └── page.tsx         # Registration page
│   │   ├── globals.css          # Global styles with Tailwind
│   │   ├── layout.tsx           # Root layout
│   │   └── page.tsx             # Landing page
│   ├── components/
│   │   ├── Header.tsx           # Navigation header
│   │   ├── JobCard.tsx          # Job recommendation card
│   │   ├── ScoreGauge.tsx       # ATS score circular gauge
│   │   └── SkillBadge.tsx       # Skill tag component
│   └── lib/
│       ├── api.ts               # API client and types
│       └── utils.ts             # Helper functions
├── public/                      # Static assets
├── tailwind.config.js          # Tailwind configuration
└── package.json
```

## API Integration

The frontend integrates with the Go backend API:

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/auth/profile` - Get user profile

### Resume Operations
- `POST /api/resume/upload` - Upload and analyze resume
- `GET /api/resume/all` - Get all user resumes
- `GET /api/resume/:id` - Get specific resume
- `GET /api/resume/:id/jobs` - Get job recommendations

## Color Scheme

### Light Mode
- Primary: Blue (#0ea5e9)
- Secondary: Purple (#a855f7)
- Background: White
- Text: Slate gray

### Dark Mode
- Primary: Light blue (#38bdf8)
- Secondary: Light purple (#c084fc)
- Background: Slate 900
- Text: Slate 100

## Key Components

### ScoreGauge
Circular progress indicator showing ATS score (0-100) with color-coded labels:
- **80-100**: Excellent (Green)
- **60-79**: Good (Yellow)
- **0-59**: Needs Improvement (Red)

### SkillBadge
Color-coded skill tags:
- **Green**: Matching skills
- **Red**: Missing skills

### JobCard
Job recommendation card displaying:
- Job title and company
- Location and salary
- Job type badge
- Description preview
- Apply button linking to external URL

## Development

### Build for Production
```bash
npm run build
```

### Run Production Build
```bash
npm start
```

### Linting
```bash
npm run lint
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API base URL | `http://localhost:8080` |

## Features in Detail

### Dashboard
- **Resume Upload**: PDF file picker with validation
- **Job Description**: Optional textarea for targeted analysis
- **Real-time Analysis**: Loading states during processing
- **Results Display**: Animated score gauge, skills, and jobs
- **Responsive Layout**: Mobile-friendly grid system

### Authentication
- **Form Validation**: Client-side validation
- **Error Handling**: User-friendly error messages
- **Auto-redirect**: Navigate to dashboard after login
- **Token Management**: Secure JWT storage

### Dark Mode
- **Auto-detection**: Respects system preference
- **Manual Toggle**: Header button to switch themes
- **Persistent**: Saves preference to localStorage
- **Smooth Transitions**: CSS transitions for theme changes

## Troubleshooting

### API Connection Issues
- Verify backend is running on port 8080
- Check CORS configuration in backend
- Ensure `.env.local` has correct API URL

### Build Errors
- Clear `.next` folder: `rm -rf .next`
- Delete `node_modules` and reinstall: `rm -rf node_modules && npm install`
- Check TypeScript errors: `npm run lint`

### Dark Mode Not Working
- Check browser localStorage
- Clear cache and reload page
- Verify Tailwind config has `darkMode: 'class'`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

MIT License - feel free to use this project for your own purposes.
