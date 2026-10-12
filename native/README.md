# ScreenShield — حماية أصلية ضد تسجيل الشاشة (Android / iOS)

وحدة الحماية الأصلية لتطبيق ويانا. **الويب (Telegram WebApp) لا يقدر أن يمنع
تسجيل الشاشة على مستوى النظام** — هذه الوحدة هي الطبقة التي تفعل ذلك عندما
يُبنى التطبيق كتطبيق أصلي عبر Capacitor.

## ماذا تفعل كل منصة؟

| | Android | iOS |
|---|---|---|
| **المنع** | `FLAG_SECURE` — لقطة الشاشة ومسجل النظام يلتقطان **شاشة سوداء** (منع حقيقي على مستوى النظام) | لا يوجد منع نظامي من آبل |
| **الكشف** | Android 14+ : `registerScreenCaptureCallback` (لقطة/تسجيل) | `UIScreen.isCaptured` (تسجيل/مشاركة AirPlay) + إشعار اللقطة |
| **الرد** | — (الالتقاط نفسه ممنوع) | إيقاف فوري + تحذير + تبليغ `screen_recording`/`screenshot_attempt` (10 نقاط = حد الحظر) |

## التثبيت

```bash
cd frontend
npx cap add android
npx cap add ios
```

### Android
1. انسخ `android/ScreenShieldPlugin.java` إلى:
   `android/app/src/main/java/iq/wiyana/screenshield/ScreenShieldPlugin.java`
2. استبدل `android/app/src/main/java/<package>/MainActivity.java` بمحتوى
   `android/MainActivity.secure.java` (عدّل سطر الـ package حسب مشروعك).
3. `npx cap sync android`

> ملاحظة: حزمة `@capacitor-community/privacy-screen` (مثبتة أصلاً في
> `package.json`) تفعل `FLAG_SECURE` أيضاً — الـ `MainActivity.secure.java`
> يضمن تفعيله **قبل رسم أول إطار** (مهم: بدونه قد تتسرب لقطة البداية).

### iOS
1. انسخ `ios/ScreenShieldPlugin.swift` إلى `ios/App/App/` وأضفه لـ Xcode target.
2. `npx cap sync ios`

لا حاجة لتعديل Podfile (يعتمد على Capacitor و UIKit فقط).

## كيف يعمل مع الواجهة؟

`frontend/src/native/screenShield.ts` هو الجسر:
- يعمل **فقط** على المنصات الأصلية (`Capacitor.isNativePlatform()`) — على الويب
  يرجع فوراً وتستمر طبقات حماية الويب.
- `ProtectedVideoPlayer.vue` يشغّله عند بدء التشغيل ويوقفه عند الخروج.
- عند تأكيد الالتقاط: إيقاف الفيديو + إلغاء الجلسة + شاشة تحذير حرجة + تبليغ
  الباك إند (`screen_recording` = 10 نقاط → حظر فوري عند حد الـ 10).

## حدود صريحة (مهم)

- **أندرويد**: `FLAG_SECURE` يمنع مسجل النظام ولقطات الشاشة داخل التطبيق، لكنه
  لا يمنع التصوير **بهاتف ثانٍ** — تبقى العلامة المائية المحروقة بالفيديو
  (اسم الطالب + رقمه) هي الرادع والتتبع.
- **iOS**: المنع مستحيل نظامياً — الكشف والرد الفوري + العلامة المائية.
- بعض أجهزة أندرويد المعدّلة (root) قد تتجاوز `FLAG_SECURE` — لا يوجد حل كامل
  ضد جهاز مخترق، والعلامة المائية تبقى خط الدفاع الأخير.
