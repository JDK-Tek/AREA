# AREA Web Application

## Installation

To install the application, run the following command:

```bash
npm install
```

## Launch

To start the application, use the following command:

```bash
npm start
```

## Deployment | CI/CD

The application is deployed on a Vercel domain.
<br>
The production deployment is automatically updated when a new commit is pushed to the `main` branch.
<br>
There are also preview deployments for each commits on branches.
<br>
You can access it [here](https://area-jeepg.vercel.app/).

## Environment Variables

The application requires a `.env` file with the following environment variable:

```
REACT_APP_BACKEND_URL=<your-backend-url>
```

Make sure to replace `<your-backend-url>` with the actual URL of your backend.