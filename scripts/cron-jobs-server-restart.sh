#!/bin/bash

# Navigate to the cmd folder
cd cron-jobs

# Create the build file using go build and wait for it to complete
echo "Building the project..."
if go build; then
  echo "Build completed successfully."
else
  echo "Build failed." >&2
  exit 1
fi

# Restart the pm2 process with the new build
echo "Restarting the pm2 process..."
pm2 restart cron-jobs --update-env

if [ $? -eq 0 ]; then
  echo "pm2 process restarted successfully."
else
  echo "Failed to restart pm2 process." >&2
  exit 1
fi
