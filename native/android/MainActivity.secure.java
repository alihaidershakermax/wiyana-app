package iq.wiyana.app;

import android.os.Bundle;
import android.view.WindowManager;

import com.getcapacitor.BridgeActivity;

import iq.wiyana.screenshield.ScreenShieldPlugin;

/**
 * Drop-in MainActivity for the Wiyana Capacitor Android app with screen-capture
 * protection baked in.
 *
 * What it does:
 *  1. FLAG_SECURE on the window — screenshots and the system screen recorder
 *     capture a BLACK screen for this app. This is the only real OS-level
 *     prevention of on-device ("internal") screen recording on Android.
 *  2. Registers ScreenShieldPlugin so the web layer can query capture state
 *     and receive capture events (Android 14+).
 *
 * Install:
 *   npx cap add android
 *   # then replace android/app/src/main/java/<pkg>/MainActivity.java with this
 *   # file (adjust the package line), and copy ScreenShieldPlugin.java into
 *   # android/app/src/main/java/iq/wiyana/screenshield/
 */
public class MainActivity extends BridgeActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        // Block screenshots + screen recording at the OS level BEFORE anything
        // renders. There must be no frame drawn without FLAG_SECURE.
        getWindow().setFlags(
                WindowManager.LayoutParams.FLAG_SECURE,
                WindowManager.LayoutParams.FLAG_SECURE);

        registerPlugin(ScreenShieldPlugin.class);

        super.onCreate(savedInstanceState);
    }
}
