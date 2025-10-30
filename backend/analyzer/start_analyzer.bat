@echo off
cd /d "%~dp0"
echo Starting Resume Analyzer Service...
python -m uvicorn app:app --host 127.0.0.1 --port 8000 --workers 1
pause
