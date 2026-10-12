# =============================================================================
# WEYANA — SUPABASE SETUP & LINK SCRIPT
# Project Reference: hixigpxoxfnumnjhbnmj
# =============================================================================

Write-Host "==================================================" -ForegroundColor Cyan
Write-Host "   WEYANA | ويانا — إعداد وربط قاعدة بيانات Supabase" -ForegroundColor Yellow
Write-Host "   Project Ref: hixigpxoxfnumnjhbnmj" -ForegroundColor Green
Write-Host "==================================================" -ForegroundColor Cyan
Write-Host ""

# 1. Login to Supabase
Write-Host "[1/3] تسجيل الدخول إلى Supabase CLI..." -ForegroundColor Cyan
Write-Host "سيتم فتح المتصفح لتسجيل الدخول أو طلب Access Token..." -ForegroundColor Gray
npx -y supabase login

# 2. Link Project
Write-Host ""
Write-Host "[2/3] ربط المشروع مع Supabase (hixigpxoxfnumnjhbnmj)..." -ForegroundColor Cyan
Write-Host "قد يطلب كلمة مرور قاعدة البيانات (Database Password)..." -ForegroundColor Gray
npx -y supabase link --project-ref hixigpxoxfnumnjhbnmj

# 3. Push Database Migrations (Schema)
Write-Host ""
Write-Host "[3/3] رفع الجداول والـ Schema إلى Supabase..." -ForegroundColor Cyan
npx -y supabase db push

Write-Host ""
Write-Host "==================================================" -ForegroundColor Green
Write-Host "   تم ربط وإعداد Supabase بنجاح! 🚀" -ForegroundColor Green
Write-Host "==================================================" -ForegroundColor Green
