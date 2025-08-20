# Reviews Browser
Reviews Browser is a web application that allows users to browse and search through recent App Store reviews for iOS apps. The application consists of a frontend built with React and a backend API built with Go.

# Possible Improvements

## Backend
- Use a real database
- Integrate a CI workflow
- Separate out schedulers to cron jobs with their own deployments
- Add more unit tests, using properly mocked dependencies
- Integrate a logging and metrics solution like Datadog
- Implement a circuit breaker for the App Store API client
- For a multi-deployment environment, consider use of messaging queues (Kafka, SQS, etc) to address race conditions

## Frontend
- Implement pagination for reviews
- Fetch available app IDs from the backend and display them in a drop-down
- Allow adding a new App ID to the system
- Show more App related details, eg, app name, icon, etc.
- Add tests
