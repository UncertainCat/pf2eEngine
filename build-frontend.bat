@echo off
echo Building the React frontend...

cd frontend
call npm install
call npm run build

echo Frontend build complete.
cd ..