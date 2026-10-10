@echo off
:: =============================================================================
:: start-ecommerce.bat — One-shot local K8s deployment for Windows.
::
:: Works with Docker Desktop or Podman (auto-detected; Docker wins when both
:: are running). Force a runtime:
::   start-ecommerce.bat -Runtime podman
::
:: Launches k3d-setup.ps1, which installs k3d/kubectl if missing.
:: =============================================================================

chcp 65001 >nul 2>&1
title Ecommerce Microservices — K8s Setup

echo.
echo  ╔══════════════════════════════════════════════════════════╗
echo  ║     Ecommerce Microservices ^· Local K8s Setup           ║
echo  ║     Docker or Podman  ^+  k3d  ^+  NGINX Ingress           ║
echo  ╚══════════════════════════════════════════════════════════╝
echo.

:: ── Locate PowerShell (prefer pwsh 7+, fallback to Windows PowerShell 5.1) ──
where pwsh >nul 2>&1
if %errorlevel% equ 0 (
    echo  [OK]  Launching k3d-setup.ps1 via PowerShell 7 ^(pwsh^)...
    echo.
    pwsh -NoProfile -ExecutionPolicy Bypass -File "%~dp0k3d-setup.ps1" %*
    goto :done
)

where powershell >nul 2>&1
if %errorlevel% equ 0 (
    echo  [OK]  Launching k3d-setup.ps1 via Windows PowerShell...
    echo.
    powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0k3d-setup.ps1" %*
    goto :done
)

echo  [ERROR] PowerShell not found on this system.
echo          Windows 10/11 ships Windows PowerShell 5.1 by default.
pause
exit /b 1

:done
echo.
pause
