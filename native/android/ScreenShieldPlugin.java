package iq.wiyana.screenshield;

import android.os.Build;
import android.os.Handler;
import android.os.Looper;
import android.view.WindowManager;

import com.getcapacitor.JSObject;
import com.getcapacitor.Plugin;
import com.getcapacitor.PluginCall;
import com.getcapacitor.PluginMethod;
import com.getcapacitor.annotation.CapacitorPlugin;

/**
 * ScreenShield — native screen-capture protection for the Wiyana app.
 *
 * - enableSecure(): sets FLAG_SECURE on the activity window. This is the OS-level
 *   block: screenshots and the system screen recorder capture a BLACK screen for
 *   this app. This is the only true prevention of on-device ("internal")
 *   screen recording on Android.
 * - isCaptured(): best-effort "is being captured right now" (API 34+ only).
 * - Events: "screenCaptureChanged" (API 34+ registerScreenCaptureCallback —
 *   fires for screenshots AND screen recordings), "screenshotTaken" is not
 *   needed on Android because FLAG_SECURE already blocks the capture itself.
 *
 * Install: copy this file into android/app/src/main/java/iq/wiyana/screenshield/
 * and register it in MainActivity (see MainActivity.secure.java).
 */
@CapacitorPlugin(name = "ScreenShield")
public class ScreenShieldPlugin extends Plugin {

    private volatile boolean captureActive = false;
    private final Handler mainHandler = new Handler(Looper.getMainLooper());

    @PluginMethod
    public void enableSecure(PluginCall call) {
        if (getActivity() == null) {
            call.reject("no activity");
            return;
        }
        getActivity().runOnUiThread(() -> {
            getActivity().getWindow().setFlags(
                    WindowManager.LayoutParams.FLAG_SECURE,
                    WindowManager.LayoutParams.FLAG_SECURE);
            call.resolve();
        });
    }

    @PluginMethod
    public void disableSecure(PluginCall call) {
        if (getActivity() == null) {
            call.reject("no activity");
            return;
        }
        getActivity().runOnUiThread(() -> {
            getActivity().getWindow().clearFlags(
                    WindowManager.LayoutParams.FLAG_SECURE);
            call.resolve();
        });
    }

    @PluginMethod
    public void isCaptured(PluginCall call) {
        JSObject ret = new JSObject();
        ret.put("captured", captureActive);
        // NOTE: on Android the enforcement is FLAG_SECURE (the capture is
        // blocked, so "is recording" is less meaningful). The API 34+
        // callback below flips captureActive for a short window so the web
        // layer can still warn + score the attempt.
        call.resolve(ret);
    }

    @Override
    public void load() {
        // Android 14 (API 34)+: the system notifies the app when the user
        // takes a screenshot or starts a screen recording while our activity
        // is visible — even though FLAG_SECURE makes the capture itself black.
        if (Build.VERSION.SDK_INT >= 34 && getActivity() != null) {
            getActivity().registerScreenCaptureCallback(
                    getActivity().getMainExecutor(),
                    () -> {
                        captureActive = true;
                        JSObject data = new JSObject();
                        data.put("captured", true);
                        data.put("platform", "android");
                        notifyListeners("screenCaptureChanged", data, true);
                        // The callback fires once per capture event; reset the
                        // flag after a window so polling sees a fresh state.
                        mainHandler.postDelayed(() -> captureActive = false, 8000);
                    });
        }
    }
}
