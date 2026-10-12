# ويانا | Wiyana — منصة الكورسات (Telegram Mini App)

منصة كورسات للمنهج العراقي (تشمل السادس العلمي) تعمل كـ **Telegram Mini App**:
واجهة **Vue 3 + Vite + Capacitor**، وباك إند **Go** واحد، وقاعدة بيانات **SQLite**
(ملف واحد، بدون سيرفر خارجي)، وفيديوهات محمية بتشفير **AES-128 HLS** مع علامة
مائية لكل طالب.

---

## 🧱 البنية

```
wena-app-platfrom/
├── frontend/            # Vue 3 + TypeScript + Vite (الواجهة + لوحة التحكم)
├── backend-go/          # Go API server (الوحيد — البايثون والـ Convex محذوفات)
│   ├── cmd/server/      # نقطة الدخول + مسارات HTTP (admin / stream / playback)
│   ├── internal/
│   │   ├── db/          # SQLite: السكيما + كل عمليات CRUD للوحة التحكم
│   │   ├── video/       # الكتالوج + HLS/ffmpeg + منح المشاهدة
│   │   └── security/    # الجلسات + المنح + نقاط المخالفات + الحظر
│   └── static/spa/      # نسخة الواجهة المبنية (يخدمها السيرفر نفسه)
├── Dockerfile           # بناء واحد: فرونت + باك إند Go
├── render.yaml / koyeb.yaml   # نشر الباك إند (سيرفر دائم — ليس Vercel)
└── vercel.json          # نشر الواجهة على Vercel (تُبنى من frontend/)
```

**قاعدة البيانات:** SQLite (`data/wiyana.db` افتراضياً، يتغير عبر `SQLITE_PATH`).
الجداول: `users, subjects, courses, sections, lessons, attachments, enrollments,
subscriptions, progress, videos, audit_logs` — تُبنى تلقائياً أول تشغيل.

---

## 🚀 التشغيل محلياً

### المتطلبات
- Go 1.24+‎
- Node.js 18+‎
- ffmpeg (لتشفير الفيديو والعلامة المائية)

### الخطوات
```bash
# 1) بناء الواجهة
cd frontend
npm install
npm run build
cp -r dist ../backend-go/static/spa

# 2) بناء وتشغيل السيرفر
cd ../backend-go
go build -o server ./cmd/server

# 3) التشغيل (وضع التطوير يسمح بدخول تجريبي)
SQLITE_PATH=data/wiyana.db \
ALLOW_DEV_AUTH=true \
JWT_SECRET=غيّرني \
PORT=8000 \
./server
```

الواجهة: `http://localhost:8000` — الـ API: `http://localhost:8000/api/v1`

> في وضع التطوير (`ALLOW_DEV_AUTH=true`) تقدر تدخل كمدير عبر
> `POST /api/v1/auth/dev-login` بدون تليجرام. **أطفئه في الإنتاج.**

---

## ⚙️ متغيرات البيئة

| المتغير | مطلوب | الوصف |
|---|---|---|
| `TELEGRAM_BOT_TOKEN` | إنتاج | توكن بوت تليجرام (للتحقق من الدخول) |
| `BOT_ADMIN_IDS` | إنتاج | أرقام تليجرام المدراء، مفصولة بفواصل (يأخذون دور owner) |
| `JWT_SECRET` | دائماً | سر توقيع الجلسات (عشوائي طويل) |
| `SQLITE_PATH` | اختياري | مسار ملف قاعدة البيانات (افتراضي `data/wiyana.db`) |
| `PORT` | اختياري | المنفذ (افتراضي `8000`) |
| `WEBAPP_URL` | إنتاج | رابط الواجهة (يُستخدم في روابط البوت) |
| `UPLOAD_DIR` | اختياري | مجلد رفع الفيديوهات (افتراضي `uploads`) |
| `HLS_CACHE_DIR` | اختياري | مجلد مقاطع HLS المشفرة (افتراضي `hls_cache`) |
| `REDIS_URL` | اختياري | إن وُجد يُستخدم للجلسات/الكاش، وإلا ذاكرة محلية |
| `ALLOW_DEV_AUTH` | تطوير فقط | `true` يفعّل الدخول التجريبي |

---

## 🛠️ لوحة التحكم

`POST /api/v1/auth/admin-login` بالـ identifier + دور `admin/owner/teacher`
(أو دخول تليجرام برقم ضمن `BOT_ADMIN_IDS`).

أهم المسارات (`/api/v1`):

| المجموعة | أمثلة |
|---|---|
| المواد/الكورسات | `GET/POST /subjects` ،`GET/POST /courses` ،`PATCH/DELETE /courses/{id}` ،`POST /courses/{id}/publish\|archive` |
| الأقسام/الدروس | `POST /courses/{id}/sections` ،`POST /sections/{id}/lessons` ،`POST /courses/{id}/sections/reorder` ،`POST /lessons/{id}/reorder` |
| الطلاب/المدرسين | `GET /students` ،`POST /students/{id}/suspend\|unsuspend` ،`POST /students/{id}/enrollments` ،`GET/POST /teachers` |
| الإحصائيات | `GET /dashboard/overview` ،`GET /dashboard/analytics?period=7d` ،`GET /dashboard/recent-courses` |
| الرفع | `POST /lessons/{id}/upload-video` (multipart) ،`POST /lessons/{id}/attachments` |
| المشغّل | `GET /lessons/{id}/player` (بيانات الدرس + الكورس + التقدم) |

---

## 🛡️ حماية الفيديوهات ومنع تسريبها

### سلسلة المشاهدة (كلها مربوطة ببعض)
1. **جلسة تشغيل** `POST /api/v1/playback/session` — تحقق من الهوية + عدم الحظر + الصلاحية (اشتراك/درس مجاني).
2. **منحة بث** `POST /api/stream/hls/{video_id}` — مربوطة **ببصمة الجهاز** (`X-Device-ID`) وبوقت انتهاء قصير.
3. **HLS مشفر AES-128** — كل منحة لها **مفتاح خاص بها**؛ الـ playlist والمقاطع والمفتاح كلها تُطلب عبر مسار المنحة فقط.
4. **علامة مائية مزدوجة**: محروقة داخل الفيديو عبر ffmpeg (اسم الطالب + رقمه + الجلسة) + علامة متحركة داخل المشغّل.
5. **بدون منحة صالحة = 403** — لا توجد روابط mp4 مباشرة أبداً.

### ضد تصوير الشاشة — بصراحة
لا توجد تقنية ويب تمنع تصوير الشاشة على مستوى نظام التشغيل (أي شخص يقدر
يصور بهاتف ثانٍ). لذلك الحماية هنا **طبقات**، كل طبقة توقف ما تقدر عليه:

| الطبقة | ماذا تفعل | ماذا لا تقدر عليه |
|---|---|---|
| حجب مسجلات المتصفح | اعتراض `getDisplayMedia` و`captureStream` و`MediaRecorder` — أي محاولة تسجيل عبر المتصفح تُرفض فوراً وتُسجَّل مخالفة | مسجل خارجي / هاتف ثانٍ |
| علامة مائية محروقة | اسم الطالب ورقمه داخل الفيديو نفسه — أي تسريب **يكشف صاحبه** (أقوى رادع حقيقي) | لا شيء — تبقى حتى في التسريب |
| علامة متحركة | تتحرك دورياً داخل المشغّل — يصعب قصّها من التسجيل | — |
| كشف السلوك | إخفاء الصفحة، فتح DevTools، الخروج من ملء الشاشة، العبث بالمشغّل — كلها إشارات تُجمَع | — |
| نقاط المخالفات | لقطة=1، أدوات المطور=2، تسجيل=5 — تحذير عند 5، **حظر 30 يوم** عند 10 (يُطبَّق خادمياً فورياً) | — |
| ربط الجهاز | المنحة والمفتاح لا يعملان إلا على الجهاز الذي أصدرهما | مشاركة الحساب محدودة بجلسة واحدة نشطة |

> **لتطبيق أندرويد الأصلي (Capacitor):** عند بناء التطبيق (`npx cap add android`)
> أضف `FLAG_SECURE` في `MainActivity` — وهو **الطريقة الوحيدة** التي تمنع لقطات
> الشاشة والتسجيل فعلياً على أندرويد:
> ```java
> @Override protected void onCreate(Bundle savedInstanceState) {
>     super.onCreate(savedInstanceState);
>     getWindow().setFlags(WindowManager.LayoutParams.FLAG_SECURE,
>                          WindowManager.LayoutParams.FLAG_SECURE);
> }
> ```
> على iOS لا يوجد منع نظامي — فقط **كشف** اللقطة (`UIScreen.isCaptured` /
> `userDidTakeScreenshotNotification`) ويُفضَّل ربطه بنظام النقاط أعلاه.
>
> **نظام ScreenShield الجاهز** (`native/`): وحدة أصلية كاملة ضد التسجيل الداخلي —
> `ScreenShieldPlugin.java` (أندرويد: FLAG_SECURE + كشف الالتقاط في Android 14+)،
> `ScreenShieldPlugin.swift` (iOS: كشف التسجيل/المشاركة/اللقطة)، وجسر
> `frontend/src/native/screenShield.ts` مربوط بالمشغّل (إيقاف فوري + تبليغ
> `screen_recording`). راجع `native/README.md` للتثبيت.

---

## 🌐 النشر

- **الواجهة → Vercel:** بناء `frontend/` ونشر `dist` (انظر `vercel.json`
  للـ SPA rewrites). متغير البناء الوحيد: `VITE_API_URL` = رابط الباك إند.
- **الباك إند → Render/Koyeb/VPS (Docker):** يحتاج سيرفراً دائماً (SQLite ملف
  محلي + ffmpeg + تحويل HLS) — **لا يعمل على Vercel** لأنه serverless بدون
  قرص دائم. الإعدادات جاهزة في `render.yaml` و`koyeb.yaml` و`Dockerfile`.
- **بوت تليجرام:** اربط الـ WebApp في BotFather (`/newapp`) برابط Vercel،
  وضع دومين Vercel في `WEBAPP_URL` و`CORS_ORIGINS` للباك إند.

---

## 📁 هيكل الباك إند (Go)

```
backend-go/
├── cmd/server/
│   ├── main.go            # الإعداد + اختيار الكتالوج (SQLite أولاً)
│   ├── http_admin.go      # ~45 مساراً للوحة التحكم (CRUD كامل)
│   ├── http_playback.go   # الجلسات + المنح + بث HLS
│   └── http_spa.go        # تقديم الواجهة المبنية
└── internal/
    ├── db/
    │   ├── db.go          # فتح SQLite + السكيما + الترحيل التلقائي
    │   └── admin.go       # كل عمليات لوحة التحكم
    ├── video/
    │   ├── sqlite_catalog.go  # الكتالوج فوق SQLite (Find/CanAccess)
    │   └── stream.go          # تحويل HLS + تشفير AES-128 + watermark
    └── security/
        └── service.go     # الجلسات + المنح + النقاط/الحظر
```

## ✅ تم التحقق منه
- بناء Go نظيف (`go build ./...` + `go vet`) وبناء الواجهة نظيف (`vue-tsc` + `vite build`)
- سلسلة كاملة: دخول ← مادة ← كورس ← قسم ← درس ← نشر ← تسجيل طالب ← إحصائيات
- سلسلة الحماية: جلسة ← منحة مربوطة بالجهاز ← playlist مشفر ← مفتاح لكل منحة؛
  بدون منحة 403، جهاز خاطئ 403، درس مدفوع بدون اشتراك 403، المقاطع مشفرة فعلياً
